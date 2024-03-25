package usecases

import (
	"context"
	"time"

	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/ports"
)

var _ GetUserUseCase = (*getUserUseCase)(nil)

type GetUserUseCase interface {
	Execute(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}

type GetUserInput struct {
	ID int64
}

type GetUserOutput struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type getUserUseCase struct {
	userDatabaseGateway ports.UserDatabaseGateway
}

func NewGetUserUseCase(userDatabaseGateway ports.UserDatabaseGateway) *getUserUseCase {
	return &getUserUseCase{
		userDatabaseGateway: userDatabaseGateway,
	}
}

func (u *getUserUseCase) Execute(ctx context.Context, input GetUserInput) (GetUserOutput, error) {
	user, err := u.userDatabaseGateway.GetUser(ctx, input.ID)
	if err != nil {
		return GetUserOutput{}, err
	}

	return u.fromEntityToOutput(user), nil
}

func (u *getUserUseCase) fromEntityToOutput(user *entities.User) GetUserOutput {
	return GetUserOutput{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
