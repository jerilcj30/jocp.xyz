package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/jerilcj30/jocp/helpers"
)

type SessionTable struct {
	ID int
	// Token is only set when creating a new session. When look up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token
	// this value is computed and not stored in the database
	Token     string
	UserID    int
	TokenHash string
}

type SesssionService struct {
	DB *sql.DB
}

// create a session
func (ss *SesssionService) Create(userID int) (*SessionTable, error) {
	//Todo: Create the session token
	token, err := helpers.GenerateSessionToken()

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO: Hash the session token
	session := SessionTable{
		Token:     token,
		UserID:    userID,
		TokenHash: ss.hash(token),
		//TODO: Set the TokenHash
	}

	// Try to update the users sessions
	// If err, create a new session

	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1 RETURNING id;`, session.UserID, session.TokenHash)

	err = row.Scan(&session.ID)

	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash) 
			VALUES ($1, $2) 
			RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create:%w", err)
	}

	return &session, nil
}

func (ss *SesssionService) User(token string) (*UserTable, error) {
	tokenHash := ss.hash(token)
	var user UserTable
	row := ss.DB.QueryRow(`
		SELECT user_id FROM sessions WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	row = ss.DB.QueryRow(`
		select email, password_hash
		FROM users WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

// Delete a session
func (ss *SesssionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	// Exec is used as no needs to be returned from the table
	_, err := ss.DB.Exec(`
			DELETE FROM sessions
			WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete:%w", err)
	}

	return nil
}

func (ss *SesssionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
