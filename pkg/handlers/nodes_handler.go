package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/unification-com/unode-onboard-api/pkg/common"
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

type NodesHandler struct {
	service NodesServiceInterface
}

func NewNodesHandler(service NodesServiceInterface) *NodesHandler {
	return &NodesHandler{
		service: service,
	}
}

// TODO: add this api to admin
func (n *NodesHandler) GetAllNodesByAccountIDHandler(c *fiber.Ctx) error {
	accountIDInterface := c.Locals("accountid")
	accountID, err := common.ConvertInterfaceToUint64(accountIDInterface)
	if err != nil {
		return common.HandleError(c, err)
	}

	nodes, err := n.service.GetAllNodesByAccountID(accountID)
	if err != nil {
		logrus.Errorf("failed to GetAllNodesByAccountID in GetAllNodesByAccountIDHandler. Err: %v", err)
		return common.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(nodes)
}

func (n *NodesHandler) AddNodeHandler(c *fiber.Ctx) error {
	var node models.Nodes
	accountIDInterface := c.Locals("accountid")
	accountID, err := common.ConvertInterfaceToUint64(accountIDInterface)
	if err != nil {
		return common.HandleError(c, common.ErrBadRequest)
	}

	if err := c.BodyParser(&node); err != nil {
		logrus.Errorf("failed to parse AddNodeHandler request body. Err: %v", err)
		return common.HandleError(c, common.ErrBadRequest)
	}

	node.AccountID = accountID

	if err := n.service.AddNode(node); err != nil {
		logrus.Errorf("failed to AddNode in AddNodeHandler. Err: %v", err)
		return common.HandleError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Node created successfully"})
}

func (n *NodesHandler) GetNodeByIDHandler(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	if !IfStringIDExists(nodeID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Node ID is required",
		})
	}

	nodeIDInt, err := common.StringToUint64(nodeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	node, err := n.service.GetNodeByID(nodeIDInt)
	if err != nil {
		return common.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(node)
}

func (n *NodesHandler) UpdateNodeHandler(c *fiber.Ctx) error {
	var node models.Nodes
	nodeID := c.Params("id")
	nodeIDInt, err := common.StringToUint64(nodeID)
	if err != nil {
		return common.HandleError(c, err)
	}
	if err := c.BodyParser(&node); err != nil {
		logrus.Errorln("failed to parse UpdateNode request body. Err: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	node.NodeID = nodeIDInt

	if err := ValidateNode(node); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	if err := n.service.UpdateNode(node); err != nil {
		return common.HandleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Updated successfully",
	})
}

func (n *NodesHandler) DeleteNodeHandler(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	if !IfStringIDExists(nodeID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Node ID is required",
		})
	}

	nodeIDInt, err := common.StringToUint64(nodeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := n.service.DeleteNode(nodeIDInt); err != nil {
		return common.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Node deleted successfully",
	})
}
