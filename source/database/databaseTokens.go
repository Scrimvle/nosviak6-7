package database

import "errors"

type Token struct {
	ID      int
	Token   string
	Plan    string
	Created int64
	Expiry  int64
	Parent  int
	Owner   int
}

// NewToken creates a brand new token inside the database
func (conn *Database) NewToken(token *Token, self *User) error {
	if token, err := conn.GetToken(token.Token); err == nil && token != nil {
		return errors.New("token exists")
	}

	return conn.execute("INSERT INTO `tokens` (`id`, `token`, `plan`, `created`, `expiry`, `parent`, `owner`) VALUES (NULL, ?, ?, ?, ?, ?, 0)", token.Token, token.Plan, token.Created, token.Expiry, self.ID)
}

// ClaimToken will claim the token to the account specified
func (conn *Database) ClaimToken(token *Token, user *User) error {
	return conn.execute("UPDATE `tokens` SET `owner` = ? WHERE `token` = ? AND `id` = ?", user.ID, token.Token, token.ID)
}

// GetToken will index for the token inside the database
func (conn *Database) GetToken(token string) (*Token, error) {
	statement, err := conn.db.Prepare("SELECT `id`, `token`, `plan`, `created`, `expiry`, `parent`, `owner` FROM `tokens` WHERE `token` = ?")
	if err != nil {
		return nil, err
	}

	defer statement.Close()
	return conn.scanToken(statement.QueryRow(token))
}

// GetTokens will fetch all the tokens from the database
func (conn *Database) GetTokens(parent *User) ([]*Token, error) {
	statement, err := conn.db.Prepare("SELECT `id`, `token`, `plan`, `created`, `expiry`, `parent`, `owner` FROM `tokens` WHERE `parent` = ? AND `owner` = 0")
	if err != nil {
		return nil, err
	}

	defer statement.Close()
	query, err := statement.Query(parent.ID)
	if err != nil {
		return nil, err
	}
	
	tokens := make([]*Token, 0)

	for query.Next() {
		token, err := conn.scanToken(query)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// GetClaimedTokens will fetch all the claimed tokens
func (conn *Database) GetClaimedTokens(parent *User) ([]*Token, error) {
	statement, err := conn.db.Prepare("SELECT `id`, `token`, `plan`, `created`, `expiry`, `parent`, `owner` FROM `tokens` WHERE `parent` = ? AND `owner` >= 2")
	if err != nil {
		return nil, err
	}

	defer statement.Close()
	query, err := statement.Query(parent.ID)
	if err != nil {
		return nil, err
	}
	
	tokens := make([]*Token, 0)

	for query.Next() {
		token, err := conn.scanToken(query)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// scanToken will scan the database for the token presented
func (conn *Database) scanToken(query QueryScan) (*Token, error) {
	token := new(Token)
	if err := query.Scan(&token.ID, &token.Token, &token.Plan, &token.Created, &token.Expiry, &token.Parent, &token.Owner); err != nil {
		return nil, err
	}

	return token, nil
}