package source

import (
	"Nosviak4/modules/goconfig"
	"Nosviak4/modules/gologr"
	"Nosviak4/modules/gotable2"
	"context"
	"encoding/json"

	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/bogdanovich/dns_resolver"
)

// ASSETS is a directory which maintains all our instancest
const ASSETS string = "assets"

// VERSION is the concurrent version of Nosviak4
const VERSION string = "v1.3.4"

// LICENSE is the provided domain for the licensing system
const LICENSE string = "nosviak.com"

// BRANDING references the file for branding
const BRANDING string = "branding"

// COMMANDS references the file for commands
const COMMANDS string = "commands"

// LOGGER defines the paramaters for logging into the terminal and files which surround it
var LOGGER *gologr.GoLogr = gologr.NewGoLogr(filepath.Join(ASSETS, "logs", "recycle"), os.Stdout)

// OPTIONS defines the global GoConfig object
var OPTIONS *goconfig.Options = &goconfig.Options{
	Config: goconfig.NewConfig(),
}

// TABLEVIEWS stores all the styles for the table config
var TABLEVIEWS map[string]*gotable2.Style = make(map[string]*gotable2.Style)

// Resolver is our sample base resolver incase we have issues with startup compatibility
var Resolver *dns_resolver.DnsResolver = dns_resolver.New([]string{"1.1.1.1", "8.8.8.8"})

// CancelContextSpinners will stop all spinners processes
var CancelContextSpinners context.CancelFunc = nil

// OpenOptions will trigger all the functions for GoConfig to init
func OpenOptions() error {
	MethodConfig = new(Attacks)
	Methods = make(map[string]*Method)

	/* NewConfig */
	savePoint := goconfig.NewConfig()
	savePoint.NewInclusion(".json", func(b []byte, s string, m map[string]any) error {
		switch filepath.Join(s) {
		case filepath.Join(ASSETS, COMMANDS, "styles.json"):
			return configureTableStyles(b)

		case filepath.Join(ASSETS, "attacks", "apis.json"):
			return configureAttackMethods(b)

		case filepath.Join(ASSETS, "attacks", "suggestions.json"):
			return json.Unmarshal(b, &Suggestions)

		default:
			return json.Unmarshal(b, &m)
		}
	})

	savePoint.NewInclusion(".toml", func(b []byte, s string, m map[string]any) error {
		switch filepath.Join(s) {

		case filepath.Join(ASSETS, "spinners.toml"):
			if CancelContextSpinners != nil {
				CancelContextSpinners()
			}

			ctx, canceller := context.WithCancel(context.Background())
			CancelContextSpinners = canceller
			return createSpinner(ctx, b)

		case filepath.Join(ASSETS, "attacks.toml"):
			if err := toml.Unmarshal(b, &m); err != nil {
				return err
			}

			return toml.Unmarshal(b, &MethodConfig)

		case filepath.Join(ASSETS, "plans.toml"):
			return configurePlanPresets(b)

		case filepath.Join(ASSETS, "themes.toml"):
			vals, err := configureThemes(b, m)
			if err != nil {
				return err
			}

			for k, v := range vals {
				m[k] = v
			}

			fallthrough

		default:
			return toml.Unmarshal(b, &m)
		}
	})

	err := savePoint.Parse(ASSETS)
	if err != nil {
		return err
	}

	OPTIONS, err = savePoint.Options()
	if err != nil {
		return err
	}

	Resolver = dns_resolver.New(OPTIONS.Strings("attacks", "resolver"))
	return nil
}
