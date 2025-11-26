package database

import (
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// User maintains our user values
type User struct {
	ID         int
	Parent     int
	Username   string
	Password   []byte
	Key        []byte
	Roles      []string
	SSHKey     []byte
	Theme      string
	Maxtime    int
	Cooldown   int
	Conns      int
	Banned     bool
	Created    int64
	Expiry     int64
	API        bool
	APIKey     []byte
	NewUser    bool
	Telegram   int
	Sessions   int
	MaxAttacks int
}

// SYSTEM is the system administration user. it's never inserted into the SQL
var SYSTEM *User = &User{
	ID: 1,
}

// NewUser will attempt to create a new user inside the database
func (database *Database) NewUser(user, parent *User, hookFunc func(string, map[string]any)) error {
	state, err := database.GetUser(user.Username)
	if state != nil && err == nil {
		return err
	}

	/* check if they have the member role, meaning they have access. */
	if i := "member"; !slices.Contains(user.Roles, i) {
		user.Roles = append(user.Roles, i)
	}

	/* generates and sorts the byte information for the user */
	user.Key = *NewSalt(16)
	user.APIKey = make([]byte, 0)
	user.Created = time.Now().Unix()
	user.Password = NewHash(user.Password, &user.Key)

	/* modifies some last second information */
	user.NewUser = true
	user.SSHKey = make([]byte, 0)

	err = database.execute("INSERT INTO `users` (`id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", parent.ID, user.Username, user.Password, user.Key, strings.Join(user.Roles, ","), user.Theme, user.Maxtime, user.Cooldown, user.Conns, user.Banned, user.Created, user.Expiry, user.API, user.APIKey, user.Telegram, user.Sessions, user.MaxAttacks, user.SSHKey)
	if err != nil || hookFunc == nil {
		return err
	}

	hookFunc("user_created", map[string]any{"user": user.User(), "actor": parent.User()})
	return nil
}

// DeleteUser will attempt to delete the user from the database
func (database *Database) DeleteUser(user *User) error {
	return database.execute("DELETE FROM `users` WHERE `username` = ?", user.Username)
}

// EditUser will attempt to modify the user directly inside the database
func (database *Database) EditUser(user, actor *User, hookFunc func(string, map[string]any)) error {
	if hookFunc != nil {
		before, err := database.GetUser(user.Username)
		if err != nil {
			return err
		}

		hookFunc("user_edited", map[string]any{"before": before.User(), "user": user.User(), "actor": actor.User()})
	}

	return database.execute("UPDATE `users` SET `password` = ?, `roles` = ?, `themes` = ?, `maxtime` = ?, `cooldown` = ?, `conns` = ?, `banned` = ?, `expiry` = ?, `api` = ?, `apiKey` = ?, `newUser` = ?, `created` = ?, `telegramID` = ?, `maxSessions` = ?, `max_attacks` = ?, `sshkey` = ? WHERE `username` = ?", user.Password, strings.Join(user.Roles, ","), user.Theme, user.Maxtime, user.Cooldown, user.Conns, user.Banned, user.Expiry, user.API, user.APIKey, user.NewUser, user.Created, user.Telegram, user.Sessions, user.MaxAttacks, user.SSHKey, user.Username)
}

// BanUser will ban the account
func (database *Database) BanUser(user *User) error {
	return database.execute("UPDATE `users` SET `banned` = 1 WHERE `username` = ?", user.Username)
}

// UnBanUser will unban the account
func (database *Database) UnBanUser(user *User) error {
	return database.execute("UPDATE `users` SET `banned` = 0 WHERE `username` = ?", user.Username)
}

// GetUser will index at `username` for the user
func (database *Database) GetUser(user string) (*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `username` = ?")
	if err != nil {
		return nil, err
	}

	defer context.Close()
	return database.scanUser(context.QueryRow(user))
}

// GetUserWithID will index at `id` for the user
func (database *Database) GetUserWithID(id int) (*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `id` = ?")
	if err != nil {
		return nil, err
	}

	defer context.Close()
	return database.scanUser(context.QueryRow(id))
}

// GetUserWithPublicKey will index at `sshkey` for the user
func (database *Database) GetUserWithPublicKey(publicKey []byte) (*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `sshkey` = ?")
	if err != nil {
		return nil, err
	}

	defer context.Close()
	return database.scanUser(context.QueryRow(publicKey))
}

// GetUsers will attempt to retrieve all the users inside the database
func (database *Database) GetUsers() ([]*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users`")
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	query, err := context.Query()
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	users := make([]*User, 0)
	for query.Next() {
		user, err := database.scanUser(query)
		if err != nil {
			return make([]*User, 0), err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetBannedUsers will attempt to retrieve all the users inside the database which are banned
func (database *Database) GetBannedUsers() ([]*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `banned` = 1")
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	query, err := context.Query()
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	users := make([]*User, 0)
	for query.Next() {
		user, err := database.scanUser(query)
		if err != nil {
			return make([]*User, 0), err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetExpiredUsers will attempt to retrieve all the users inside the database which are expired
func (database *Database) GetExpiredUsers() ([]*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `created` + `expiry` < ?")
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	query, err := context.Query(time.Now().Unix())
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	users := make([]*User, 0)
	for query.Next() {
		user, err := database.scanUser(query)
		if err != nil {
			return make([]*User, 0), err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetActiveUsers will fetch all the active users from the database.
func (database *Database) GetActiveUsers() ([]*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `banned` = 0 AND `created` + `expiry` < ?")
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	query, err := context.Query(time.Now())
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	users := make([]*User, 0)
	for query.Next() {
		user, err := database.scanUser(query)
		if err != nil {
			return make([]*User, 0), err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserAsParentalFigure will attempt to find the user as a parent
func (database *Database) GetUserAsParentalFigure(username string, parent *User) (*User, error) {
	users, err := database.GetUsersAsParent(parent)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username != username {
			continue
		}

		return user, nil
	}

	return nil, errors.New("bad username")
}

// GetUsersAsParent will attempt to retrieve the children of the user
func (database *Database) GetUsersAsParent(parent *User) ([]*User, error) {
	if slices.Contains(parent.Roles, "admin") || database.IsSuperuser(parent) {
		return database.GetUsers()
	}

	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `parent` = ?")
	if err != nil {
		return make([]*User, 0), err
	}

	defer context.Close()
	query, err := context.Query(parent.ID)
	if err != nil {
		return make([]*User, 0), err
	}

	children := make([]*User, 0)
	defer query.Close()
	for query.Next() {
		user, err := database.scanUser(query)
		if err != nil {
			return make([]*User, 0), err
		}

		children = append(children, user)
		if user.Username == parent.Username {
			continue
		}

		grandchildren, err := database.GetUsersAsParent(user)
		if err != nil {
			return make([]*User, 0), err
		}

		children = append(children, grandchildren...)
	}

	return children, nil
}

// GetUserAPIKey will find the owner of that APIKey
func (database *Database) GetUserAPIKey(apiKey string) (*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `apiKey` = ? AND `api` = 1")
	if err != nil {
		return nil, err
	}

	defer context.Close()
	return database.scanUser(context.QueryRow(sha256.New().Sum([]byte(apiKey))))
}

// GetUserTelegram will find the owner of that telegram ID
func (database *Database) GetUserTelegram(id int) (*User, error) {
	context, err := database.db.Prepare("SELECT `id`, `parent`, `username`, `password`, `key`, `roles`, `themes`, `maxtime`, `cooldown`, `conns`, `banned`, `created`, `expiry`, `api`, `apiKey`, `newUser`, `telegramID`, `maxSessions`, `max_attacks`, `sshkey` FROM `users` WHERE `telegramID` = ?")
	if err != nil {
		return nil, err
	}

	defer context.Close()
	return database.scanUser(context.QueryRow(id))
}

type QueryScan interface {
	Scan(...any) error
}

// scanUser will use the response from a query and scan into user
func (database *Database) scanUser(row QueryScan) (*User, error) {
	pointUser := new(User)

	/* before we split */
	var roles string

	if err := row.Scan(&pointUser.ID, &pointUser.Parent, &pointUser.Username, &pointUser.Password, &pointUser.Key, &roles, &pointUser.Theme, &pointUser.Maxtime, &pointUser.Cooldown, &pointUser.Conns, &pointUser.Banned, &pointUser.Created, &pointUser.Expiry, &pointUser.API, &pointUser.APIKey, &pointUser.NewUser, &pointUser.Telegram, &pointUser.Sessions, &pointUser.MaxAttacks, &pointUser.SSHKey); err != nil {
		return nil, err
	}

	pointUser.Roles = strings.Split(roles, ",")
	return pointUser, nil
}

// GetRoles will attempt to fetch every single unique role inside the database
func (database *Database) GetRoles(standard []string) ([]string, error) {
	users, err := database.GetUsers()
	if err != nil {
		return make([]string, 0), err
	}

	/* ranges through all the users indexed */
	for _, user := range users {
		for _, role := range user.Roles {
			if slices.Contains(standard, role) {
				continue
			}

			standard = append(standard, role)
		}
	}

	return standard, nil
}

// IsSuperuser happens whenever a user is parental to themselves.
func (database *Database) IsSuperuser(user *User) bool {
	return user.ID == user.Parent
}
