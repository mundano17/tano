// Package auth: Handles all the login, delete, logout, refresh, etc
package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	db "tanoserver/internal/db/generated"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidCredentials error = errors.New("invalid credentials")
	ErrmismatchedToken    error = errors.New("token mismatched")
)

type AuthService struct {
	query *db.Queries
}

func NewAuthService(query *db.Queries) *AuthService {
	return &AuthService{query}
}

type Service interface {
	RegisterUserService(ctx context.Context, email string, pasword string) error
	LoginUserService(ctx context.Context, email string, password string) (accessTokenString string, refreshTokenString string, err error)
	DeleteUserService(ctx context.Context, email string, password string) error
	RefreshTokenService(ctx context.Context, refreshTokenString string) (string, string, error)
	LogoutService(ctx context.Context, userID int64) error
}

func (service *AuthService) RegisterUserService(ctx context.Context, email string, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	_, err = service.query.CreateUser(ctx, db.CreateUserParams{Email: email, Password: string(hash)})
	return err
}

func (service *AuthService) LoginUserService(ctx context.Context, email string, password string) (accessTokenString string, refreshTokenString string, err error) {
	res, err := service.query.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	match, err := argon2id.ComparePasswordAndHash(password, res.Password)
	if err != nil {
		return "", "", err
	}
	if !match {
		return "", "", ErrInvalidCredentials
	}
	accessTokenString, err = createAccessToken(res.UserID)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err = createRefreshToken(res.UserID)
	if err != nil {
		return "", "", err
	}
	err = service.query.UpdateRefreshToken(
		ctx,
		db.UpdateRefreshTokenParams{
			RefreshToken: pgtype.Text{
				String: hashToken(refreshTokenString),
				Valid:  true,
			},
			UserID: res.UserID,
		},
	)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	return accessTokenString, refreshTokenString, nil
}

func (service *AuthService) DeleteUserService(ctx context.Context, email string, password string) error {
	res, err := service.query.GetUserByEmail(ctx, email)
	errorInvalidCredentials := errors.New("invalid credentials")
	if err != nil {
		return errorInvalidCredentials
	}
	match, err := argon2id.ComparePasswordAndHash(password, res.Password)
	if err != nil {
		return err
	}
	if !match {
		return errorInvalidCredentials
	}
	err = service.query.DeleteUserByEmail(ctx, email)
	if err != nil {
		return errorInvalidCredentials
	}
	return nil
}

func (service *AuthService) RefreshTokenService(ctx context.Context, refreshTokenString string) (string, string, error) {
	userID, err := VerifyToken(refreshTokenString)
	if err != nil {
		return "", "", ErrmismatchedToken
	}
	hashedRefreshToken := hashToken(refreshTokenString)
	dbrefreshToken, err := service.query.GetRefreshTokenByUserID(ctx, userID)
	if err != nil {
		return "", "", err
	}
	if hashedRefreshToken != dbrefreshToken.String {
		_ = service.query.SetRefreshTokenToNULL(ctx, userID)
		return "", "", ErrmismatchedToken
	}
	newrefreshTokenString, err := createRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	newAccessTokenString, err := createAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	err = service.query.UpdateRefreshToken(
		ctx,
		db.UpdateRefreshTokenParams{
			RefreshToken: pgtype.Text{
				String: hashToken(newrefreshTokenString),
				Valid:  true,
			},
			UserID: userID,
		},
	)
	if err != nil {
		return "", "", err
	}
	return newAccessTokenString, newrefreshTokenString, nil
}

func (service *AuthService) LogoutService(ctx context.Context, userID int64) error {
	err := service.query.SetRefreshTokenToNULL(ctx, userID)
	return err
}

func hashToken(tokenstring string) string {
	hash := sha256.Sum256([]byte(tokenstring))
	hashString := fmt.Sprintf("%x", hash)
	return hashString
}
