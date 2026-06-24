package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type contextKey string

const UserIDKey contextKey = "userID"

var ErrUnauthorized = errors.New("unauthorized")

type AuthController struct {
	service Service
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthController(service Service) *AuthController {
	return &(AuthController{service: service})
}

func (c *AuthController) RegisterController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds := new(credentials)
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusUnauthorized)
		return
	}
	err = c.service.RegisterUserService(ctx, creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Server Error", 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *AuthController) LoginController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds := new(credentials)
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	accTokenString, refreshTokenString, err := c.service.LoginUserService(ctx, creds.Email, creds.Password)
	if err != nil && !errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "Server Error", 500)
		return
	}
	if err != nil && errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	res := LoginResponse{AccessToken: accTokenString, RefreshToken: refreshTokenString}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(res)
}

func (c *AuthController) LogoutController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDstr, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	err = c.service.LogoutService(ctx, userID)
	if errors.Is(err, ErrUnauthorized) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *AuthController) DeleteController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds := new(credentials)
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	err = c.service.DeleteUserService(ctx, creds.Email, creds.Password)
	if err != nil && !errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "Server Error", 500)
		return
	}
	if err != nil && errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
