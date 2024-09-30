package account

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"emnavisa/webserver/infrastructure/kernel"
	"emnavisa/webserver/infrastructure/models"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(app *kernel.Application) *Service {
	return &Service{
		db: app.Database,
	}
}

func (s *Service) Authenticate(username, password string) (int, error) {
	var hashedPassword string
	var userID int
	query := "SELECT id, password FROM users WHERE username = $1"
	err := s.db.QueryRow(query, username).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found")
	} else if err != nil {
		return 0, err
	}
	sha := sha256.New()
	sha.Write([]byte(password))
	if hex.EncodeToString(sha.Sum(nil)) != hashedPassword {
		return 0, errors.New("invalid password")
	}
	return userID, nil
}
func (s *Service) UserCreate(username, password string) error {
	sha := sha256.New()
	sha.Write([]byte(password))
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := s.db.Exec(query, username, hex.EncodeToString(sha.Sum(nil)))
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}
	return nil
}

func (s *Service) GetUserByToken(token string) (models.Account, error) {
	var user models.Account
	query := "SELECT id, username FROM users WHERE access_token = $1"
	err := s.db.QueryRow(query, token).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (s *Service) StoreAccessToken(userID int) (string, error) {
	data := fmt.Sprintf("%d-%s", userID, time.Now().Format(time.RFC3339Nano))
	hash := md5.New()
	hash.Write([]byte(data))
	accessToken := hex.EncodeToString(hash.Sum(nil))

	query := "UPDATE users SET access_token = $1 WHERE id = $2"
	_, err := s.db.Exec(query, accessToken, userID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return accessToken, nil
}
