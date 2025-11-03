package main

import (
	"context"

	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"petclinic/data"
	"petclinic/logger"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func getJWTSecret() []byte {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return []byte(v)
	}
	return []byte("dev_secret_change_me")
}

func hashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func checkPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func generateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling user registration request")

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid registration request: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.Warn("Registration attempt with empty email or password")
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	logger.Debug("Hashing password for user: %s", req.Email)
	hash, err := hashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		http.Error(w, "failed to process request", http.StatusInternalServerError)
		return
	}

	logger.Debug("Creating user in database: %s", req.Email)
	user, err := data.CreateUser(DB, req.Email, hash)

	if err != nil {
		logger.Error("Failed to create user %s: %v", req.Email, err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully registered user: %s (ID: %d)", req.Email, user.ID)
	tok, err := generateToken(user.ID, req.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{Token: tok})
}

func Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling login request")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid login request: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.Warn("Login attempt with empty email or password")
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	logger.Debug("Looking up user: %s", req.Email)
	user, err := data.FindUserByEmail(DB, req.Email)
	if err != nil {
		logger.Warn("Login failed - user not found: %s", req.Email)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	logger.Debug("Verifying password for user: %s", req.Email)
	if err := checkPassword(user.PasswordHash, req.Password); err != nil {
		logger.Warn("Login failed - invalid password for user: %s", req.Email)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	logger.Debug("Generating JWT token for user: %s (ID: %d)", req.Email, user.ID)
	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		logger.Error("Failed to generate token for user %s: %v", req.Email, err)
		http.Error(w, "failed to process request", http.StatusInternalServerError)
		return
	}

	logger.Info("Successful login for user: %s", req.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{Token: token})
}

type ctxKey string

const ctxUserIDKey ctxKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Processing authentication for request: %s %s", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format")
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		logger.Debug("Validating JWT token")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			logger.Warn("Invalid JWT token: %v", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.Warn("Invalid token claims")
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			logger.Warn("Invalid user ID in token")
			http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		userIDInt := int(userID)
		logger.Debug("Successfully authenticated user ID: %d", userIDInt)
		ctx := context.WithValue(r.Context(), ctxUserIDKey, userIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
