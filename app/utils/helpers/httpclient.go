package helpers

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type Requestor interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockReqestor struct {
	mock.Mock
}

func (m *MockReqestor) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

var _ Requestor = (*MockReqestor)(nil)
