package database

import "time"

type Warn struct {
	ID          int
	User        int
	Issuer      int
	Reason      string
	Created     int64
	WeightedFor int64
}

// CreateWarn will insert the warning into the database
func (database *Database) NewWarn(warn *Warn) error {
	return database.execute("INSERT INTO `warns` (`id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for`) VALUES (NULL, ?, ?, ?, ?, ?)", warn.User, warn.Issuer, warn.Reason, warn.Created, warn.WeightedFor)
}

// GetWarns will return all warnings inside the database
func (database *Database) GetWarns() ([]*Warn, error) {
	context, err := database.db.Prepare("SELECT `id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for` FROM `warns`")
	if err != nil {
		return make([]*Warn, 0), nil
	}

	defer context.Close()
	query, err := context.Query()
	if err != nil {
		return make([]*Warn, 0), nil
	}

	warns := make([]*Warn, 0)

	defer query.Close()
	for query.Next() {
		warn, err := database.scanWarns(query)
		if err != nil {
			return make([]*Warn, 0), err
		}

		warns = append(warns, warn)
	}

	return warns, nil
}

// GetOngoingWarnings will get all the warnings which are currently active
func (database *Database) GetOngoingWarnings(user *User) ([]*Warn, error) {
	context, err := database.db.Prepare("SELECT `id`, `user`, `issuer`, `reason`, `created_at`, `weighted_for` FROM `warns` WHERE `user` = ?")
	if err != nil {
		return make([]*Warn, 0), nil
	}

	defer context.Close()
	query, err := context.Query(user.ID)
	if err != nil {
		return make([]*Warn, 0), nil
	}

	warns := make([]*Warn, 0)

	defer query.Close()
	for query.Next() {
		warn, err := database.scanWarns(query)
		if err != nil {
			return make([]*Warn, 0), err
		}

		if warn.User != user.ID || warn.Created + warn.WeightedFor < time.Now().Unix() {
			continue
		}

		warns = append(warns, warn)
	}

	return warns, nil
}

// RemoveWarn will remove the warning which associated with that id
func (database *Database) RemoveWarn(warnID int) error {
	return database.execute("DELETE FROM `warns` WHERE `id` = ?", warnID)
}

// scanWarns will scan in from a database query result
func (database *Database) scanWarns(row QueryScan) (*Warn, error) {
	var warning *Warn = new(Warn)
	if err := row.Scan(&warning.ID, &warning.User, &warning.Issuer, &warning.Reason, &warning.Created, &warning.WeightedFor); err != nil {
		return nil, err
	}

	return warning, nil
}