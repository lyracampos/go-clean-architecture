package usecases

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	mock "github.com/lyracampos/go-clean-architecture/internal/domain/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetUserExecute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockUserDatabaseGateway(ctrl)
	fmt.Println(c)

	type args struct {
		ctx   context.Context
		input *GetUserInput
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(userDatabase *mock.MockUserDatabaseGateway)
		want       GetUserOutput
		wantErr    bool
		err        error
	}{
		{
			name: "success getting user",
			args: args{
				ctx:   ctx,
				input: &GetUserInput{ID: 1},
			},
			beforeTest: func(userDatabase *mock.MockUserDatabaseGateway) {
				userDatabase.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(defaultGetUserResult())
			},
			want:    defaultGetUserOutput(),
			wantErr: false,
		},
		{
			name: "fail getting user when id not found",
			args: args{
				ctx:   ctx,
				input: &GetUserInput{ID: 1},
			},
			beforeTest: func(userDatabase *mock.MockUserDatabaseGateway) {
				userDatabase.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(getUserNotFoundResult())
			},
			wantErr: true,
			err:     domain.ErrUserDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDatabase := mock.NewMockUserDatabaseGateway(ctrl)

			w := &getUserUseCase{
				userDatabaseGateway: mockUserDatabase,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockUserDatabase)
			}

			getUser, err := w.Execute(tt.args.ctx, *tt.args.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("getUser.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(getUser, tt.want) {
				t.Errorf("registerUserUseCase.Execute() = %v, want %v", getUser, tt.want)
			}
			if tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("registerUserUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func defaultGetUserResult() (*entities.User, error) {
	return &entities.User{
		ID:        1,
		FirstName: "firstName",
		LastName:  "lastName",
		Email:     "useremail@domain.com",
		Role:      "admin",
	}, nil
}

func getUserNotFoundResult() (*entities.User, error) {
	return &entities.User{}, domain.ErrUserDoesNotExist
}

func defaultGetUserOutput() GetUserOutput {
	return GetUserOutput{
		ID:        1,
		FirstName: "firstName",
		LastName:  "lastName",
		Email:     "useremail@domain.com",
		Role:      "admin",
	}
}
