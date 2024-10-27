package container

import (
	"database/sql"
	"log"
	"reflect"
)

// Common interfaces
type Repository interface {
	SetDB(db *sql.DB)
}

type Service interface {
	InjectRepositories(container *Container)
}

type Handler interface {
	InjectServices(container *Container)
}

// Container holds all dependencies
type Container struct {
	db           *sql.DB
	repositories map[string]Repository
	services     map[string]Service
	handlers     map[string]Handler
}

func NewContainer(db *sql.DB) *Container {
	return &Container{
		db:           db,
		repositories: make(map[string]Repository),
		services:     make(map[string]Service),
		handlers:     make(map[string]Handler),
	}
}

// RegisterRepository automatically registers a repository
func (c *Container) RegisterRepository(repo Repository) {
	name := reflect.TypeOf(repo).Elem().Name()
	repo.SetDB(c.db)
	c.repositories[name] = repo
}

// RegisterService automatically registers a service
func (c *Container) RegisterService(service Service) {
	name := reflect.TypeOf(service).Elem().Name()
	service.InjectRepositories(c)
	c.services[name] = service
}

// RegisterHandler automatically registers a handler
func (c *Container) RegisterHandler(handler Handler) {
	name := reflect.TypeOf(handler).Elem().Name()
	handler.InjectServices(c)
	c.handlers[name] = handler
}

// Example repository implementation
type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) SetDB(db *sql.DB) {
	r.db = db
}

// Example service implementation
type UserService struct {
	userRepo *UserRepository
}

func (s *UserService) InjectRepositories(c *Container) {
	s.userRepo = c.repositories["UserRepository"].(*UserRepository)
}

// Example handler implementation
type UserHandler struct {
	userService *UserService
}

func (h *UserHandler) InjectServices(c *Container) {
	h.userService = c.services["UserService"].(*UserService)
}

// Helper function to automatically register multiple repositories
func (c *Container) RegisterRepositories(repos ...Repository) {
	for _, repo := range repos {
		c.RegisterRepository(repo)
	}
}

// Helper function to automatically register multiple services
func (c *Container) RegisterServices(services ...Service) {
	for _, service := range services {
		c.RegisterService(service)
	}
}

// Helper function to automatically register multiple handlers
func (c *Container) RegisterHandlers(handlers ...Handler) {
	for _, handler := range handlers {
		c.RegisterHandler(handler)
	}
}

func main() {
	// Initialize database
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Create container
	container := NewContainer(db)

	// Register all repositories at once
	container.RegisterRepositories(
		&UserRepository{},
		&OrderRepository{},
		&ProductRepository{},
		&CategoryRepository{},
		// Add as many repositories as you need...
	)

	// Register all services at once
	container.RegisterServices(
		&UserService{},
		&OrderService{},
		&ProductService{},
		&CategoryService{},
		// Add as many services as you need...
	)

	// Register all handlers at once
	container.RegisterHandlers(
		&UserHandler{},
		&OrderHandler{},
		&ProductHandler{},
		&CategoryHandler{},
		// Add as many handlers as you need...
	)
}

// Additional example implementations
type OrderRepository struct {
	db *sql.DB
}

func (r *OrderRepository) SetDB(db *sql.DB) {
	r.db = db
}

type OrderService struct {
	orderRepo *OrderRepository
	userRepo  *UserRepository
}

func (s *OrderService) InjectRepositories(c *Container) {
	s.orderRepo = c.repositories["OrderRepository"].(*OrderRepository)
	s.userRepo = c.repositories["UserRepository"].(*UserRepository)
}

type OrderHandler struct {
	orderService *OrderService
}

func (h *OrderHandler) InjectServices(c *Container) {
	h.orderService = c.services["OrderService"].(*OrderService)
}
