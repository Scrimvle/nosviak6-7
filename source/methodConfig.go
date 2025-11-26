package source

import "strings"

// MethodGroup represents a method group within the array
type MethodGroup struct {
	Name     string `json:"name"`
	Conns    int    `json:"conns"`
	Cooldown int    `json:"cooldown"`
}

// Attacks is the current structure used within attacks.toml
type Attacks struct {
	Attacks struct {
		PortThenDuration bool           `toml:"port_then_duration"`
		AttackPrefix     string         `toml:"attack_prefix"`
		KvPrefix         string         `toml:"kv_prefix"`
		Resolver         []string       `toml:"resolver"`
		Groups           []*MethodGroup `toml:"groups"`
		Timeout          int            `toml:"timeout"`
		SendDiv          int            `toml:"send_div"`
		APIIgnore        bool           `toml:"api_ignore"`
		Global           struct {
			Enabled  bool `toml:"enabled"`
			MaxTime  int  `toml:"max_time"`
			Conns    int  `toml:"conns"`
			Cooldown int  `toml:"cooldown"`
		} `toml:"global"`
	} `toml:"attacks"`
}

// AttackSuggestion is the config holder for recommendations on typing
type AttackSuggestion struct {
	Forced        bool     `json:"force"`
	Enabled       bool     `json:"enabled"`
	Organisations []string `json:"organisations"`
	ASNs          []int    `json:"asns"`
}

// MethodConfig the configuration for the methods
var MethodConfig *Attacks = new(Attacks)

// Suggestions stores all the suggestions for every method registered
var Suggestions map[string]*AttackSuggestion = make(map[string]*AttackSuggestion)

// FindGroup will return the group via the name presented
func (attacks *Attacks) FindGroup(name string) (*MethodGroup, bool) {
	for _, iter := range attacks.Attacks.Groups {
		if !strings.EqualFold(iter.Name, name) {
			continue
		}

		return iter, true
	}
	return nil, false
}
