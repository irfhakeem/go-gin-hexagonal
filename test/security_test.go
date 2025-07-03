package test

import "github.com/stretchr/testify/mock"

type MockSecurityService struct {
	mock.Mock
}

func (m *MockSecurityService) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityService) Verify(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func NewMockSecurityService() *MockSecurityService {
	return &MockSecurityService{}
}
