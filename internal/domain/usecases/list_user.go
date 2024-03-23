package usecases

import (
	"context"

	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/ports"
)

var _ ListUserUseCase = (*listUserUseCase)(nil)

type ListUserUseCase interface {
	Execute(ctx context.Context, input ListUserInput) (ListUserOutput, error)
}

type ListUserInput struct {
	Email string
	Role  string
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
	users, err := u.userDatabaseGateway.ListUser(ctx)
	if err != nil {
		return ListUserOutput{}, nil
	}

	return ListUserOutput{
		Users: users,
		Total: len(users),
	}, nil
}
