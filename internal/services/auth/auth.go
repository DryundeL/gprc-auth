package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"grpc-auth/internal/domain/models"
	"grpc-auth/internal/lib/jwt"
	"grpc-auth/internal/storage"
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
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
)

// New returns a new instance of the Auth service
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

// Login checks if user with given credentials exist in the system
//
// If user exist, but pass is incorrect, returns error
// If user doesn`t exist, returns error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempt to login user")

	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to login", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	const op = "auth.registerNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err)

		return 0, fmt.Errorf("%s %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user not found", err)

			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to save user", err)

		return 0, fmt.Errorf("%s %w", op, err)
	}

	log.Info("user created")

	return id, nil
}

// IsAdmin checks if user is admin
func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	op := "auth.isAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userId", fmt.Sprint(userId)),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
