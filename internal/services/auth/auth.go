package auth

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"grpc-sso/internal/domain/models"
	"grpc-sso/internal/grpc/auth"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (app models.App, err error)
}

// New returns a new instance of Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

var _ auth.Auth = &Auth{}

// Login checks is user exists.
// If user does not exist, returns error.
// If user exists, but password is incorrect, returns error.
func (a *Auth) Login(
	ctx context.Context,
	mail, password string,
	appID int,
) (token string, err error) {
	panic("implement me")
}

// RegisterNewUser registers new user and returns user ID.
// If user with this email already exists, returns error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))

	log.Info("Registering new user", slog.String("email", email))

	passHash, cryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if cryptErr != nil {
		log.Error("Failed to generate password hash", slog.String("error", cryptErr.Error()))

		return models.EmptyUserID, fmt.Errorf("%s: %w", op, cryptErr)
	}

	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("Failed to save user", slog.String("error", err.Error()))

		return models.EmptyUserID, fmt.Errorf("%s: %w", op, cryptErr)
	}

	log.Info("User registered")

	return userID, nil
}

// IsAdmin checks is user is admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (isAdmin bool, err error) {
	panic("implement me")
}
