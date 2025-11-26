package source

import "encoding/json"

type Method struct {
	Disabled        bool          `json:"disabled"`
	Description     string        `json:"description"`
	Access          []string      `json:"access"`
	IPAllowed       bool          `json:"ip_allowed"`
	URLAllowed      bool          `json:"url_allowed"`
	DefaultPort     int           `json:"default_port"`
	DefaultDuration int           `json:"default_duration"`
	Options         MethodOptions `json:"options"`
	URLEncoding     bool          `json:"url_encoding"`
	URLs            []string      `json:"urls"`
}

type MethodOptions struct {
	EnhancedMaxtime int                        `json:"enhanced_maxtime"`
	MaxtimeOverride int                        `json:"maxtime_override"`
	MinimumDuration int                        `json:"minimum_duration"`
	MethodGroup     string                     `json:"method_group"`
	API             bool                       `json:"api"`
	Bot             bool                       `json:"bot"`
	BypassCooldown  bool                       `json:"bypass_cooldown"`
	OngoingCap      int                        `json:"ongoing_cap"`
	KeyValues       map[string]*MethodKeyValue `json:"key_value"`
	AttackSend      string                     `json:"attack_send"` // priority
}

type MethodKeyValue struct {
	Default   string `json:"default"`
	Type      string `json:"type"`
	IntMax    int    `json:"int_max"`
	StringMax int    `json:"string_max"`
	Required  bool   `json:"required"`
}

// Methods stores every single registered method
var Methods map[string]*Method = make(map[string]*Method)

// configureAttackMethods will attempt to parse all the floods for the cnc
func configureAttackMethods(content []byte) error {
	return json.Unmarshal(content, &Methods)
}
