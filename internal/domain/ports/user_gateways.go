package ports

import (
	"context"

	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
)

type UserDatabaseGateway interface {
	ListUser(ctx context.Context) ([]*entities.User, error)
	GetUser(ctx context.Context, id int64) (*entities.User, error)
	InsertUser(ctx context.Context, user *entities.User) (*entities.User, error)
}
