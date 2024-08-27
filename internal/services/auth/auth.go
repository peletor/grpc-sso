package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"grpc-sso/internal/domain/models"
	"grpc-sso/internal/grpc/auth"
	"grpc-sso/internal/lib/jwt"
	"grpc-sso/internal/storage"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

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
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("Try to login user")
	log.Debug("User", slog.String("email", email))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("filed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("App not found",
				slog.Int("appID", appID),
				slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.Error("filed to get app",
			slog.Int("appID", appID),
			slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("filed to create token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
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

	log.Info("Registering new user")
	log.Debug("User", slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate password hash", slog.String("error", err.Error()))

		return models.EmptyUserID, fmt.Errorf("%s: %w", op, err)
	}

	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("User already exists", slog.String("error", err.Error()))

			return models.EmptyUserID, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("Failed to save user", slog.String("error", err.Error()))

		return models.EmptyUserID, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered")
	log.Debug("User", slog.String("email", email))

	return userID, nil
}

// IsAdmin checks if user is admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (isAdmin bool, err error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID))

	log.Info("Checking if user is admin")

	isAdmin, err = a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found", slog.String("error", err.Error()))

			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Filed to get user", slog.String("error", err.Error()))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
