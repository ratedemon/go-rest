package dbtesting

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	maindb "github.com/ratedemon/go-rest/datastore/db"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	user     = "test"
	password = "test_password"
	db       = "test"
	port     = "5432"
	dialect  = "postgres"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	idleConn = 25
	maxConn  = 25
)

func Inject(f func(*testing.T, *maindb.DB)) func(*testing.T) {
	return func(t *testing.T) {
		inject(t, func(db *maindb.DB) {
			f(t, db)
		})
	}
}

func inject(t testing.TB, f func(*maindb.DB)) {
	r := require.New(t)

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		t.Fatalf("Could not start resource: %s", err.Error())
	}
	defer resource.Close()

	var gormDB *gorm.DB
	if err = pool.Retry(func() error {
		gormDB, err = gorm.Open(postgres.New(
			postgres.Config{
				DSN: fmt.Sprintf(dsn, user, password, port, db),
			},
		), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		return err
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err.Error())
	}
	r.NoError(err)

	path := filepath.Join("../../..", "datastore", "postgres", "scripts/init.sql")
	file, err := ioutil.ReadFile(path)
	r.NoError(err)

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_ = gormDB.Exec(request, nil)
	}

	DB := maindb.NewDB(gormDB)

	f(DB)

	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}
