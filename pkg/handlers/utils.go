package handlers

import (
	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/models"
	mockRepo "github.com/unification-com/unode-onboard-api/pkg/repositories/mocks"
	"github.com/unification-com/unode-onboard-api/pkg/services"
	mockSvc "github.com/unification-com/unode-onboard-api/pkg/services/mocks"
)

func IfStringIDExists(id string) bool {
	if id == "" {
		return false
	}

	return true
}

func SetupNodesService() (*NodesHandler, *mockRepo.MockNodesRepo) {
	mockRepo := new(mockRepo.MockNodesRepo)
	service := services.NewNodesService(mockRepo)
	handler := NewNodesHandler(service)

	return handler, mockRepo
}

func SetupNodesHandler() (*NodesHandler, *mockSvc.MockNodesService) {
	mockService := new(mockSvc.MockNodesService)
	handler := NewNodesHandler(mockService)
	return handler, mockService
}

func ValidateNode(node models.Nodes) error {
	if node.NodeID <= 0 {
		return common.ErrInvalidID
	}

	return nil
}
