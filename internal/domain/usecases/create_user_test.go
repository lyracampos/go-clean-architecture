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

func TestCreateUserExecute(t *testing.T) {
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
			name: "fail creating new user with empty first name",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputFirstName("")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'FirstName' should not be empty; "),
		},
		{
			name: "fail creating new user with empty last name",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputLastName("")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'LastName' should not be empty; "),
		},
		{
			name: "fail creating new user with empty email",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputEmail("")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'Email' should not be empty; "),
		},
		{
			name: "fail creating new user with invalid email",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputEmail("invalid@domain")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'Email' is invalid; "),
		},
		{
			name: "fail creating new user when email is already in use",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(),
			},
			beforeTest: func(userDatabase *mock.MockUserDatabaseGateway) {
				userDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(emailIsAlreadyInUseResult())
			},
			wantErr: true,
			err:     errors.New("failed to create user into database: email is arealdy in use"),
		},
		{
			name: "fail creating new user with empty role",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputRole("")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'Role' should not be empty; "),
		},
		{
			name: "fail creating new user with invalid role",
			args: args{
				ctx:   ctx,
				input: defaultCreateUserInput(WithCreateUserInputRole("invalid")),
			},
			beforeTest: nil,
			wantErr:    true,
			err:        errors.New("input is invalid: the field 'Role' is not valid, expected one of [admin contributor]; "),
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

func emailIsAlreadyInUseResult() (*entities.User, error) {
	return &entities.User{}, domain.ErrEmailAlreadyInUse
}
