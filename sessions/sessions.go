package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Session struct {
	SessionID string `json:"session_id"`
	Expires   int64  `json:"expires"`
	// UserID    string `json:"user_id"`
	// Username  string `json:"username"`
	// Email     string `json:"email"`
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// you can set the expiration to an hour and extend it every 30 minutes but set an absolute expiration of 12 hours so sessions won't last for longer than that.
func CreateSession() *Session {
	id, err := GenerateSessionID()
	if err != nil {
		fmt.Println(err)
	}
	return &Session{
		SessionID: id,
		Expires:   int64(time.Now().Add(time.Hour * 12).Unix()),
	}
}

func GetSession(sessionID string) (*Session, bool) {
	return CreateSession(), true
}

func ValidateSessionID(sessionID string) (*Session, error) {
	session, ok := GetSession(sessionID)
	if !ok {
		return nil, errors.New("invalid session id")
	}
	if session.Expires < int64(time.Now().Unix()) {
		return nil, errors.New("session expired")
	}
	return session, nil
}
