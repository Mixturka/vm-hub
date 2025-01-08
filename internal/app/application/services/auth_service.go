package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
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
	ctx := context.Background()
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

func (as *AuthService) ValidateRegisterDto(dto dtos.RegisterDto) error {
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
