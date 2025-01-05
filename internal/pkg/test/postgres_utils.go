package test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type PostgresTestUtil struct {
	t       TestingT
	db      *pgxpool.Pool
	connStr string
	refDb   string
	poolMu  sync.Mutex
}

func NewPostgresTestUtilWithIsolatedDb(t TestingT) *PostgresTestUtil {
	return newPostgresTestUtil(t).cloneFromReference()
}

func NewPostgresTestUtilWithIsolatedSchema(t TestingT) *PostgresTestUtil {
	return newPostgresTestUtil(t).createSchema(t)
}

func NewPostgresTestUtilWithTransactionalCleanUp(t TestingT) pgx.Tx {
	const databaseName = "transaction"
	ptutil := newPostgresTestUtil(t)
	ptutil = ptutil.replaceDBName(databaseName)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	tx, err := ptutil.DB().BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, tx.Rollback(ctx))
	})

	return tx
}

func newPostgresTestUtil(t TestingT) *PostgresTestUtil {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	projRoot, err := putils.GetProjectRoot(cwd)
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}

	if err := godotenv.Load(projRoot + "/.env.test"); err != nil {
		log.Fatalf("Error loading environment file: %v", err)
	}

	connStr := MustGetEnv("TEST_POSTGRES_URL")
	migrationsPath := MustGetEnv("POSTGRES_MIGRATIONS_PATH")
	absoluteMigrationsPath := GetAbsolutePath(projRoot, migrationsPath)

	ApplyMigrations(connStr, absoluteMigrationsPath)
	return &PostgresTestUtil{
		t:       t,
		connStr: connStr,
		refDb:   "reference",
	}
}

func (p *PostgresTestUtil) DB() *pgxpool.Pool {
	p.poolMu.Lock()
	defer p.poolMu.Unlock()

	if p.db == nil {
		config, err := pgxpool.ParseConfig(p.connStr)
		require.NoError(p.t, err)

		pool, err := pgxpool.ConnectConfig(context.Background(), config)
		require.NoError(p.t, err)

		p.t.Cleanup(func() {
			pool.Close()
		})

		p.db = pool
	}

	return p.db
}

func (p *PostgresTestUtil) createSchema(t TestingT) *PostgresTestUtil {
	schemaName := newUniqueHumanReadableDatabaseName(p.t)
	schemaName = strings.ToLower(schemaName)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Create the schema
	{
		query := fmt.Sprintf(`CREATE SCHEMA "%s";`, schemaName)
		_, err := p.DB().Exec(ctx, query)
		require.NoError(t, err)
	}

	// Clean up the schema after the test
	t.Cleanup(func() {
		query := fmt.Sprintf(`DROP SCHEMA "%s" CASCADE;`, schemaName)
		_, err := p.DB().Exec(ctx, query)
		require.NoError(t, err)
	})

	// Set search path to the newly created schema
	pgurl := setSearchPath(t, p.connStr, schemaName)

	// Load migrations

	// Apply migrations to the new schema
	// ApplyMigrations(pgurl.String(), absoluteMigrationsPath)

	// Return a new PostgresTestUtil instance with the updated connection string
	return &PostgresTestUtil{
		t:       p.t,
		refDb:   p.refDb,
		connStr: pgurl.String(),
	}
}

func (p *PostgresTestUtil) cloneFromReference() *PostgresTestUtil {
	newDBName := newUniqueHumanReadableDatabaseName(p.t)
	p.t.Logf("Database name for this test: %s", newDBName)

	ctx := context.Background()
	query := fmt.Sprintf(`CREATE DATABASE %q WITH TEMPLATE %q;`, newDBName, p.refDb)
	_, err := p.DB().Exec(ctx, query)
	require.NoError(p.t, err)

	p.t.Cleanup(func() {
		query := fmt.Sprintf(`DROP DATABASE %q WITH {FORCE};`, newDBName)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err := p.DB().Exec(ctx, query)
		require.NoError(p.t, err)
	})

	return p.replaceDBName(newDBName)
}

func (p *PostgresTestUtil) replaceDBName(newDBName string) *PostgresTestUtil {
	o := p.clone()
	o.connStr = replaceDBName(p.t, o.connStr, newDBName)
	return o
}

func (p *PostgresTestUtil) clone() *PostgresTestUtil {
	return &PostgresTestUtil{
		t:       p.t,
		connStr: p.connStr,
		refDb:   p.refDb,
	}
}

func replaceDBName(t TestingT, dataSourceURL, dbname string) string {
	r, err := url.Parse(dataSourceURL)
	require.NoError(t, err)
	r.Path = dbname
	return r.String()
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in environment", key)
	}
	return value
}

func GetAbsolutePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

func ApplyMigrations(connStr, migrationsPath string) {
	migration, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}

func newUniqueHumanReadableDatabaseName(t TestingT) string {
	output := strings.Builder{}

	const maxIdentifierLengthBytes = 63
	uid := genUnique8BytesID(t)
	maxHumanReadableLenBytes := maxIdentifierLengthBytes - len(uid)

	lastSymbolIsHyphen := false
	for _, r := range t.Name() {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			output.WriteRune(r)
			lastSymbolIsHyphen = false
		} else {
			if !lastSymbolIsHyphen {
				output.WriteRune('-')
			}
			lastSymbolIsHyphen = true
		}
		if output.Len() >= maxHumanReadableLenBytes {
			break
		}
	}
	output.WriteString(uid)
	return output.String()
}

func genUnique8BytesID(t TestingT) string {
	bs := make([]byte, 6)
	_, err := rand.Read(bs)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(bs)
}

func setSearchPath(t TestingT, pgURL string, schemaName string) *url.URL {
	parsedURL, err := url.Parse(pgURL)
	require.NoError(t, err)
	query := parsedURL.Query()
	query.Set("search_path", schemaName)
	parsedURL.RawQuery = query.Encode()
	return parsedURL
}
