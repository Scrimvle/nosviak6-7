package interactions

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/terminal"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ADMIN    = "admin"
	MEMBER   = "member"
	MOD      = "mod"
	RESELLER = "reseller"

	// associations links the short string to the representation
	associations map[string]string = map[string]string{
		ADMIN:    "admin.tfx",
		MEMBER:   "member.tfx",
		MOD:      "mod.tfx",
		RESELLER: "reseller.tfx",
	}

	DEFAULT []string = []string{ADMIN, MOD, RESELLER, MEMBER}
)

// PopulateStringWithRoles will populate a string with the roles associated
func PopulateStringWithRoles(term *terminal.Terminal, roles ...string) (string, error) {
	global, err := database.DB.GetRoles(DEFAULT)
	if err != nil {
		return "", err
	}

	if len(roles) == 0 {
		roles = append(roles, global...)
	}

	placeholder, err := term.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "roles", "placeholder.tfx")
	if err != nil {
		return "", err
	}

	destination := make([]string, 0)
	for _, role := range global {
		if !slices.Contains(roles, role) {
			destination = append(destination, placeholder)
			continue
		}

		content, err := ExecuteRole(term, role)
		if err != nil {
			return "", err
		}

		destination = append(destination, content)
	}

	return strings.Join(destination, " "), nil
}

// ExecuteRole will grab the symbol and execute it
func ExecuteRole(term *terminal.Terminal, name string) (string, error) {
	index, ok := associations[strings.ToLower(name)]
	if !ok {
		/* checks for a preset association role */
		value, ok := source.OPTIONS.Config.Renders[filepath.Join(source.ASSETS, source.BRANDING, "roles", name + ".tfx")]
		if !ok || len(value) == 0 {
			index = "custom.tfx"
		} else {
			index = name + ".tfx"
		}
	}

	return term.ExecuteBrandingToString(map[string]any{"role": name}, source.ASSETS, source.BRANDING, "roles", index)
}