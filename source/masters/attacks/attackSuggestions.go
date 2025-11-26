package attacks

import (
	"Nosviak4/source"
	"Nosviak4/source/functions/iplookup"
	"net"

	"golang.org/x/exp/slices"
)

// Suggest will suggest a methods based on the suggestions.json file
func (t *target) Suggest() (*source.Method, string, error) {
	ip, ok := t.ValidateWithEndpoints()
	if !ok || ip == nil {
		return nil, "", nil
	}

	if len(ip) == 0 {
		ip = append(ip, net.ParseIP(t.target))
	}

	data, err := iplookup.Lookup(ip[0].String())
	if err != nil || data == nil {
		return nil, "", nil
	}

	// stores all of the attack suggestions
	for methodName, suggestion := range source.Suggestions {
		method, ok := source.Methods[methodName]
		if !ok || !suggestion.Enabled {
			continue
		}

		// implements the logic to checking if its the right thing
		if !slices.Contains(suggestion.Organisations, data.Organization) && !slices.Contains(suggestion.ASNs, int(data.ASN)) {
			continue
		}

		return method, methodName, nil
	}

	return nil, "", nil
}
