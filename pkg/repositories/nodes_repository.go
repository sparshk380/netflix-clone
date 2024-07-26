package repositories

import (
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"

	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

type PgNodesRepo struct {
	DB *pg.DB
}

func NewPgNodesRepository(db *pg.DB) *PgNodesRepo {
	return &PgNodesRepo{
		DB: db,
	}
}

func (db *PgNodesRepo) GetAllNodesByAccountID(accountID uint64) ([]models.Nodes, error) {
	var nodes []models.Nodes

	if err := db.DB.Model(&nodes).Where("account_id = ?", accountID).Select(); err != nil {
		logrus.Errorf("failed to GetAllNodes. Err: %v", err)
		if errors.Is(err, pg.ErrNoRows) {
			return nodes, common.ErrRecordNotFound
		}
		return nodes, common.ErrInternalServerError
	}

	return nodes, nil
}

func (db *PgNodesRepo) GetNodeByID(nodeID uint64) (models.Nodes, error) {
	var node models.Nodes
	if err := db.DB.Model(&node).Where("node_id = ?", nodeID).Select(); err != nil {
		logrus.Errorf("failed to GetNodeByID. Err: %v", err)
		if errors.Is(err, pg.ErrNoRows) {
			return node, common.ErrRecordNotFound
		}
		return node, common.ErrInternalServerError
	}

	return node, nil
}

func (db *PgNodesRepo) AddNode(node models.Nodes) error {
	if _, err := db.DB.Model(&node).Insert(); err != nil {
		logrus.Errorf("failed to AddNode. Err: %v", err)

		if pgErr, ok := err.(pg.Error); ok {
			if ok && pgErr.IntegrityViolation() {
				return common.ErrDuplicateKey
			}
		}
		return common.ErrInternalServerError
	}

	return nil
}

func (db *PgNodesRepo) UpdateNode(node models.Nodes) error {
	result, err := db.DB.Model(&node).WherePK().UpdateNotZero()
	if err != nil {
		logrus.Errorf("failed to UpdateNode. Err: %v", err)

		if errors.Is(err, pg.ErrNoRows) {
			return common.ErrRecordNotFound
		}
		return common.ErrInternalServerError
	}

	if result.RowsAffected() == 0 {
		return common.ErrRecordNotFound
	}

	return nil
}

func (db *PgNodesRepo) DeleteNode(nodeID uint64) error {
	result, err := db.DB.Model(&models.Nodes{}).Where("node_id = ?", nodeID).Delete()
	if err != nil {
		logrus.Errorf("failed to DeleteNode. Err: %v", err)
		return common.ErrInternalServerError
	}

	if result.RowsAffected() == 0 {
		return common.ErrRecordNotFound
	}

	return nil
}
