package database

import (
	"Nosviak4/source"
	"math"
	"net/url"
	"time"
)

type Attack struct {
	ID       int    `swash:"id"`
	User     int    `swash:"user"`
	Group    string `swash:"group"`
	Method   string `swash:"method"`
	Target   string `swash:"target"`
	Resolved string `swash:"resolved"`
	Port     int    `swash:"port"`
	Duration int    `swash:"duration"`
	Created  int    `swash:"created"`
}

// NewAttack creates a new attack inside the database
func (db *Database) NewAttack(attack *Attack) error {
	group, ok := source.MethodConfig.FindGroup(attack.Group)
	if !ok || group == nil {
		attack.Group = ""
	}

	// strips any path from the url
	if method, ok := source.Methods[attack.Method]; ok && method != nil && method.URLAllowed {
		url, err := url.Parse(attack.Target)
		if err != nil || url == nil {
			return db.execute("INSERT INTO `attacks` (`id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?)", attack.User, attack.Group, attack.Method, attack.Target, attack.Resolved, attack.Port, attack.Duration, attack.Created)
		}

		/* if host length is larger than 0 */
		if len(url.Host) > 0 {
			attack.Target = url.Host
		}
	}

	return db.execute("INSERT INTO `attacks` (`id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?)", attack.User, attack.Group, attack.Method, attack.Target, attack.Resolved, attack.Port, attack.Duration, attack.Created)
}

// GetOngoing will return a list of all ongoing attacks
func (db *Database) GetOngoing() ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix())
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingTarget will fetch all the ongoing attacks for that target specified
func (db *Database) GetOngoingTarget(target string) ([]*Attack, error) {
	if url, err := url.Parse(target); err == nil && url != nil && len(url.Host) > 0 {
		target = url.Host
	}

	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ? AND `target` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix(), target)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingUser will return a list of all ongoing attacks for a user
func (db *Database) GetOngoingUser(user *User) ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ? AND `user` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix(), user.ID)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingUserMethod will return a list of all ongoing attacks for a user
func (db *Database) GetOngoingUserMethod(user *User, method string) ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ? AND `user` = ? AND `method` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix(), user.ID, method)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingMethod will return a list of all ongoing attacks for a method
func (db *Database) GetOngoingMethod(method string) ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ? AND `method` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix(), method)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingGroup will return a list of all ongoing attacks for a group
func (db *Database) GetOngoingGroup(group string) ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` + `duration` >= ? AND `group` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Now().Unix(), group)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetAttacks will return a list of all attacks sent globally
func (db *Database) GetAttacks() ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks`")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query()
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetTodaysAttacks will return all the attacks launched today by that use 
func (db *Database) GetTodaysAttacks(user *User) ([]*Attack, error) {
	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `created` >= ? AND `user` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Unix(), user.ID)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetUserAttacks returns all attacks launched by a user
func (db *Database) GetUserAttacks(user string) ([]*Attack, error) {
	sender, err := db.GetUser(user)
	if err != nil {
		return make([]*Attack, 0), err
	}

	statement, err := db.db.Prepare("SELECT `id`, `user`, `group`, `method`, `target`, `resolved`, `port`, `duration`, `created` FROM `attacks` WHERE `user` = ?")
	if err != nil {
		return make([]*Attack, 0), err
	}

	query, err := statement.Query(sender.ID)
	if err != nil {
		return make([]*Attack, 0), err
	}

	attacks := make([]*Attack, 0)

	defer query.Close()
	for query.Next() {
		attack, err := db.scanAttack(query)
		if err != nil {
			return make([]*Attack, 0), err
		}

		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetOngoingMostUsedMethod will return the method which has the most ongoing attacks
func (db *Database) GetOngoingMostUsedMethod() (string, error) {
	ongoing, err := db.GetOngoing()
	if err != nil || len(ongoing) == 0 {
		return "", err
	}

	// stack will contain all the methods
	stack := make(map[string]int)
	stack[ongoing[0].Method] = 1

	for _, attack := range ongoing[0:] {
		if _, ok := stack[attack.Method]; ok {
			stack[attack.Method]++
			continue
		}

		stack[attack.Method] = 1
	}

	/* works out which method has the most ongoing attacks */
	var min, minValue = "", math.MinInt
	for key, val := range stack {
		if val >= minValue {
			min = key
			val = minValue
		}
	}

	return min, nil 
}

// scanAttack will directly scan the database query
func (db *Database) scanAttack(row QueryScan) (*Attack, error) {
	attack := new(Attack)

	if err := row.Scan(&attack.ID, &attack.User, &attack.Group, &attack.Method, &attack.Target, &attack.Resolved, &attack.Port, &attack.Duration, &attack.Created); err != nil {
		return nil, err
	}

	return attack, nil
}