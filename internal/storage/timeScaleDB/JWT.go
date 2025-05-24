package timescaledb

import (
	"time"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
	"golang.org/x/crypto/bcrypt"
)

func (db *Database) SaveRefreshToken(hashed_token string, user_id int) error {
	_, err := db.pool.Exec(db.cntxt, `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		hashed_token, user_id, time.Now().Add(24*time.Hour).Format(time.RFC3339))
	if err != nil {
		return e.Wrap("Can not add refresh token", err)
	}

	return nil
}

type refreshToken struct {
	Id         int
	Token      string
	User_id    int
	Expires_at time.Time
	Created_at time.Time
}

func (db *Database) IfRefreshTokenValid(unhashed_token string, user_id int) (bool, error) {
	rows, err := db.pool.Query(db.cntxt, `SELECT * FROM refresh_tokens WHERE user_id = $1`,
		user_id)
	if err != nil {
		return false, e.Wrap("Can not get day data", err)
	}

	defer rows.Close()

	for rows.Next() {
		var refreshToken refreshToken
		if err := rows.Scan(
			&refreshToken.Id,
			&refreshToken.Token,
			&refreshToken.User_id,
			&refreshToken.Expires_at,
			&refreshToken.Created_at); err != nil {
			return false, e.Wrap("failed to scan row", err)
		}
		if err = bcrypt.CompareHashAndPassword([]byte(refreshToken.Token), []byte(unhashed_token)); err != nil {
			if refreshToken.Expires_at.Second() > time.Now().Second() {
				return true, nil
			} else {
				_, err := db.pool.Exec(db.cntxt, `DELETE FROM refresh_tokens WHERE id = $1`, refreshToken.Id)
				if err != nil {
					return false, e.Wrap("Can not delete expired refresh token", err)
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return false, e.Wrap("rows iteration error", err)
	}

	return false, nil
}
