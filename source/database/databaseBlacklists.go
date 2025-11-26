package database

import (
	"errors"
	"time"
)

// Blacklist represents a target which is blacklisted within the database
type Blacklist struct {
	ID int
	User int
	Target string
	Created int64
	Expires int64
}

// NewBlacklist creates a new blacklisted target inside the database
func (db *Database) NewBlacklist(blacklist *Blacklist) error {
	if blacklist, err := db.GetBlacklistedTarget(blacklist.Target); err == nil && blacklist != nil {
		return errors.New("already blacklisted")
	}

	return db.execute("INSERT INTO `blacklists` (`id`, `user`, `target`, `created`, `expires`) VALUES (NULL, ?, ?, ?, ?)", blacklist.User, blacklist.Target, blacklist.Created, blacklist.Expires)
}

// GetBlacklistedTarget will return the object which shows if the target is blacklisted, follows expiry rules.
func (db *Database) GetBlacklistedTarget(target string) (*Blacklist, error) {
	query, err := db.db.Prepare("SELECT `id`, `user`, `target`, `created`, `expires` FROM `blacklists` WHERE `target` = ?")
	if err != nil {
		return nil, err
	}

	defer query.Close()
	blacklist, err := db.scanBlacklist(query.QueryRow(target))
	if err != nil || blacklist == nil || blacklist.Created + blacklist.Expires <= time.Now().Unix() {
		return nil, err
	}

	return blacklist, nil
}

// GetBlacklistedTargets will return a list of blacklisted targets
func (db *Database) GetBlacklistedTargets() ([]*Blacklist, error) {
	query, err := db.db.Prepare("SELECT `id`, `user`, `target`, `created`, `expires` FROM `blacklists`")
	if err != nil {
		return nil, err
	}

	defer query.Close()
	results, err := query.Query()
	if err != nil {
		return nil, err
	}

	blacklists := make([]*Blacklist, 0)

	defer results.Close()
	for results.Next() {
		blacklisted, err := db.scanBlacklist(results)
		if err != nil {
			return nil, err
		}

		blacklists = append(blacklists, blacklisted)
	}

	return blacklists, nil
}

// GetUserBlacklistedTargets will return a list of blacklisted targets from a given user
func (db *Database) GetUserBlacklistedTargets(user *User) ([]*Blacklist, error) {
	query, err := db.db.Prepare("SELECT `id`, `user`, `target`, `created`, `expires` FROM `blacklists` WHERE `created` + `expires` > ? AND `user` = ?")
	if err != nil {
		return nil, err
	}

	defer query.Close()
	queried, err := query.Query(time.Now().Unix(), user.ID)
	if err != nil {
		return nil, err
	}

	blacklists := make([]*Blacklist, 0)

	defer queried.Close()
	for queried.Next() {
		blacklisted, err := db.scanBlacklist(queried)
		if err != nil {
			return nil, err
		}

		blacklists = append(blacklists, blacklisted)
	}

	return blacklists, nil
}

// RemoveBlacklist will remove all the blacklists against that target inside the database.
func (db *Database) RemoveBlacklist(target string) error {
	return db.execute("DELETE FROM `blacklists` WHERE `target` = ?", target)
}

// scanBlacklist will scan into the object the results for the blacklist
func (db *Database) scanBlacklist(query QueryScan) (*Blacklist, error) {
	var dest *Blacklist = new(Blacklist)
	if err := query.Scan(&dest.ID, &dest.User, &dest.Target, &dest.Created, &dest.Expires); err != nil {
		return nil, err
	}

	return dest, nil
}