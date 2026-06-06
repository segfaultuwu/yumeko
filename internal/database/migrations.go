package database

import "database/sql"

func Migrate(db *sql.DB) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			guild_id TEXT NOT NULL,
			content TEXT NOT NULL,
			done INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			guild_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS guild_settings (
					guild_id TEXT PRIMARY KEY,
					log_channel_id TEXT,
					welcome_channel_id TEXT,
					autorole_id TEXT,
					welcome_enabled INTEGER NOT NULL DEFAULT 0,
					welcome_message TEXT NOT NULL DEFAULT '👋 Welcome <@{user_id}> to **{server_name}**!',
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS self_roles (
					guild_id TEXT NOT NULL,
					role_id TEXT NOT NULL,
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (guild_id, role_id)
		);
		`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
