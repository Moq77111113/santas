package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/ent"

	// ent runtime
	_ "github.com/moq77111113/chmoly-santas/ent/runtime"
)

type (
	Container struct {
		Validator    *Validator
		Config       *config.Config
		Web          *echo.Echo
		DB           *sql.DB
		ORM          *ent.Client
		Repositories *Repositories
	}

	Repositories struct {
		Group     *GroupRepo
		Exclusion *ExclusionRepo
	}
)

func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initDb()
	c.initORM()
	c.initRepos()

	return c
}

func (c *Container) Shutdown() error {
	if err := c.ORM.Close(); err != nil {
		return err
	}

	if err := c.DB.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() {
	conf, err := config.GetConfig()

	if err != nil {
		panic(fmt.Sprintf("unable to load config: %v", err))
	}
	c.Config = &conf

	slog.SetLogLoggerLevel(slog.LevelDebug) // TODO: move to config
}

func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

func (c *Container) initWeb() {
	c.Web = echo.New()
	c.Web.Validator = c.Validator
	c.Web.HideBanner = true
}

func (c *Container) initDb() {
	var err error

	connection := c.Config.Database.Connection
	c.DB, err = sql.Open(c.Config.Database.Driver, connection)

	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err))
	}
}

func (c *Container) initORM() {
	driver := entsql.OpenDB(c.Config.Database.Driver, c.DB)
	c.ORM = ent.NewClient(ent.Driver(driver))

	if err := c.ORM.Schema.Create(context.Background()); err != nil {
		panic(err)
	}
}

func (c *Container) initRepos() {
	c.Repositories = new(Repositories)
	c.Repositories.Group = NewGroupRepo(c.ORM)
	c.Repositories.Exclusion = NewExclusionRepo(c.ORM)
}
