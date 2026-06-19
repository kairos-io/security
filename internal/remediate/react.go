package remediate

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
