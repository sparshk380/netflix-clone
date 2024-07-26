package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/handlers"
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

func TestIfStringIDExists(t *testing.T) {
	tests := []struct {
		name     string
		chainID  string
		expected bool
	}{
		{
			name:     "Empty chainID",
			chainID:  "",
			expected: false,
		},
		{
			name:     "Non-empty chainID",
			chainID:  "123",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := handlers.IfStringIDExists(tt.chainID)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestValidateNode(t *testing.T) {

	testCases := []struct {
		name          string
		inputNode     models.Nodes
		expecetdError error
	}{
		{
			name: "Successful case",
			inputNode: models.Nodes{
				NodeID: 1,
			},
			expecetdError: nil,
		},
		{
			name: "Node id 0 case",
			inputNode: models.Nodes{
				NodeID: 0,
			},
			expecetdError: common.ErrInvalidID,
		},
		{
			name:          "Node id empty case",
			inputNode:     models.Nodes{},
			expecetdError: common.ErrInvalidID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := handlers.ValidateNode(tc.inputNode)
			t.Logf("Error: %v", err)
			assert.Equal(t, tc.expecetdError, err)
		})
	}

}
