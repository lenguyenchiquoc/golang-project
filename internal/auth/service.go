package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"managahub/pkg/models"
)

type AuthService struct {
	DB  *sql.DB
	JWT string
}

var activeSessions = struct {
	mu   sync.Mutex
	data map[string]string
}{
	data: make(map[string]string),
}

func NewAuthService(db *sql.DB, JWT string) *AuthService {
	return &AuthService{DB: db, JWT: JWT}
}

func (a *AuthService) login(req models.LoginRequest) (*models.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("Can not emptyu")
	}

	var user models.User
	query := `SELECT id, username, email, password_hash FROM users WHERE username = ?`
	err := a.DB.QueryRow(query, req.Username).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("1")
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("2")
	}

	token, err := a.generateJWT(user)
	if err != nil {
		return nil, errors.New("Can not create token" + err.Error())
	}

	activeSessions.mu.Lock()
	defer activeSessions.mu.Unlock()
	if oldToken, exist := activeSessions.data[user.ID]; exist {
		_, err := a.ValidateJWT(oldToken)
		if err == nil {
			return nil, errors.New("already_logged_in")
		}

	}
	activeSessions.data[user.ID] = token

	return &models.LoginResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}, nil

}

func Logout(userID string) {
	activeSessions.mu.Lock()
	defer activeSessions.mu.Unlock()

	delete(activeSessions.data, userID)
}

func (a *AuthService) register(request models.RegisterRequest) (*models.RegisterResponse, error) {
	if request.Username == "" || len(request.Username) < 6 {
		return nil, errors.New("1")
	}

	if validEmailRegex(request.Email) == false {
		return nil, errors.New("2")
	}

	if request.Password == "" || len(request.Password) < 10 {
		return nil, errors.New("3")
	}

	if request.RePassword == "" || request.RePassword != request.Password {
		return nil, errors.New("0")
	}

	if a.ExitOrNot(request.Email, request.Username) == true {
		return nil, errors.New("4")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("5")
	}

	userID := generateUserID()

	query := `INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)`
	_, err = a.DB.Exec(query, userID, request.Username, request.Email, string(hashedPassword))
	if err != nil {
		return nil, errors.New("Can not create user " + err.Error())
	}

	return &models.RegisterResponse{
		UserID:   userID,
		Username: request.Username,
		Email:    request.Email,
		Message:  "Register successfully",
	}, nil
}

func (a *AuthService) ExitOrNot(email string, username string) bool {
	query := "SELECT id FROM users WHERE email = ? OR username = ?"
	row := a.DB.QueryRow(query, email, username)

	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {

		return false
	}
	if err != nil {
		return false
	}
	return true
}

func validEmailRegex(email string) bool {
	valid := strings.Index(email, "@")
	if valid != -1 {
		return true
	}
	if strings.Contains(email[valid+1:], ".") {
		return true
	}
	return false
}

func generateUserID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}
func (s *AuthService) generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWT))
}
func (s *AuthService) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Not valid signing method")
		}
		return []byte(s.JWT), nil
	})

	if err != nil {
		return nil, errors.New("Not valid token" + err.Error())
	}

	if !token.Valid {
		return nil, errors.New("token is expired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Can not read claim token")
	}

	return claims, nil
}
