package postgres

import (
	"context"
	"strings"

	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/ports"
	"github.com/lyracampos/go-clean-architecture/internal/gateways/postgres/models"
	"github.com/uptrace/bun"
)

var _ ports.UserDatabaseGateway = (*userDatabase)(nil)

type userDatabase struct {
	Client *Client
}

func NewUserDatabase(client *Client) *userDatabase {
	return &userDatabase{
		Client: client,
	}
}

func (g *userDatabase) ListUser(ctx context.Context, filter ports.ListUserFilter) ([]*entities.User, error) {
	var modelList []models.User
	query := g.Client.DB.NewSelect().Model(&modelList)

	if filter.FirstName != "" {
		query.Where("first_name LIKE ?", "%"+filter.FirstName+"%")
	}

	if filter.LastName != "" {
		query.Where("last_name LIKE ?", "%"+filter.LastName+"%")
	}

	if filter.Email != "" {
		query.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	if filter.Role != "" {
		query.Where("? = ?", bun.Ident("role"), filter.Role)
	}

	if err := query.Scan(ctx); err != nil {
		if strings.Contains(err.Error(), NoRowsInResultSet) {
			return nil, newNoRowsError(models.UsersTableName, domain.ErrUserDoesNotExist)
		}

		return nil, newListError(models.UsersTableName, err)
	}

	list := make([]*entities.User, 0)
	for _, model := range modelList {
		list = append(list, model.ToEntity())
	}

	return list, nil
}

func (g *userDatabase) GetUser(ctx context.Context, id int64) (*entities.User, error) {
	model := models.User{}

	if err := g.Client.DB.NewSelect().Model(&model).
		Where("? = ?", bun.Ident("id"), id).
		Scan(ctx); err != nil {
		if strings.Contains(err.Error(), NoRowsInResultSet) {
			return nil, newNoRowsError(models.UsersTableName, domain.ErrUserDoesNotExist)
		}

		return nil, newListError(models.UsersTableName, err)
	}

	return model.ToEntity(), nil
}

func (g *userDatabase) InsertUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	model := models.NewUserModel(user)

	_, err := g.Client.DB.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), DuplicateKeyPrefix) {
			return nil, domain.ErrEmailAlreadyInUse
		}

		return nil, newInsertError(models.UsersTableName, err)
	}

	return model.ToEntity(), nil
}
