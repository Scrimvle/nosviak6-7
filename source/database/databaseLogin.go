package database

import "time"

// Login represents all the authentication attempts received
type Login struct {
	ID       int
	User     int
	Username string
	Created  int
	Terminal string
	IP       string
}

// NewLogin will perform the query to the database to insert the login request
func (db *Database) NewLogin(user *User, version,ip string) error {
	return db.execute("INSERT INTO `logins` (`id`, `user`, `username`, `created`, `terminal`, `ip`) VALUES (NULL, ?, ?, ?, ?, ?)", user.ID, user.Username, int(time.Now().Unix()), version, ip)
}

// GetLogins will fetch all the logins from the database
func (db *Database) GetLogins() ([]*Login, error) {
	context, err := db.db.Prepare("SELECT `id`, `user`, `username`, `created`, `terminal`, `ip` FROM `logins`")
	if err != nil {
		return make([]*Login, 0), err
	}

	defer context.Close()
	query, err := context.Query()
	if err != nil {
		return make([]*Login, 0), err
	}

	logins := make([]*Login, 0)

	defer query.Close()
	for query.Next() {
		login, err := db.scanLogin(query)
		if err != nil {
			return make([]*Login, 0), err
		}

		logins = append(logins, login)
	}

	return logins, nil
}

// GetLogins will fetch all the logins from the database who are assigned to a specific user
func (db *Database) GetUserLogins(user *User) ([]*Login, error) {
	context, err := db.db.Prepare("SELECT `id`, `user`, `username`, `created`, `terminal`, `ip` FROM `logins` WHERE `user` = ?")
	if err != nil {
		return make([]*Login, 0), err
	}

	defer context.Close()
	query, err := context.Query(user.ID)
	if err != nil {
		return make([]*Login, 0), err
	}

	logins := make([]*Login, 0)

	defer query.Close()
	for query.Next() {
		login, err := db.scanLogin(query)
		if err != nil {
			return make([]*Login, 0), err
		}

		logins = append(logins, login)
	}

	return logins, nil
}


// scanLogin will scan the database response for the login
func (database *Database) scanLogin(row QueryScan) (*Login, error) {
	var login *Login = new(Login)
	if err := row.Scan(&login.ID, &login.User, &login.Username, &login.Created, &login.Terminal, &login.IP); err != nil {
		return nil, err
	}

	return login, nil
}