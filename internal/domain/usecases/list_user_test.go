package usecases

import (
	"context"
	"reflect"
	"testing"

	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	mock "github.com/lyracampos/go-clean-architecture/internal/domain/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListUserExecute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx   context.Context
		input *ListUserInput
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(userDatabase *mock.MockUserDatabaseGateway)
		want       ListUserOutput
		wantErr    bool
		err        error
	}{
		{
			name: "success listing users",
			args: args{
				ctx:   ctx,
				input: &ListUserInput{},
			},
			beforeTest: func(userDatabase *mock.MockUserDatabaseGateway) {
				userDatabase.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(defaultListUserResult())
			},
			want:    defaultListUsersOutput(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDatabase := mock.NewMockUserDatabaseGateway(ctrl)

			usecase := &listUserUseCase{
				userDatabaseGateway: mockUserDatabase,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockUserDatabase)
			}

			listUsersResult, err := usecase.Execute(tt.args.ctx, *tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("listUsers.Execute() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(listUsersResult, tt.want) {
				t.Errorf("listUsers.Execute() = %v, want %v", listUsersResult, tt.want)
			}
			if tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("listUsers.Execute() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func defaultListUserResult() ([]*entities.User, error) {
	return []*entities.User{
		{
			ID:        1,
			FirstName: "firstName_admin",
			LastName:  "lastName_admin",
			Email:     "admin@domain.com",
			Role:      "admin",
		},
		{
			ID:        2,
			FirstName: "firstName_contributor",
			LastName:  "lastName_contributor",
			Email:     "contributor@domain.com",
			Role:      "contributor",
		},
	}, nil
}

func defaultListUsersOutput() ListUserOutput {
	users := []*entities.User{
		{
			ID:        1,
			FirstName: "firstName_admin",
			LastName:  "lastName_admin",
			Email:     "admin@domain.com",
			Role:      "admin",
		},
		{
			ID:        2,
			FirstName: "firstName_contributor",
			LastName:  "lastName_contributor",
			Email:     "contributor@domain.com",
			Role:      "contributor",
		},
	}
	output := ListUserOutput{
		Users: users,
		Total: 2,
	}
	return output
}
