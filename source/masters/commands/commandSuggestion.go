package commands

import (
	"Nosviak4/source/masters/sessions"
	"strings"
)

// IndexPrefixReturnsBundle will look through all the commands with that prefix and return their aliases
func (node *Command) IndexPrefixReturnsBundle(prefix string, ctx *ArgContext, session *sessions.Session, pos int) []string {
	recommended := make([]string, 0)
	if node := ctx.Pos(pos - 1); node != nil && node.Callback != nil {
		recommended = append(recommended, node.Callback(ctx, session, pos - 1)...)
	}

	if node.Callback != nil {
		recommended = append(recommended, node.Callback(ctx, session, pos)...)	
	}

	/* algorithm to check the depth which depends on tokens length */
	if len(node.Parents()) + 1 == len(ctx.Tokens) && len(node.Parents()) <= 0 || len(node.Parents()) > 0 && len(node.Subcommands) > 0 {
		for _, subcommand := range node.Subcommands {
			recommended = append(recommended, subcommand.Aliases...)
		}
	}

	/* worker which removes everything which isn't relevant to the command */
	removedInvalidSuggestions := make([]string, 0)
	for _, recommendation := range recommended {
		if !strings.HasPrefix(recommendation, prefix) {
			continue
		}

		removedInvalidSuggestions = append(removedInvalidSuggestions, recommendation)
	}

	return removedInvalidSuggestions
}
