package timescaledb

import (
	"github.com/HladCode/RMonitoringServer/internal/lib/e"
	"github.com/HladCode/RMonitoringServer/storage"
	"golang.org/x/crypto/bcrypt"
)

func (db *Database) AddUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return e.Wrap("can not hash password", err)
	}

	_, err = db.pool.Exec(db.cntxt, `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`,
		username, email, hashedPassword)
	if err != nil {
		return e.Wrap("Can not add data", err)
	}

	return nil
}

func (db *Database) GetUser(username string) (storage.User_data, error) {
	user_data := storage.User_data{
		User_name: username,
	}

	rows, err := db.pool.Query(db.cntxt, `SELECT id, email, password_hash FROM users WHERE username = $1`, username)
	if err != nil {
		return storage.User_data{}, e.Wrap("Can not get user", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user_data.ID, &user_data.Email, &user_data.Hashed_password); err != nil {
			return storage.User_data{}, e.Wrap("failed to scan row", err)
		}
	} else {
		return storage.User_data{}, e.Wrap("user was not found", err)
	}

	return user_data, nil
}

func (db *Database) IfUserInCompany(user_id int) (bool, error) {
	rows, err := db.pool.Query(db.cntxt, `SELECT company_name FROM users_companies WHERE user_id = $1`,
		user_id)
	if err != nil {
		return false, e.Wrap("Can not check if user in any company", err)
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (db *Database) GetUserID(username string) (int, error) {
	rows, err := db.pool.Query(db.cntxt, `SELECT id FROM users WHERE username = $1`,
		username)
	if err != nil {
		return 0, e.Wrap("Can not check if user in any company", err)
	}

	defer rows.Close()

	var user_id int

	if rows.Next() {
		if err := rows.Scan(&user_id); err != nil {
			return 0, e.Wrap("failed to scan row", err)
		}
	} else {
		return 0, e.Wrap("There is no user like this", err)
	}

	return user_id, nil
}
