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

// swagger:model
type CreateUserInput struct {
	// user first name
	//
	// required: true
	FirstName string `validate:"required"`
	// user last name
	//
	// required: true
	LastName string `validate:"required"`
	// user email
	//
	// required: true
	Email string `validate:"required,email"`
	// user role [admin or contributor]
	//
	// required: true
	Role string `validate:"required,oneof=admin contributor"`
}

type CreateUserInputOption func(*CreateUserInput)

func NewCreateUserInput(opts ...CreateUserInputOption) *CreateUserInput {
	input := &CreateUserInput{}
	for _, opt := range opts {
		opt(input)
	}

	return input
}

func WithCreateUserInputFirstName(firstName string) CreateUserInputOption {
	return func(u *CreateUserInput) {
		u.FirstName = firstName
	}
}

func WithCreateUserInputLastName(lastName string) CreateUserInputOption {
	return func(u *CreateUserInput) {
		u.LastName = lastName
	}
}

func WithCreateUserInputEmail(email string) CreateUserInputOption {
	return func(u *CreateUserInput) {
		u.Email = email
	}
}

func WithCreateUserInputRole(role string) CreateUserInputOption {
	return func(u *CreateUserInput) {
		u.Role = role
	}
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
