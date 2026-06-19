package remediate

import (
	"fmt"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type ReactionKind string

const (
	ReactReply  ReactionKind = "reply"
	ReactAdjust ReactionKind = "adjust"
	ReactClose  ReactionKind = "close"
	ReactNone   ReactionKind = "none"
)

type Reaction struct {
	Kind      ReactionKind
	ReplyBody string
	ToVersion string
}

func DecideReaction(c Classification) Reaction {
	switch c.Intent {
	case "nack":
		return Reaction{Kind: ReactClose, ReplyBody: c.Reply}
	case "request-change":
		if c.Version != "" {
			return Reaction{Kind: ReactAdjust, ToVersion: c.Version, ReplyBody: c.Reply}
		}
		return Reaction{Kind: ReactReply, ReplyBody: c.Reply}
	case "question":
		return Reaction{Kind: ReactReply, ReplyBody: c.Reply}
	default:
		return Reaction{Kind: ReactNone}
	}
}

const botLogin = "kairos-security-bot"

type Adjuster interface {
	Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error)
}

func seen(entry *state.LedgerEntry, id string) bool {
	for _, s := range entry.SeenComments {
		if s == id {
			return true
		}
	}
	return false
}

func ReactToComments(entry *state.LedgerEntry, gh ghclient.GitHub, cls CommentClassifier, adj Adjuster, prTitle, runID string, dryRun bool) error {
	if entry.State != "open" || entry.PRNumber == 0 {
		return nil
	}
	comments, err := gh.ListPRComments(entry.Repo, entry.PRNumber)
	if err != nil {
		return err
	}
	for _, cm := range comments {
		if cm.Author == botLogin || seen(entry, cm.ID) {
			continue
		}
		cl, err := cls.Classify(prTitle, cm.Author, cm.Body)
		if err != nil {
			// Do not guess; leave the comment unseen so a later run can retry.
			entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "needs-human", Detail: "classify failed: " + err.Error()})
			continue
		}
		r := DecideReaction(cl)
		action := string(r.Kind)
		if dryRun {
			fmt.Printf("[dry-run] would react to %s#%d comment %s: %s\n", entry.Repo, entry.PRNumber, cm.ID, action)
		} else {
			switch r.Kind {
			case ReactReply:
				if err := gh.PostPRComment(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
					return err
				}
			case ReactAdjust:
				updated, err := adj.Adjust(*entry, r.ToVersion, runID)
				if err != nil {
					return err
				}
				*entry = updated
				if r.ReplyBody != "" {
					if err := gh.PostPRComment(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
						return err
					}
				}
			case ReactClose:
				if err := gh.ClosePR(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
					return err
				}
				entry.State = "closed"
			case ReactNone:
			}
		}
		entry.SeenComments = append(entry.SeenComments, cm.ID)
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "reacted", Detail: action + " to " + cm.ID})
	}
	return nil
}
