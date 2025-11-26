package source

import "github.com/BurntSushi/toml"

// Presets are all the plans possible
var Presets map[string]*Preset = make(map[string]*Preset)

// Preset is a value of the Presets map
type Preset struct {
	Theme        string   `toml:"theme"`
	Roles        []string `toml:"roles"`
	Maxtime      int      `toml:"maxtime"`
	Cooldown     int      `toml:"Cooldown"`
	Concurrents  int      `toml:"concurrents"`
	Days         int      `toml:"days"`
	DailyAttacks int      `toml:"daily_attacks"`
}

// configurePlanPresets will promptly parse all the imported plan presets
func configurePlanPresets(jsonBytes []byte) error {
	return toml.Unmarshal(jsonBytes, &Presets)
}
