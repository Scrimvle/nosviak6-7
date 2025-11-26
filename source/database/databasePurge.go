package database

// Trunacate will remove everything from the table provided
func (database *Database) Truncate(table string) error {
	statement, err := database.db.Prepare("DELETE FROM " + table)
	if err != nil || statement == nil {
		return err
	}

	defer statement.Close()
	if _, err := statement.Exec(); err != nil {
		return err
	}

	return nil
}