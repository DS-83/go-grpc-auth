package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TokenRepoMock struct {
	mock.Mock
}

func (m *TokenRepoMock) RevokeToken(c context.Context, t string) error {
	args := m.Called(t)
	return args.Error(0)
}
func (m *TokenRepoMock) IsRevoked(c context.Context, t string) (bool, error) {
	args := m.Called(t)
	return args.Get(0).(bool), args.Error(1)
}
