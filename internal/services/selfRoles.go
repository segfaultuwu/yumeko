package services

import "database/sql"

func AddSelfRole(db *sql.DB, guildID, roleID string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO self_roles (guild_id, role_id)
		VALUES (?, ?)
	`, guildID, roleID)

	return err
}

func RemoveSelfRole(db *sql.DB, guildID, roleID string) error {
	_, err := db.Exec(`
		DELETE FROM self_roles
		WHERE guild_id = ? AND role_id = ?
	`, guildID, roleID)

	return err
}

func IsSelfRole(db *sql.DB, guildID, roleID string) (bool, error) {
	var exists int

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM self_roles
			WHERE guild_id = ? AND role_id = ?
		)
	`, guildID, roleID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func ListSelfRoles(db *sql.DB, guildID string) ([]string, error) {
	rows, err := db.Query(`
		SELECT role_id
		FROM self_roles
		WHERE guild_id = ?
		ORDER BY created_at ASC
	`, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string

	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}

		roles = append(roles, roleID)
	}

	return roles, rows.Err()
}
