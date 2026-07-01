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

type tokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errBody struct {
	Error string `json:"error"`
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
	if err != nil && errors.Is(err, ErrBadPassword) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "bad password"}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	if err != nil {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "server error"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *AuthController) LoginController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds := new(credentials)
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "invalid body"}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	accTokenString, refreshTokenString, err := c.service.LoginUserService(ctx, creds.Email, creds.Password)
	if err != nil && !errors.Is(err, ErrInvalidCredentials) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "server error"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	if err != nil && errors.Is(err, ErrInvalidCredentials) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "invalid credentials"}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	res := tokenResponseBody{AccessToken: accTokenString, RefreshToken: refreshTokenString}
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
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "invalid body"}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "invalid body"}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}
	err = c.service.LogoutService(ctx, userID)
	if errors.Is(err, ErrUnauthorized) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		errBody := errBody{Error: "invalid credentials"}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errBody)
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
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	if err != nil && errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

type refreshTokenBody struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *AuthController) RefreshTokenController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body refreshTokenBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	accTokenString, refreshTokenString, err := c.service.RefreshTokenService(ctx, body.RefreshToken)

	if err != nil && !errors.Is(err, ErrmismatchedToken) {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	if err != nil && errors.Is(err, ErrmismatchedToken) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	res := tokenResponseBody{AccessToken: accTokenString, RefreshToken: refreshTokenString}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(res)
}
