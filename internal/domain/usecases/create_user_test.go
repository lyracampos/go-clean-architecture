package usecases

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	mock "github.com/lyracampos/go-clean-architecture/internal/domain/ports/mocks"
	"go.uber.org/mock/gomock"
)

func Test_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockUserDatabaseGateway(ctrl)
	fmt.Println(c)

	type args struct {
		ctx   context.Context
		input *CreateUserInput
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(userDatabase *mock.MockUserDatabaseGateway)
		want       CreateUserOutput
		wantErr    bool
		err        error
	}{
		{
			name: "success creating new user",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(),
			},
			beforeTest: func(userDatabase *mock.MockUserDatabaseGateway) {
				userDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(defaultInsertedUserResult())
			},
			want: defaultCreateUserOutput(),
		},
		{
			name: "fail creating new user with empty name",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputFirstName("")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'FirstName' should not be empty; "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDatabase := mock.NewMockUserDatabaseGateway(ctrl)
			validator := domain.NewValidatorService()

			w := &createUserUseCase{
				userDatabase: mockUserDatabase,
				validator:    validator,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockUserDatabase)
			}

			createdUser, err := w.Execute(tt.args.ctx, *tt.args.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("registerUserUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(createdUser, tt.want) {
				t.Errorf("registerUserUseCase.Execute() = %v, want %v", createdUser, tt.want)
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

func defaultCreateUserInput(options ...CreateUserInputOption) *CreateUserInput {
	options = append([]CreateUserInputOption{
		WithCreateUserInputFirstName("Guilherme"),
		WithCreateUserInputLastName("Lyra"),
		WithCreateUserInputEmail("lyracampos@gmail.com"),
		WithCreateUserInputRole("admin"),
	}, options...)

	return NewCreateUserInput(options...)
}

func defaultCreateUserOutput() CreateUserOutput {
	return CreateUserOutput{
		ID:        1,
		FirstName: "firstName",
		LastName:  "lastName",
		Email:     "useremail@domain.com",
		Role:      "admin",
	}
}

func defaultInsertedUserResult() (*entities.User, error) {
	return &entities.User{
		ID:        1,
		FirstName: "firstName",
		LastName:  "lastName",
		Email:     "useremail@domain.com",
		Role:      "admin",
	}, nil
}
