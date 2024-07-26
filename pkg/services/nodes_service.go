package services

import (
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

type NodesService struct {
	repo NodesRepository
}

func NewNodesService(repo NodesRepository) *NodesService {
	return &NodesService{
		repo: repo,
	}
}

func (n *NodesService) GetAllNodesByAccountID(accountID uint64) ([]models.Nodes, error) {
	return n.repo.GetAllNodesByAccountID(accountID)
}

func (n *NodesService) GetNodeByID(nodeID uint64) (models.Nodes, error) {
	return n.repo.GetNodeByID(uint64(nodeID))
}

func (n *NodesService) AddNode(node models.Nodes) error {
	return n.repo.AddNode(node)
}

func (n *NodesService) UpdateNode(node models.Nodes) error {
	return n.repo.UpdateNode(node)
}

func (n *NodesService) DeleteNode(nodeID uint64) error {
	return n.repo.DeleteNode(nodeID)
}
