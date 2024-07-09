package employees

import (
	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) FindEmployeeByUserID(gctx golly.Context, userID uuid.UUID) (Employee, error) {
	args := m.Called(gctx, userID)
	return args.Get(0).(Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeesByIDS(gctx golly.Context, ids uuid.UUIDs) ([]Employee, error) {
	args := m.Called(gctx, ids)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeesForTeam(gctx golly.Context, teamID uuid.UUID, excludeEmployees ...uuid.UUID) ([]Employee, error) {
	args := m.Called(gctx, teamID, excludeEmployees)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeesByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
	args := m.Called(gctx, managerID, scopes)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeesByManagerUserID(gctx golly.Context, userID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) ([]Employee, error) {
	args := m.Called(gctx, userID, scopes)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeeByEmailAndOrganizationID(gctx golly.Context, email string, organizationID uuid.UUID) (Employee, error) {
	args := m.Called(gctx, email, organizationID)
	return args.Get(0).(Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeesByManagerAndIDS(gctx golly.Context, managerID uuid.UUID, employeeIDs ...uuid.UUID) ([]Employee, error) {
	args := m.Called(gctx, managerID, employeeIDs)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *MockEmployeeService) PluckEmployeeIDsByManagerID(gctx golly.Context, managerID uuid.UUID, scopes ...func(db *gorm.DB) *gorm.DB) (uuid.UUIDs, error) {
	args := m.Called(gctx, managerID, scopes)
	return args.Get(0).(uuid.UUIDs), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeeByID(gctx golly.Context, id uuid.UUID) (Employee, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(Employee), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeeEmailsBySearch(gctx golly.Context, name string) ([]string, error) {
	args := m.Called(gctx, name)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockEmployeeService) FindEmployeeByID_Unsafe(gctx golly.Context, id uuid.UUID) (Employee, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(Employee), args.Error(1)
}

func (m *MockEmployeeService) FindTeamsByOrganizationID(gctx golly.Context, organizationID uuid.UUID) ([]Team, error) {
	args := m.Called(gctx, organizationID)
	return args.Get(0).([]Team), args.Error(1)
}

func (m *MockEmployeeService) FindTeamByID(gctx golly.Context, id uuid.UUID) (Team, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(Team), args.Error(1)
}

var _ EmployeeService = &MockEmployeeService{}

func UseMockService(gctx golly.Context, mock *MockEmployeeService) golly.Context {
	gctx.Set(serviceCtxKey, mock)
	return gctx
}
