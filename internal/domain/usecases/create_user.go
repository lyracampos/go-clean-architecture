package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/ports"
)

var _ CreateUserUseCase = (*createUserUseCase)(nil)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error)
}

type CreateUserInput struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required"`
	Role      string `validate:"required,oneof=admin contributor"`
}

type CreateUserOutput struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Role      string
	CreatedAt time.Time
}

type createUserUseCase struct {
	userDatabase ports.UserDatabaseGateway
	validator    domain.Validator
}

func NewCreateUserUseCase(userDatabase ports.UserDatabaseGateway, validator domain.Validator) *createUserUseCase {
	return &createUserUseCase{
		userDatabase: userDatabase,
		validator:    validator,
	}
}

func (u *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	err := u.validator.Validate(input)
	if err != nil {
		return CreateUserOutput{}, fmt.Errorf("input is invalid: %w", err)
	}

	newUser := u.fromInputToEntity(input)

	user, err := u.userDatabase.InsertUser(ctx, newUser)
	if err != nil {
		return CreateUserOutput{}, fmt.Errorf("failed to create user into database: %w", err)
	}

	return u.fromEntityToOutput(user), nil
}

func (u *createUserUseCase) fromInputToEntity(input CreateUserInput) *entities.User {
	return entities.NewUser(input.FirstName, input.LastName, input.Email, input.Role)
}

func (u *createUserUseCase) fromEntityToOutput(user *entities.User) CreateUserOutput {
	return CreateUserOutput{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
