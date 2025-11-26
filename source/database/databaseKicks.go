package database

import "time"

type Kick struct {
	ID          int
	User        int
	Issuer      int
	Reason      string
	Created     int64
	WeightedFor int64
}

// NewKick will insert the warning into the database
func (database *Database) NewKick(kick *Kick) error {
	return database.execute("INSERT INTO `kicks` (`id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for`) VALUES (NULL,?, ?, ?, ?, ?)", kick.User, kick.Issuer, kick.Reason, kick.Created, kick.WeightedFor)
}

// GetKicks will return all warnings inside the database
func (database *Database) GetKicks() ([]*Kick, error) {
	context, err := database.db.Prepare("SELECT `id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for` FROM `kicks`")
	if err != nil {
		return make([]*Kick, 0), nil
	}

	defer context.Close()
	query, err := context.Query()
	if err != nil {
		return make([]*Kick, 0), nil
	}

	kicks := make([]*Kick, 0)

	defer query.Close()
	for query.Next() {
		kick, err := database.scanKicks(query)
		if err != nil {
			return make([]*Kick, 0), err
		}

		kicks = append(kicks, kick)
	}

	return kicks, nil
}

// GetOngoingKicks will get all the kicks which are currently active
func (database *Database) GetOngoingKicks(user *User) ([]*Kick, error) {
	context, err := database.db.Prepare("SELECT `id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for` FROM `kicks` WHERE `user` = ?")
	if err != nil {
		return make([]*Kick, 0), nil
	}

	defer context.Close()
	query, err := context.Query(user.ID)
	if err != nil {
		return make([]*Kick, 0), nil
	}

	kicks := make([]*Kick, 0)

	defer query.Close()
	for query.Next() {
		kick, err := database.scanKicks(query)
		if err != nil {
			return make([]*Kick, 0), err
		}

		if kick.User != user.ID || kick.Created + kick.WeightedFor < time.Now().Unix() {
			continue
		}

		kicks = append(kicks, kick)
	}

	return kicks, nil
}

// RemoveKick will remove the warning which associated with that id
func (database *Database) RemoveKick(kickID int) error {
	return database.execute("DELETE FROM `kicks` WHERE `id` = ?", kickID)
}

// scanKicks will scan in from a database query result
func (database *Database) scanKicks(row QueryScan) (*Kick, error) {
	var kick *Kick = new(Kick)
	if err := row.Scan(&kick.ID, &kick.User, &kick.Issuer, &kick.Reason, &kick.Created, &kick.WeightedFor); err != nil {
		return nil, err
	}

	return kick, nil
}