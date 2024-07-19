package tests

import (
	"time"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/models"
	"github.com/unification-com/unode-onboard-api/pkg/repositories/mocks"
)

func TestGetAllNodesByAccountID(t *testing.T) {

	testCases := []struct {
		name           string
		inputAccountID uint64
		expectedResult []models.Nodes
		expectedErr    error
	}{
		{
			name:           "Successful case",
			inputAccountID: 1,
			expectedResult: []models.Nodes{
				{
					NodeID:           1,
					AccountID:        1,
					ChainName:        "Ethereum",
					ChainType:        "mainnet",
					Status:           "active",
					PublicIP:         "192:22:22:21",
					RPCPort:          "5542",
					RPCWebSocketPort: "4545",
					RestHTTPPort:     "",
					CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
		{
			name:           "No record found case",
			inputAccountID: 1,
			expectedResult: []models.Nodes{},
			expectedErr:    common.ErrRecordNotFound,
		},
		{
			name:           "Internal server error case",
			inputAccountID: 1,
			expectedResult: []models.Nodes{},
			expectedErr:    common.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the behavior for GetAll method
			mockRepo := new(mocks.MockNodesRepo)
			mockRepo.On("GetAllNodesByAccountID", tc.inputAccountID).Return(tc.expectedResult, tc.expectedErr)

			// Call the tested function
			result, err := mockRepo.GetAllNodesByAccountID(tc.inputAccountID)

			// // Assert the results
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResult, result)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetNodeByID(t *testing.T) {
	testCases := []struct {
		name           string
		inputNodeID    uint64
		expectedResult models.Nodes
		expectedErr    error
	}{
		{
			name:           "Successful case",
			inputNodeID:    1,
			expectedResult: models.Nodes{NodeID: 1, ChainName: "ethereum", Status: "active", PublicIP: "127.0.0.1", RPCPort: "8545", RPCWebSocketPort: "8546", RestHTTPPort: "8547", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			expectedErr:    nil,
		},
		{
			name:           "Record not found case",
			inputNodeID:    2,
			expectedResult: models.Nodes{},
			expectedErr:    common.ErrRecordNotFound,
		},
		{
			name:           "Internal server error case",
			inputNodeID:    3,
			expectedResult: models.Nodes{},
			expectedErr:    common.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockNodesRepo)
			mockRepo.On("GetNodeByID", tc.inputNodeID).Return(tc.expectedResult, tc.expectedErr)

			result, err := mockRepo.GetNodeByID(tc.inputNodeID)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResult, result)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddNode(t *testing.T) {

	testCases := []struct {
		name        string
		inputNode   models.Nodes
		expectedErr error
	}{
		{
			name: "Successful case",
			inputNode: models.Nodes{
				NodeID:           1,
				ChainName:        "ethereum",
				Status:           "active",
				PublicIP:         "127.0.0.1",
				RPCPort:          "8545",
				RPCWebSocketPort: "8546",
				RestHTTPPort:     "8547",
				CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},

		{
			name: "Record already exists case",
			inputNode: models.Nodes{
				NodeID:           1,
				ChainName:        "ethereum",
				Status:           "active",
				PublicIP:         "127.0.0.1",
				RPCPort:          "8545",
				RPCWebSocketPort: "8546",
				RestHTTPPort:     "8547",
				CreatedAt:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: common.ErrDuplicateKey,
		},
		{
			name: "Internal server error case",
			inputNode: models.Nodes{
				NodeID:           2,
				ChainName:        "ethereum",
				Status:           "active",
				PublicIP:         "127.0.0.1",
				RPCPort:          "8545",
				RPCWebSocketPort: "8546",
				RestHTTPPort:     "8547",
				CreatedAt:        time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: common.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the behavior for GetAll method
			mockRepo := new(mocks.MockNodesRepo)
			mockRepo.On("AddNode", tc.inputNode).Return(tc.expectedErr)

			err := mockRepo.AddNode(tc.inputNode)
			assert.Equal(t, tc.expectedErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateNode(t *testing.T) {
	testCases := []struct {
		name        string
		inputNode   models.Nodes
		expectedErr error
	}{{
		name: "Successful case",
		inputNode: models.Nodes{
			NodeID:           1,
			ChainName:        "ethereum",
			Status:           "active",
			PublicIP:         "172.5.0.1",
			RPCPort:          "8545",
			RPCWebSocketPort: "8546",
			RestHTTPPort:     "8547",
			CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		expectedErr: nil,
	},
		{
			name: "Not found case",
			inputNode: models.Nodes{
				NodeID:           2,
				ChainName:        "ethereum",
				Status:           "active",
				PublicIP:         "172.5.0.1",
				RPCPort:          "8545",
				RPCWebSocketPort: "8546",
				RestHTTPPort:     "8547",
				CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: common.ErrRecordNotFound,
		},
		{
			name: "Internal server error case",
			inputNode: models.Nodes{
				NodeID:           1,
				ChainName:        "ethereum",
				Status:           "active",
				PublicIP:         "172.5.0.1",
				RPCPort:          "8545",
				RPCWebSocketPort: "8546",
				RestHTTPPort:     "8547",
				CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: common.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the behavior for GetAll method
			mockRepo := new(mocks.MockNodesRepo)
			mockRepo.On("UpdateNode", tc.inputNode).Return(tc.expectedErr)

			err := mockRepo.UpdateNode(tc.inputNode)
			assert.Equal(t, tc.expectedErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteNode(t *testing.T) {
	testCases := []struct {
		name        string
		inputNodeID uint64
		expectedErr error
	}{
		{
			name:        "Successful case",
			inputNodeID: 1,
			expectedErr: nil,
		},
		{
			name:        "Internal server error case",
			inputNodeID: 2,
			expectedErr: common.ErrInternalServerError,
		},
		{
			name:        "Not found case",
			inputNodeID: 3,
			expectedErr: common.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockNodesRepo)
			mockRepo.On("DeleteNode", tc.inputNodeID).Return(tc.expectedErr)

			err := mockRepo.DeleteNode(tc.inputNodeID)
			assert.Equal(t, tc.expectedErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
