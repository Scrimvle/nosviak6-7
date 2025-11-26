package database

import (
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type PublicUser struct {
	ID           int                 `swash:"id"`
	Username     string              `swash:"username"`
	Theme        string              `swash:"theme"`
	Maxtime      int                 `swash:"maxtime"`
	MaxTimeFunc  func() string       `swash:"maxtime"`
	Cooldown     int                 `swash:"cooldown"`
	Conns        int                 `swash:"conns"`
	ConnsFunc    func() string       `swash:"conns"`
	HasRank      func(string) bool   `swash:"has_role"`
	LastLogin    *PublicLogin        `swash:"last_login"`
	Warns        func() int          `swash:"warns"`
	Expiry       int64               `swash:"expiry"`
	API          bool                `swash:"api"`
	NewUser      bool                `swash:"new_user"`
	Telegram     int                 `swash:"telegram"`
	Sessions     int                 `swash:"max_sessions"`
	Created      int64               `swash:"created"`
	LastAttack   *Attack             `swash:"last_attack"`
	Roles        func(string) string `swash:"roles"`
	AttacksToday func() int          `swash:"attacks_today"`
	MaxAttacks   func() string       `swash:"max_attacks"`
}

// User will publish a new PublicUser which is designed for Swash
func (u *User) User() *PublicUser {
	logins, err := DB.GetUserLogins(u)
	if err != nil || logins == nil || len(logins)-1 <= 0 {
		logins = append(logins, &Login{ID: 0, IP: "<n/a>", User: 0, Created: 0, Username: "<n/a>", Terminal: "<n/a>"})
	}

	lastAttack, err := DB.GetUserAttacks(u.Username)
	if err != nil || len(lastAttack) <= 0 {
		lastAttack = make([]*Attack, 1)
		lastAttack[0] = new(Attack)
	}

	dailyAttacks, err := DB.GetTodaysAttacks(u)
	if err != nil || dailyAttacks == nil {
		dailyAttacks = make([]*Attack, 0)
	}

	return &PublicUser{
		ID:           u.ID,
		API:          u.API,
		Conns:        u.Conns,
		Theme:        u.Theme,
		Expiry:       u.Created + u.Expiry,
		NewUser:      u.NewUser,
		Maxtime:      u.Maxtime,
		Username:     u.Username,
		Cooldown:     u.Cooldown,
		Created:      u.Created,
		Sessions:     u.Sessions,
		Telegram:     u.Telegram,
		LastLogin:    logins[len(logins)-1].Login(),
		LastAttack:   lastAttack[len(lastAttack)-1],

		// max_attacks()
		AttacksToday: func () int { 
			return len(dailyAttacks)
		},

		// has_role()
		HasRank: func(s string) bool {
			return slices.Contains(u.Roles, s)
		},

		// roles()
		Roles: func(sep string) string {
			return strings.Join(u.Roles, sep)
		},

		// conns()
		ConnsFunc: func() string {
			if u.Conns >= 1 {
				return strconv.Itoa(u.Conns)
			}

			return "∞"
		},

		// maxtime()
		MaxTimeFunc: func() string {
			if u.Maxtime >= 1 {
				return strconv.Itoa(u.Maxtime)
			}

			return "∞"
		},

		MaxAttacks: func() string {
			if u.MaxAttacks >= 1 {
				return strconv.Itoa(u.MaxAttacks)
			}

			return "∞"
		},

		// warns()
		Warns: func() int {
			ongoingWarns, err := DB.GetOngoingWarnings(u)
			if err == nil && len(ongoingWarns) > 0 {
				return len(ongoingWarns)
			}

			return 0
		},
	}
}

// PublicLogin will publish a new Public login structure which is designed for Swash
type PublicLogin struct {
	ID       int    `swash:"id"`
	Username string `swash:"username"`
	Terminal string `swash:"terminal"`
	Created  int    `swash:"created"`
	IP       string `swash:"ip"`
}

// Login will publish the swash connection structure
func (l *Login) Login() *PublicLogin {
	return &PublicLogin{
		ID:       l.ID,
		Username: l.Username,
		Terminal: l.Terminal,
		Created:  l.Created,
		IP:       l.IP,
	}
}
