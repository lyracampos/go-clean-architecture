package usecases

import (
	"context"
	"fmt"

	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/ports"
)

var _ ListUserUseCase = (*listUserUseCase)(nil)

type ListUserUseCase interface {
	Execute(ctx context.Context, input ListUserInput) (ListUserOutput, error)
}

type ListUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Role      string
}

type ListUserOutput struct {
	Users []*entities.User
	Total int
}

type listUserUseCase struct {
	userDatabaseGateway ports.UserDatabaseGateway
}

func NewListUserUseCase(userDatabaseGateway ports.UserDatabaseGateway) *listUserUseCase {
	return &listUserUseCase{
		userDatabaseGateway: userDatabaseGateway,
	}
}

func (u *listUserUseCase) Execute(ctx context.Context, input ListUserInput) (ListUserOutput, error) {
	users, err := u.userDatabaseGateway.ListUser(ctx, ports.ListUserFilter{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Role:      input.Role,
	})
	if err != nil {
		return ListUserOutput{}, fmt.Errorf("failed to list user from database: %w", err)
	}

	return ListUserOutput{
		Users: users,
		Total: len(users),
	}, nil
}
