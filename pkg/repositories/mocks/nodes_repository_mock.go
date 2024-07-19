package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

type MockNodesRepo struct {
	mock.Mock
}

func (m *MockNodesRepo) GetAllNodesByAccountID(accountID uint64) ([]models.Nodes, error) {
	args := m.Called(accountID)
	return args.Get(0).([]models.Nodes), args.Error(1)
}

func (m *MockNodesRepo) GetNodeByID(nodeID uint64) (models.Nodes, error) {
	args := m.Called(nodeID)
	return args.Get(0).(models.Nodes), args.Error(1)
}

func (m *MockNodesRepo) AddNode(node models.Nodes) error {
	args := m.Called(node)
	return args.Error(0)
}

func (m *MockNodesRepo) UpdateNode(node models.Nodes) error {
	args := m.Called(node)
	return args.Error(0)
}

// TODO: refactor
func (m *MockNodesRepo) DeleteNode(nodeID uint64) error {
	args := m.Called(nodeID)
	return args.Error(0)
}
