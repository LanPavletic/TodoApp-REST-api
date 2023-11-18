package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LanPavletic/go-rest-server/internal/middleware"
	"github.com/LanPavletic/go-rest-server/internal/responds"
	"github.com/LanPavletic/go-rest-server/internal/stores"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Server stores pointers to gin engine, a mongo client connection
// and a task and user store/controllers
type Server struct {
	store         *stores.TaskStore
	userStore     *stores.UserStore
	engine        *gin.Engine
	mongoDBclient *mongo.Client
}

// New creates a new server with default logger and recovery middleware
// The funcction also allocates a new task and user store
func New() *Server {
	ginEngine := gin.New()

	// default middleware for logger and panic recovery
	ginEngine.Use(gin.Logger())
	ginEngine.Use(gin.Recovery())
	mongoDBclient := stores.InitMongoDB()

	return &Server{
		store:         stores.NewTaskStore(mongoDBclient),
		userStore:     stores.NewUserStore(mongoDBclient),
		engine:        ginEngine,
		mongoDBclient: mongoDBclient,
	}
}

// RegisterRoutes registers all routes for the server
// Some routes are protected with AuthRequired middleware
func (s *Server) RegisterRoutes() {
	s.engine.POST("/register", s.HandleRegister)
	s.engine.POST("/login", s.HandleLogin)

	protected := s.engine.Group("/task")

	// task routes require a valid JWT token in the Authorization header
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", s.HandleGetAllTasks)
		protected.POST("/", s.HandleCreateTask)
		protected.GET("/:id", s.HandleGetTask)
		protected.DELETE("/:id", s.HandleDeleteTask)
		protected.PUT("/:id", s.HandleUpdateTask)
	}

}

// Starts the server on port 8080 with TLS
func (s *Server) Start() {
	// read cert and key from repo
	sslPath := os.Getenv("SSL_CERT_PATH")
	err := s.engine.RunTLS(":8080", sslPath+"server.crt", sslPath+"server.key")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Stop() error {
	return s.mongoDBclient.Disconnect(context.Background())
}

// Handler for /task GET request
func (s *Server) HandleGetAllTasks(c *gin.Context) {
	tasks, err := s.store.GetAll()
	if err != nil {
		responds.InternalServerError(c)
	}
	c.JSON(http.StatusOK, tasks)
}

// Handler for /task POST request
func (s *Server) HandleCreateTask(c *gin.Context) {
	// read request body
	var requestedTask stores.Task
	c.BindJSON(&requestedTask)

	// add task to store
	id := s.store.Create(&requestedTask)

	// respond with created tasks id
	c.JSON(http.StatusCreated, struct{ Id interface{} }{id})
}

// Handler for /task/:id GET request
func (s *Server) HandleGetTask(c *gin.Context) {
	// parse id parameter from POST request
	id := ParseId(c)

	fmt.Println(id)
	// get task from store
	task, err := s.store.Get(id)
	if err != nil {
		responds.NotFound(c)
		return
	}
	c.JSON(http.StatusOK, task)
}

// Handler for /task/:id DELETE request
func (s *Server) HandleDeleteTask(c *gin.Context) {
	id := ParseId(c)

	err := s.store.Delete(id)
	if err != nil {
		responds.NotFound(c)
	}
	c.JSON(http.StatusNoContent, nil)
}

// Handler for /task/:id PUT request
func (s *Server) HandleUpdateTask(c *gin.Context) {
	var requestedTask stores.Task
	c.BindJSON(&requestedTask)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {

		responds.BadRequest(c)
	}

	// update task in store
	err = s.store.Update(id, &requestedTask)
	if err != nil {
		responds.NotFound(c)
	}

	// respond with no content
	c.JSON(http.StatusNoContent, struct{}{})
}

// Handler for /register POST request
func (s *Server) HandleRegister(c *gin.Context) {
	var userCredentials stores.UserCredentials
	c.BindJSON(&userCredentials)

	id, err := s.userStore.Create(&userCredentials)
	if err != nil {
		responds.InternalServerError(c)
	}
	c.JSON(http.StatusCreated, struct{ primitive.ObjectID }{id})
}

// Handler for /login POST request
func (s *Server) HandleLogin(c *gin.Context) {
	var userCredentials stores.UserCredentials
	c.BindJSON(&userCredentials)

	token, err := s.userStore.Login(&userCredentials)
	if err != nil {
		responds.Unauthorized(c)
	}
	c.JSON(http.StatusOK, struct{ Token string }{token})
}

// ParseId takes id from request url and parses it to primitive.ObjectID
func ParseId(c *gin.Context) primitive.ObjectID {
	ParsedId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responds.BadRequest(c)
	}
	return ParsedId
}
