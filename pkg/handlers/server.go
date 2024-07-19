package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/unification-com/unode-onboard-api/pkg/middleware/cors"
)

// type Server struct {
// 	DBClient     db.Client
// 	NodesHandler *NodesHandler
// }

// func NewServer(dbClient db.Client) *Server {

// 	nodesRepo := repositories.NewPgNodesRepository(dbClient.DB)
// 	nodesService := services.NewNodesService(nodesRepo)
// 	nodesHandler := NewNodesHandler(nodesService)

// 	return &Server{
// 		DBClient:     dbClient,
// 		NodesHandler: nodesHandler,
// 	}
// }

type Server struct {
}

func NewServer() *Server {

	return &Server{}
}

func (s *Server) setupRoutes(app *fiber.App) {
	app.Get("/", fiber.Handler(func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	}))

	app.Use(cors.CORSMiddleware)
	// app.Use(auth.Middleware(&s.DBClient))

	// app.Get("/node/:id", s.NodesHandler.GetNodeByIDHandler)
	// app.Get("/node", s.NodesHandler.GetAllNodesByAccountIDHandler)
	// app.Post("/node", s.NodesHandler.AddNodeHandler)
	// app.Put("/node/:id", s.NodesHandler.UpdateNodeHandler)
	// app.Delete("/node/:id", s.NodesHandler.DeleteNodeHandler)
}

func (s *Server) RunServer() error {
	// if err := s.DBClient.InitializeDB(); err != nil {
	// 	return fmt.Errorf("error initialising DB %s", err.Error())
	// }

	app := fiber.New()
	s.setupRoutes(app)
	return app.Listen(":8000")
}
