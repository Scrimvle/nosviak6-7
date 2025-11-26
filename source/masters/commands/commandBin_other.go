//go:build !linux

package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/masters/sessions"
)

// Unsupported operating system
func executeBinCommand(bin *source.Bin, session *sessions.Session) error {
	return session.ExecuteBranding(make(map[string]any), "feature_disabled.tfx")
}
