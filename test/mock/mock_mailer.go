package adapter

import (
	"go-gin-hexagonal/internal/domain/dto"

	"github.com/stretchr/testify/mock"
)

type MockMailerManager struct {
	mock.Mock
}

func (m *MockMailerManager) LoadEmailTemplate(templateName string, data any) (string, error) {
	args := m.Called(templateName, data)
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

func (m *MockEmailService) SendNewUserEmail(to string, data *dto.NewUserData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func (m *MockEmailService) SendVerifyEmail(to string, data *dto.VerifyEmailData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func (m *MockEmailService) SendRequestResetPassword(to string, data *dto.ResetPasswordData) error {
	args := m.Called(to, data)
	return args.Error(0)
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}
