package functions

import (
	"Nosviak4/source/database"
	"golang.org/x/exp/slices"
)

// CanAccessThemPermissions decides if the user can access the permissions
func CanAccessThemPermissions(user *database.User, permissions ...string) bool {
	if len(permissions) == 0 || len(permissions[0]) == 0 {
		return true 
	}

	for _, permission := range permissions {
		indexFlip := true
		switch permission[0] {

		case '!': // flip the result
			permission = permission[1:]
			indexFlip = !indexFlip
			fallthrough

		default:
			if permission[0] == '@' && user.Username == permission[1:] {
				return indexFlip
			}

			if slices.Contains(user.Roles, permission) {
				return indexFlip
			}
		}
	}

	return false
}

// RemovePermissionRights will directly remove their rights to that permission
func RemovePermissionRights(user *database.User, permission string) {
	index := 0

	for pos, role := range user.Roles {
		if role != permission {
			continue
		}

		index = pos
		break
	}

	user.Roles = append(user.Roles[:index], user.Roles[index+1:]...)
}