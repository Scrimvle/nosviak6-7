package attacks

import (
	"Nosviak4/source"
	"bytes"
	"net"
	"net/url"
	"regexp"
	"unicode"
)

type target struct {
	target string
	method *source.Method

	/* Hosts is a list returned whenever a URL is passed */
	Hosts []net.IP 
}

// NewTarget creates a brand new target schema
func NewTarget(t string, m *source.Method) *target {
	return &target{
		target: t, method: m, Hosts: make([]net.IP, 0),
	}
}

// Validate will check if the target is valid or not.
func (t *target) Validate() bool {
	if net.ParseIP(t.target) != nil && t.method.IPAllowed {
		return true
	}

	u, err := url.ParseRequestURI(t.target)
	if u != nil && err == nil && t.method.URLAllowed {
		return true
	}

	u, err = url.Parse(t.target)
	if u != nil && err == nil && t.method.URLAllowed {
		t.target = u.Hostname()
	}

	t.Hosts, err = source.Resolver.LookupHost(t.target)
	if err != nil || !t.method.URLAllowed {
		return false
	}
	
	return len(t.Hosts) > 0
}

// ValidateWithEndpoints will validate the target and return the endpoints.
func (t *target) ValidateWithEndpoints() ([]net.IP, bool) {
	ok := t.Validate()
	return t.Hosts, ok
}

// HostStrings will return a list of strings which represent the resv endpoints.
func (t *target) HostStrings() []string {
	buf := make([]string, 0)
	if len(t.Hosts) <= 0 {
		return buf
	}

	for _, host := range t.Hosts {
		buf = append(buf, host.String())
	}

	return buf[:1]
}

// MaskTarget will return a string represents of target masked
func (t *target) MaskTarget(char string) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, content := range t.target {
		if !unicode.IsDigit(content) && !unicode.IsLetter(content) {
			buf.WriteRune(content)
			continue
		}

		buf.WriteString(char)
	}

	return buf.String()
}

// Blacklisted will check the type of the target and then attempt to continue
func (t *target) Blacklisted() bool {
	if len(t.Hosts) == 0 && !t.Validate() {
		return true
	}

	/* checks against the ip blacklist. */
	if t.method.IPAllowed {
		for _, ip := range source.OPTIONS.Strings("ips") {
			if ip == t.target {
				return true
			}

			re, err := regexp.Compile(ip)
			if err != nil {
				continue
			}

			/* if matches we know it's blacklisted. */
			if re.MatchString(t.target) {
				return true
			}
		}
	}

	if t.method.URLAllowed {
		host, err := url.ParseRequestURI(t.target)
		if err == nil && host != nil {
			t.target = host.Host
		}

		for _, ip := range source.OPTIONS.Strings("domains") {
			re, err := regexp.Compile(ip)
			if err != nil {
				continue
			}

			/* if matches we know it's blacklisted. */
			if re.MatchString(t.target) {
				return true
			}
		}
	}

	return false
}