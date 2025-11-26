package autocomplete

import (
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"strings"
)

// NewUnixAutoCompleter will directly attach to the reader and offer completions for the session
func NewUnixAutoCompleter(session *sessions.Session) *terminal.AutoCompleter {
	return &terminal.AutoCompleter{
		Completer: func(content []byte, reader *terminal.Read) (bool, error) {
			args := strings.Split(string(content), " ")
			index := commands.ROOT.IndexCommands(args, session)
			if index == nil || len(strings.Split(string(content), " ")[0]) == 0 && len(strings.Split(string(content), " ")) > 1 {
				return false, nil
			}

			ctx, _ := index.ParseArgs(args[len(index.Parents()):], args, session, false)
			if ctx == nil {
				return false, nil
			}

			indexable := index.IndexPrefixReturnsBundle(args[len(args)-1], ctx, session, len(ctx.Tokens))
			if len(indexable) == 0 {
				return false, nil
			}

			pos := 0

			for {
				_, err := reader.ChangeInput([]byte(strings.Join(append(args[:len(args)-1], indexable[pos]), " ")))
				if err != nil {
					return false, err
				}

				context, err := session.Terminal.Signal.ReadWithContext(session.Reader.Context)
				if err != nil {
					return false, err
				}

				switch context[0] {

				case 9: /* in-charge of pos keeping with the cursor and tab */
					if pos + 1 >= len(indexable) {
						pos = 0
					} else {
						pos++
					}

				default:
					return reader.Buf(context, true)
				}
			}

		},
	}
}
