package database

import (
	"strings"
)

// DBFields is a collection of functions which allow users to perform complex database operations within the database.
var DBFields map[string]any = map[string]any{

	// countUsers will take a filter and return how many users meet that filter criteria.
	"countUsers": func(args ...string) int {
		if len(args) == 0 {
			users, err := DB.GetUsers()
			if err != nil {
				return 0
			}

			return len(users)
		}

		switch strings.ToLower(args[0]) {

		case "active":
			users, err := DB.GetActiveUsers()
			if err != nil {
				return 0
			}

			return len(users)

		case "banned":
			bannedUsers, err := DB.GetBannedUsers()
			if err != nil {
				return 0
			}

			return len(bannedUsers)

		case "expired":
			expiredUsers, err := DB.GetExpiredUsers()
			if err != nil {
				return 0
			}

			return len(expiredUsers)

		default:
			users, err := DB.GetUsers()
			if err != nil {
				return 0
			}

			return len(users)
		}
	},

	// ongoing is how many attacks are ongoing
	"ongoing": func() int {
		attacks, err := DB.GetOngoing()
		if err != nil {
			return 0
		}

		return len(attacks)
	},

	// ongoingUser will return all the ongoing attacks by a user
	"ongoingUser": func(username string) int {
		user, err := DB.GetUser(username)
		if err != nil {
			return 0
		}

		attacks, err := DB.GetOngoingUser(user)
		if err != nil {
			return 0
		}

		return len(attacks)
	},

	"attacks": func() int {
		sent, err := DB.GetAttacks()
		if err != nil {
			return 0
		}

		return len(sent)
	},

	"attacksUser": func(username string) int {
		attacks, err := DB.GetUserAttacks(username)
		if err != nil {
			return 0
		}

		return len(attacks)
	},
}