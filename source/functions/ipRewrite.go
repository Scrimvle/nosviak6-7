package functions

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"strings"

	"golang.org/x/exp/slices"
)

// RewriteIP will decide if we should rewrite the Ip address
func RewriteIP(ip string, user *database.User) string {
	if slices.Contains(source.OPTIONS.Strings("users"), user.Username) && source.OPTIONS.Bool("enabled") || CanAccessThemPermissions(user, source.OPTIONS.Strings("roles")...) && source.OPTIONS.Bool("enabled") {
		return source.OPTIONS.String("placerholder_ip")
	}

	return strings.Join(strings.Split(ip, ":")[:strings.Count(ip, ":")], ":")
}