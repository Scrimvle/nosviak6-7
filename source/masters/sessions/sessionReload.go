package sessions

import "Nosviak4/source"

// PushConcurrentChangesAcrossSessions will push updates for every single session online including their themes, prompts etc.
func PushConcurrentChangesAcrossSessions(exclude *Session) error {
	for _, session := range Sessions {
		if exclude != nil && exclude.Opened == session.Opened {
			continue
		}

		theme, err := source.GetTheme(session.User.Theme)
		if err != nil {
			return err
		}

		/* pushes the new theme */
		session.Theme = theme
		prompt, err := session.ExecuteBrandingToString(make(map[string]any), "prompt.tfx")
		if err != nil {
			return err
		}

		go session.Reader.Reskin(prompt)
	}

	return nil
}