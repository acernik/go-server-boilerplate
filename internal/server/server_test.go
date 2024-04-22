package server_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"go.uber.org/zap"

	"github.com/acernik/go-server-boilerplate/internal/app"
	"github.com/acernik/go-server-boilerplate/internal/server"
	"github.com/acernik/go-server-boilerplate/internal/service"
	"github.com/acernik/go-server-boilerplate/internal/store"
	"github.com/acernik/go-server-boilerplate/internal/testutil"
)

type ServerTestSuite struct {
	suite.Suite
	migrator *testutil.Migrator
	mysql    testcontainers.Container
	ctx      context.Context
	db       *sql.DB
	host     string
	port     string
	cfg      *app.Config
}

func createJsonPostRequest(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func createGetRequest(c *gin.Context) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
}

func newServerTestSuite() *ServerTestSuite {
	return &ServerTestSuite{}
}

func getMySQLContainer(ctx context.Context, cfg *app.Config) (testcontainers.Container, error) {
	mysqlC, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(cfg.Database.Name),
		mysql.WithUsername(cfg.Database.User),
		mysql.WithPassword(cfg.Database.Password),
	)
	if err != nil {
		return nil, err
	}

	return mysqlC, nil
}

func (s *ServerTestSuite) SetupSuite() {
	s.T().Setenv("DATABASE_USER", "root")
	s.T().Setenv("DATABASE_PASSWORD", "password")
	s.T().Setenv("DATABASE_NAME", "norddb")

	cfg, err := app.LoadConfig()
	s.NoError(err)

	ctx := context.Background()
	s.ctx = ctx

	mysqlContainer, err := getMySQLContainer(ctx, cfg)
	s.NoError(err)

	s.mysql = mysqlContainer

	host, err := s.mysql.Endpoint(ctx, "")
	s.NoError(err)

	hostArr := strings.Split(host, ":")
	s.host = hostArr[0]
	s.port = hostArr[1]

	cfg.Database.Host = s.host
	cfg.Database.Port = s.port

	s.cfg = cfg

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		host,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	s.NoError(err)

	s.db = db

	s.migrator = testutil.NewMigrator("../../database/migrations/api", db)
}

func (s *ServerTestSuite) TearDownSuite() {
	err := s.mysql.Terminate(s.ctx)
	s.NoError(err)
}

func (s *ServerTestSuite) SetupTest() {
	err := s.migrator.MigrateUp()
	s.NoError(err)
}

func (s *ServerTestSuite) TearDownTest() {
	err := s.migrator.MigrateDown()
	s.NoError(err)
}

func (s *ServerTestSuite) TestInsert() {
	logger, err := zap.NewProduction()
	s.NoError(err)

	testStore, err := store.NewStore(s.cfg)
	s.NoError(err)

	testService := service.NewService(testStore)

	testServer := &server.Server{
		Service:    testService,
		Logger:     logger,
		TimeFormat: "2006-01-02T15:04:05",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	createJsonPostRequest(
		ctx,
		map[string]interface{}{
			"name":        gofakeit.Name(),
			"description": gofakeit.Comment(),
		},
	)

	testServer.Insert(ctx)
	s.Equal(http.StatusCreated, w.Code)
}

func (s *ServerTestSuite) TestInsertErrorInvalidRequest() {
	logger, err := zap.NewProduction()
	s.NoError(err)

	testStore, err := store.NewStore(s.cfg)
	s.NoError(err)

	testService := service.NewService(testStore)

	testServer := &server.Server{
		Service:    testService,
		Logger:     logger,
		TimeFormat: "2006-01-02T15:04:05",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	createJsonPostRequest(
		ctx,
		map[string]interface{}{},
	)

	testServer.Insert(ctx)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ServerTestSuite) TestHealth() {
	logger, err := zap.NewProduction()
	s.NoError(err)

	testStore, err := store.NewStore(s.cfg)
	s.NoError(err)

	testService := service.NewService(testStore)

	testServer := &server.Server{
		Service:    testService,
		Logger:     logger,
		TimeFormat: "2006-01-02T15:04:05",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	createGetRequest(ctx)

	testServer.Health(ctx)
	s.Equal(http.StatusOK, w.Code)
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, newServerTestSuite())
}
