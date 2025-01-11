package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
	"github.com/Mixturka/vm-hub/pkg/security"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
)

type AuthService struct {
	userServise    *UserService
	validate       *validator.Validate
	sessionManager *session.SessionManager
}

func NewAuthService(userService *UserService, sessionManager *session.SessionManager) *AuthService {
	return &AuthService{
		userServise:    userService,
		validate:       validator.New(),
		sessionManager: sessionManager,
	}
}

func (as *AuthService) Register(dto dtos.RegisterDto, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isExists, err := as.userServise.FindByEmail(ctx, dto.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to check user existence: %w", err)
		}
	}
	if isExists != nil {
		return errors.New("registration failed: user with this email already exists. Please try to use other email or login to the existing account")
	}

	newUser, err := as.userServise.CreateUser(ctx, dto.Email, dto.Password,
		dto.Name, "", entities.Credentials, false)
	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	return as.SaveSession(newUser, w)
}

func (as *AuthService) Login(dto dtos.LoginDto, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := as.userServise.FindByEmail(ctx, dto.Email)

	if err != nil || user.Password == "" {
		return errors.New("user wasn't found. Please check entered data")
	}

	if !security.ComparePasswords(user.Password, dto.Password) {
		return errors.New("wrong password")
	}

	return as.SaveSession(user, w)
}

func (as *AuthService) Logout(w http.ResponseWriter, r *http.Request) error {
	err := as.sessionManager.DestroySession(w, r)
	if err != nil {
		return errors.New("unable to stop session: possible internal server error or session was destroyed already")
	}

	return nil
}

func (as *AuthService) SaveSession(user *entities.User, w http.ResponseWriter) error {
	sessionData := map[string]interface{}{
		"userID": user.ID,
	}

	_, err := as.sessionManager.CreateSession(w, sessionData)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (as *AuthService) ValidateDto(dto interface{}) error {
	if err := as.validate.Struct(dto); err != nil {
		var errorMessages []string
		for _, validationErr := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages,
				fmt.Sprintf("Field '%s' failed validation. Rule: '%s', Value: '%v'",
					validationErr.Field(), validationErr.Tag(), validationErr.Value()))
		}
		return errors.New("validation failed: " + strings.Join(errorMessages, ", "))
	}
	return nil
}
