package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/acernik/go-server-boilerplate/internal/app"
	srvmodel "github.com/acernik/go-server-boilerplate/internal/service/model"
	"github.com/acernik/go-server-boilerplate/internal/store"
	"github.com/acernik/go-server-boilerplate/internal/testutil"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type StoreTestSuite struct {
	suite.Suite
	migrator *testutil.Migrator
	mysql    testcontainers.Container
	ctx      context.Context
	db       *sql.DB
	host     string
	port     string
	cfg      *app.Config
}

func newStoreTestSuite() *StoreTestSuite {
	return &StoreTestSuite{}
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

func (s *StoreTestSuite) SetupSuite() {
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

func (s *StoreTestSuite) TearDownSuite() {
	err := s.mysql.Terminate(s.ctx)
	s.NoError(err)
}

func (s *StoreTestSuite) SetupTest() {
	err := s.migrator.MigrateUp()
	s.NoError(err)

}

func (s *StoreTestSuite) TearDownTest() {
	err := s.migrator.MigrateDown()
	s.NoError(err)
}

func (s *StoreTestSuite) TestInsertTestTableItem() {
	str, err := store.NewStore(s.cfg)
	s.NoError(err)

	testTableItem := srvmodel.TestTableItem{
		Name:        gofakeit.Name(),
		Description: gofakeit.Comment(),
	}

	_, err = str.InsertTestTableItem(testTableItem)
	s.NoError(err)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, newStoreTestSuite())
}
