package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/handlers"
	"github.com/unification-com/unode-onboard-api/pkg/models"
	mockSvc "github.com/unification-com/unode-onboard-api/pkg/services/mocks"
)

// TODO: refactor
func TestGetNodeByIDHandler(t *testing.T) {
	testCases := []struct {
		name string
		// inputNodeID        uint64
		inputNodeID        string
		mockResponse       models.Nodes
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Successful case",
			// inputNodeID: 1,
			inputNodeID: "1",
			mockError:   nil,
			mockResponse: models.Nodes{
				NodeID: 1, ChainName: "ethereum", Status: "active", PublicIP: "127.0.0.1", RPCPort: "8545", RPCWebSocketPort: "8546", RestHTTPPort: "8547",
			},
			expectedStatusCode: 200,
			expectedResponse: models.Nodes{
				NodeID: 1, ChainName: "ethereum", Status: "active", PublicIP: "127.0.0.1", RPCPort: "8545", RPCWebSocketPort: "8546", RestHTTPPort: "8547",
			},
		},
		{
			name: "Record not found for node id case",
			// inputNodeID:        0,
			inputNodeID:        "0",
			mockError:          common.ErrRecordNotFound,
			mockResponse:       models.Nodes{},
			expectedStatusCode: 404,
			expectedResponse:   fiber.Map{"message": common.ErrRecordNotFound.Error()},
		},
		{
			name: "Internal server error case",
			// inputNodeID:        0,
			inputNodeID:        "0",
			mockError:          common.ErrInternalServerError,
			mockResponse:       models.Nodes{},
			expectedStatusCode: 500,
			expectedResponse:   fiber.Map{"message": common.ErrInternalServerError.Error()},
		},
		{
			name: "Invalid node id case",
			// inputNodeID:        0,
			inputNodeID:        "invalid",
			mockError:          common.ErrInvalidID,
			mockResponse:       models.Nodes{},
			expectedStatusCode: 400,
			expectedResponse:   fiber.Map{"message": common.ErrInvalidID.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockSvc.MockNodesService)
			nodesHandler := handlers.NewNodesHandler(mockSvc)

			nodeID, err := common.StringToUint64(tc.inputNodeID)
			if err != nil {
				assert.Equal(t, tc.expectedStatusCode, 400)
				assert.Equal(t, tc.expectedResponse.(fiber.Map), fiber.Map{"message": err.Error()})
				return
			}
			mockSvc.On("GetNodeByID", nodeID).Return(tc.mockResponse, tc.mockError).Once()

			app := fiber.New()
			app.Get("/node/:id", nodesHandler.GetNodeByIDHandler)

			req := httptest.NewRequest("GET", fmt.Sprintf("/node/%s", tc.inputNodeID), nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while testing: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// t.Logf("Response body: %s", string(body))

			if tc.expectedStatusCode == 200 {
				var respBody models.Nodes
				if err := json.Unmarshal(body, &respBody); err != nil {
					t.Fatalf("Error unmarshaling response body to models.Nodes: %v", err)
				}
				assert.Equal(t, tc.expectedResponse.(models.Nodes), respBody)
			} else {
				var respBody fiber.Map
				if err := json.Unmarshal(body, &respBody); err != nil {
					t.Fatalf("Error unmarshaling response body to map[string]interface{}: %v", err)
				}
				assert.Equal(t, tc.expectedResponse.(fiber.Map), respBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

// TODO: refactor this test. Error: panic: interface conversion: interface {} is nil, not []models.Nodes [recovered]
func TestGetAllNodesByAccountIDHandler(t *testing.T) {
	testCases := []struct {
		name               string
		inputAccountID     uint64
		mockResponse       []models.Nodes
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:           "Successful case",
			inputAccountID: 1,
			mockResponse: []models.Nodes{
				{NodeID: 1, NodeName: "test1", ChainName: "ethereum", Status: "active", PublicIP: "127.0.0.1", RPCPort: "8545", RPCWebSocketPort: "8546", RestHTTPPort: "8547", Location: "India"},
			},
			mockError:          nil,
			expectedStatusCode: 200,
			expectedResponse: []models.Nodes{
				{NodeID: 1, NodeName: "test1", ChainName: "ethereum", Status: "active", PublicIP: "127.0.0.1", RPCPort: "8545", RPCWebSocketPort: "8546", RestHTTPPort: "8547", Location: "India"},
			},
		},
		{
			name:               "Internal server error case",
			inputAccountID:     1,
			mockResponse:       []models.Nodes{},
			mockError:          common.ErrInternalServerError,
			expectedStatusCode: 500,
			expectedResponse:   fiber.Map{"message": common.ErrInternalServerError.Error()},
		},
		{
			name:               "Record not found case",
			inputAccountID:     0,
			mockResponse:       []models.Nodes{},
			mockError:          common.ErrRecordNotFound,
			expectedStatusCode: 404,
			expectedResponse:   fiber.Map{"message": common.ErrRecordNotFound.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockSvc.MockNodesService)
			nodesHandler := handlers.NewNodesHandler(mockSvc)

			mockSvc.On("GetAllNodesByAccountID", tc.inputAccountID).Return(tc.mockResponse, tc.mockError)

			app := fiber.New()

			app.Get("/nodes", func(ctx *fiber.Ctx) error {
				// Set local "accountid" value in Fiber context
				ctx.Locals("accountid", tc.inputAccountID)
				// Call handler function
				return nodesHandler.GetAllNodesByAccountIDHandler(ctx)
			})

			req := httptest.NewRequest("GET", "/nodes", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while testing: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			if tc.expectedStatusCode == 200 {
				var respBody []models.Nodes
				if err := json.Unmarshal(body, &respBody); err != nil {
					t.Fatalf("Error unmarshaling response body to models.Nodes: %v", err)
				}
				assert.Equal(t, tc.expectedResponse.([]models.Nodes), respBody)
			} else {
				var respBody fiber.Map
				if err := json.Unmarshal(body, &respBody); err != nil {
					t.Fatalf("Error unmarshaling response body to map[string]interface{}: %v", err)
				}
				assert.Equal(t, tc.expectedResponse.(fiber.Map), respBody)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

// TODo: refactor bad request case
// func TestAddNodeHandler(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 		// inputBody          models.Nodes
// 		inputBody          interface{}
// 		inputAccountID     uint64
// 		mockError          error
// 		expectedStatusCode int
// 		expectedResponse   fiber.Map
// 	}{
// 		{
// 			name:           "Successful case",
// 			inputAccountID: 1,
// 			inputBody: models.Nodes{
// 				NodeName:  "test1",
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 				Location:  "India",
// 			},
// 			mockError:          nil,
// 			expectedStatusCode: 201,
// 			expectedResponse: fiber.Map{
// 				"message": "Node created successfully"},
// 		},
// 		{
// 			name: "Internal server error case",
// 			inputBody: models.Nodes{
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 			},
// 			mockError:          common.ErrInternalServerError,
// 			expectedStatusCode: 500,
// 			expectedResponse: fiber.Map{
// 				"message": common.ErrInternalServerError.Error()},
// 		},
// 		{
// 			name: "Recrord already exists case",
// 			inputBody: models.Nodes{
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 			},
// 			mockError:          common.ErrDuplicateKey,
// 			expectedStatusCode: 409,
// 			expectedResponse: fiber.Map{
// 				"message": common.ErrDuplicateKey.Error()},
// 		},
// 		// {
// 		// 	name:               "Bad request case",
// 		// 	inputBody:          fiber.Map{},
// 		// 	mockError:          common.ErrBadRequest,
// 		// 	expectedStatusCode: 400,
// 		// 	expectedResponse: fiber.Map{
// 		// 		"message": common.ErrBadRequest.Error()},
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockSvc := new(mockSvc.MockNodesService)
// 			nodesHandler := handlers.NewNodesHandler(mockSvc)

// 			// if tc.expectedStatusCode != 400 {
// 			// 	mockSvc.On("AddNode", tc.inputBody).Return(tc.mockError).Maybe()
// 			// }
// 			mockSvc.On("AddNode", tc.inputBody).Return(tc.mockError).Maybe()

// 			app := fiber.New()
// 			// app.Post("/nodes", nodesHandler.AddNodeHandler)

// 			app.Post("/nodes", func(ctx *fiber.Ctx) error {
// 				// Set local "accountid" value in Fiber context
// 				ctx.Locals("accountid", tc.inputAccountID)
// 				// Call handler function
// 				return nodesHandler.AddNodeHandler(ctx)
// 			})

// 			inputBodyBytes, err := json.Marshal(tc.inputBody)
// 			if err != nil {
// 				t.Fatalf("Error marshaling input body: %v", err)
// 			}
// 			req := httptest.NewRequest("POST", "/nodes", bytes.NewReader(inputBodyBytes))
// 			req.Header.Set("Content-Type", "application/json")

// 			resp, err := app.Test(req)
// 			if err != nil {
// 				t.Fatalf("Error while testing: %v", err)
// 			}
// 			defer resp.Body.Close()
// 			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

// 			// Read the response body
// 			body, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("Error reading response body: %v", err)
// 			}

// 			t.Logf("%s", string(body))

// 			var respBody fiber.Map
// 			if err := json.Unmarshal(body, &respBody); err != nil {
// 				t.Fatalf("Error unmarshaling response body to models.Nodes: %v", err)
// 			}
// 			assert.Equal(t, tc.expectedResponse, respBody)

// 			mockSvc.AssertExpectations(t)

// 			mockSvc.On("AddNode", tc.inputBody).Return(tc.mockError)
// 		})
// 	}
// }

func TestDeleteNodeHandler(t *testing.T) {
	testCases := []struct {
		name               string
		inputNodeID        uint64
		mockError          error
		expectedStatusCode int
		expectedResponse   fiber.Map
	}{
		{
			name:               "Successful case",
			inputNodeID:        1,
			mockError:          nil,
			expectedStatusCode: 200,
			expectedResponse: fiber.Map{
				"message": "Node deleted successfully",
			},
		},
		{
			name:               "Internal server error case",
			inputNodeID:        1,
			mockError:          common.ErrInternalServerError,
			expectedStatusCode: 500,
			expectedResponse: fiber.Map{
				"message": common.ErrInternalServerError.Error(),
			},
		},
		{
			name:               "Not found case",
			inputNodeID:        1,
			mockError:          common.ErrRecordNotFound,
			expectedStatusCode: 404,
			expectedResponse: fiber.Map{
				"message": common.ErrRecordNotFound.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockSvc.MockNodesService)
			nodesHandler := handlers.NewNodesHandler(mockSvc)

			mockSvc.On("DeleteNode", tc.inputNodeID).Return(tc.mockError)

			app := fiber.New()
			app.Delete("/nodes/:id", nodesHandler.DeleteNodeHandler)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/nodes/%d", tc.inputNodeID), nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error while testing: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			t.Logf("%s", string(body))

			var respBody fiber.Map
			if err := json.Unmarshal(body, &respBody); err != nil {
				t.Fatalf("Error unmarshaling response body to models.Nodes: %v", err)
			}
			assert.Equal(t, tc.expectedResponse, respBody)

			mockSvc.AssertExpectations(t)
		},
		)
	}
}

// TODO add bad request case
// func TestUpdateNodeHandler(t *testing.T) {

// 	testCases := []struct {
// 		name               string
// 		inputBody          interface{}
// 		mockError          error
// 		expectedStatusCode int
// 		expectedResponse   fiber.Map
// 	}{
// 		{
// 			name: "Successful case",
// 			inputBody: models.Nodes{
// 				NodeID:    1,
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 			},
// 			mockError:          nil,
// 			expectedStatusCode: 200,
// 			expectedResponse: fiber.Map{
// 				"message": "Updated successfully",
// 			},
// 		},
// 		{
// 			name: "Internal server error case",
// 			inputBody: models.Nodes{
// 				NodeID:    1,
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 			},
// 			mockError:          common.ErrInternalServerError,
// 			expectedStatusCode: 500,
// 			expectedResponse: fiber.Map{
// 				"message": common.ErrInternalServerError.Error(),
// 			},
// 		},
// 		{
// 			name: "Record not found case",
// 			inputBody: models.Nodes{
// 				NodeID:    1,
// 				ChainName: "ethereum",
// 				Status:    "active",
// 				PublicIP:  "127.0.0.1",
// 				RPCPort:   "8545",
// 			},
// 			mockError:          common.ErrRecordNotFound,
// 			expectedStatusCode: 404,
// 			expectedResponse: fiber.Map{
// 				"message": common.ErrRecordNotFound.Error(),
// 			},
// 		}}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockSvc := new(mockSvc.MockNodesService)
// 			nodesHandler := handlers.NewNodesHandler(mockSvc)

// 			mockSvc.On("UpdateNode", tc.inputBody).Return(tc.mockError)

// 			app := fiber.New()
// 			app.Put("/nodes", nodesHandler.UpdateNodeHandler)

// 			inputBodyBytes, err := json.Marshal(tc.inputBody)
// 			if err != nil {
// 				t.Fatalf("Error marshaling input body: %v", err)
// 			}
// 			req := httptest.NewRequest("PUT", "/nodes", bytes.NewReader(inputBodyBytes))
// 			req.Header.Set("Content-Type", "application/json")

// 			resp, err := app.Test(req)
// 			if err != nil {
// 				t.Fatalf("Error while testing: %v", err)
// 			}
// 			defer resp.Body.Close()
// 			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

// 			// Read the response body
// 			body, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("Error reading response body: %v", err)
// 			}

// 			t.Logf("%s", string(body))

// 			var respBody fiber.Map
// 			if err := json.Unmarshal(body, &respBody); err != nil {
// 				t.Fatalf("Error unmarshaling response body to models.Nodes: %v", err)
// 			}
// 			assert.Equal(t, tc.expectedResponse, respBody)

// 			mockSvc.AssertExpectations(t)
// 		})
// 	}
// }
