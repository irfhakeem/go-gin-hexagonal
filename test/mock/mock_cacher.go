package adapter

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCacherService struct {
	mock.Mock
}

func (m *MockCacherService) Set(ctx context.Context, key string, value any) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCacherService) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacherService) Get(ctx context.Context, key string) (any, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (m *MockCacherService) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacherService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacherService) SetJSON(ctx context.Context, key string, data any, ttl time.Duration) error {
	args := m.Called(ctx, key, data, ttl)
	return args.Error(0)
}

func (m *MockCacherService) GetJSON(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func NewMockCacherService() *MockCacherService {
	return &MockCacherService{}
}
