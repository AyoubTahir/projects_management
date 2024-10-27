package container

import (
	"database/sql"
	"fmt"

	"github.com/AyoubTahir/projects_management/config"
	"github.com/AyoubTahir/projects_management/internal/handlers"
	"github.com/AyoubTahir/projects_management/internal/repositories"
	"github.com/AyoubTahir/projects_management/internal/services"
	"github.com/AyoubTahir/projects_management/pkg/database"
	"github.com/AyoubTahir/projects_management/pkg/logger"
	"github.com/AyoubTahir/projects_management/pkg/orm"
)

type Container struct {
	config     *config.Config
	db         *sql.DB
	logger     *logger.Logger
	orm        *orm.Orm
	repository *repositories.Repository
	service    *services.Service
	Handler    *handlers.Handler
}

func New(cfg *config.Config) (*Container, error) {
	c := &Container{
		config: cfg,
	}

	if err := c.initLogger(); err != nil {
		return nil, err
	}

	if err := c.initDB(); err != nil {
		return nil, err
	}

	c.initORM()
	c.initRepository()
	c.initService()
	c.initHandler()

	defer c.orm.Cleanup()

	return c, nil
}

func (c *Container) initLogger() error {
	logger, err := logger.New(c.config.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	c.logger = logger
	return nil
}

func (c *Container) initDB() error {
	db, err := database.NewConnection(c.config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.db = db
	return nil
}

func (c *Container) Close() error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

func (c *Container) initORM() error {
	c.orm = orm.New(c.db, orm.Config(c.config.OrmConfig))
	return nil
}

func (c *Container) initRepository() error {
	c.repository = repositories.NewRepository(c.orm)
	return nil
}

func (c *Container) initService() error {
	c.service = services.NewService(c.repository)
	return nil
}

func (c *Container) initHandler() error {
	c.Handler = handlers.NewHandler(c.service)
	return nil
}

// Getters for dependencies
func (c *Container) Logger() *logger.Logger { return c.logger }
