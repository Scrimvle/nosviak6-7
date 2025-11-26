package commands

import (
	"Nosviak4/source/masters/sessions"
	"fmt"
	"strings"
	"time"
)

// Arg represents a placeholder which can occur after a command parent
type Arg struct {
	Type
	Callback
	NotProvided
	Name        string
	OpenEnded   bool
	Description string
}

// NotProvided allows for a field which can dynamically prompt
type NotProvided func(*sessions.Session, []string) (string, error)

// ArgValue represents what is returned from the ParseArgs function
type ArgValue struct {
	Parent *Arg
	Values []*Token
}

// ArgContext represents the bundler for the ArgValue
type ArgContext struct {
	init     time.Time
	Command  *Command
	Header   string
	Tokens   []*Token
	Args     []*ArgValue
	Text     []string
}

// ParseArgs will concurrently range through all the tokens created from the args array and parse them
func (node *Command) ParseArgs(args, text []string, session *sessions.Session, ask bool) (*ArgContext, error) {
	tokens, err := Tokenize(strings.Join(args, " "), make([]*Token, 0))
	if err != nil {
		return nil, err
	}

	ctx := &ArgContext{
		Command: node,
		Tokens:  make([]*Token, 0),
		Args:    make([]*ArgValue, 0),
		Text:    text,
		init:    time.Now(),
	}

	for pos := 0; pos < len(tokens); pos++ {
		if pos >= len(node.Args) {
			break
		}

		if tokens[pos].Type == SPACE {
			tokens = append(tokens[:pos], tokens[pos+1:]...)
			pos--
			continue
		}

		/* blanks get attached to the value of the arg at hand */
		if tokens[pos].Type == BLANK {
			tokens[pos].Type = node.Args[pos].Type
		}

		/* type checks the value */
		if tokens[pos].Type != node.Args[pos].Type && node.Args[pos].Type != ANY {
			return ctx, fmt.Errorf("type mismatch in the argument %s", node.Args[pos].Name)
		}

		ctx.Args = append(ctx.Args, &ArgValue{
			Values: append(make([]*Token, 0), tokens[pos]),
			Parent: node.Args[pos],
		})

		ctx.Tokens = append(ctx.Tokens, tokens[pos])
		if !node.Args[pos].OpenEnded {
			continue
		}

		indexTokens := tokens[pos+1:]
		for i := 0; i < len(indexTokens); i++ {
			indexToken := indexTokens[i]
			if indexToken.Type == SPACE {
				indexTokens = append(indexTokens[:i], indexTokens[i+1:]...)
				i--
				continue
			}

			if indexToken.Type == BLANK {
				indexToken.Type = node.Args[pos].Type
			}

			if indexToken.Type != node.Args[pos].Type && node.Args[pos].Type != ANY {
				return ctx, fmt.Errorf("type mismatch in the argument %s", node.Args[pos].Name)
			}

			ctx.Tokens = append(ctx.Tokens, indexTokens[i])
			ctx.Args[len(ctx.Args) - 1].Values = append(ctx.Args[len(ctx.Args)-1].Values, indexTokens[i])
		}

		break
	}

	if len(ctx.Tokens) == 0 && len(node.Args) > 0 {
		ctx.Tokens = append(ctx.Tokens, &Token{
			Type: node.Args[0].Type,
			Literal: "",
		})
	}

	for _, arg := range node.Args[len(ctx.Args):] {
		if arg.NotProvided == nil {
			return ctx, fmt.Errorf("please provide the %s argument", arg.Name)
		}

		// default is an empty string
		var provided string = ""

		if ask {
			sample, err := arg.NotProvided(session, args)
			if err != nil {
				return ctx, err
			}

			provided = sample
		}

		tokens, err := Tokenize(provided, make([]*Token, 0))
		if err != nil {
			return ctx, err
		}

		if len(tokens) > 0 && tokens[0].Type != arg.Type && arg.Type != ANY {
			return nil, fmt.Errorf("type mismatch in the argument %s", arg.Name)
		}

		ctx.Args = append(ctx.Args, &ArgValue{
			Parent: arg,
			Values: tokens,
		})
	}

	return ctx, nil
}

// Pos will attempt to find and score the current pos based on the ctx value
func (ctx *ArgContext) Pos(node int) *Arg {
	if node >= len(ctx.Command.Args) || node < 0 {
		if node < 0 || len(ctx.Args) <= 0 || !ctx.Args[len(ctx.Args) - 1].Parent.OpenEnded {
			return nil
		}

		return ctx.Command.Args[len(ctx.Command.Args) - 1]
	}

	return ctx.Command.Args[node]
}
