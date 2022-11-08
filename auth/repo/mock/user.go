// nolint
package mock

import (
	"context"
	"example-grpc-auth/models"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) CreateUser(c context.Context, u string, p string) error {
	args := m.Called(u, p)
	return args.Error(0)
}
func (m *UserRepoMock) GetUser(c context.Context, u string, p string) (*models.User, error) {
	args := m.Called(u, p)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *UserRepoMock) UpdateUser(c context.Context, f *models.User, u *models.User) (*models.User, error) {
	args := m.Called(f, u)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *UserRepoMock) DeleteUser(c context.Context, u *models.User) error {
	args := m.Called(u)
	return args.Error(0)

}
