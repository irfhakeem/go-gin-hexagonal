package mock_external

import (
	"go-gin-hexagonal/internal/domain/ports/services"

	"github.com/stretchr/testify/mock"
)

type MockMailerManager struct {
	mock.Mock
}

func (m *MockMailerManager) LoadEmailTemplate(templateName string) (string, error) {
	args := m.Called(templateName)
	return args.String(0), args.Error(1)
}

func (m *MockMailerManager) SendEmail(to string, subject string, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func NewMockMailerManager() *MockMailerManager {
	return &MockMailerManager{}
}

type MockEmailService struct {
	*MockMailerManager
}

func (m *MockEmailService) SendNewUserEmail(to string, data *services.NewUserEmailData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func (m *MockEmailService) SendVerifyEmail(to string, data *services.VerifyEmailData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func (m *MockEmailService) SendRequestResetPassword(to string, data *services.ResetPasswordData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}
