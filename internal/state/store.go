package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Save[T any](dir, name string, v T) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(filepath.Join(dir, name), b, 0o644)
}

func Load[T any](dir, name string, v *T) error {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
