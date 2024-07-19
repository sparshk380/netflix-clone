package handlers

import "github.com/unification-com/unode-onboard-api/pkg/models"

type NodesServiceInterface interface {
	GetAllNodesByAccountID(accountID uint64) ([]models.Nodes, error)
	GetNodeByID(nodeID uint64) (models.Nodes, error)
	AddNode(node models.Nodes) error
	UpdateNode(node models.Nodes) error
	DeleteNode(nodeID uint64) error
}
