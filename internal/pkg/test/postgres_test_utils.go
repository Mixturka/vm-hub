package test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type PostgresTestUtil struct {
	DB             *pgxpool.Pool
	MigrationsPath string
	ConnStr        string
}

func NewPostgresTestUtil() *PostgresTestUtil {
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

	connStr := mustGetEnv("TEST_POSTGRES_URL")
	migrationsPath := mustGetEnv("POSTGRES_MIGRATIONS_PATH")
	absoluteMigrationsPath := getAbsolutePath(projRoot, migrationsPath)

	applyMigrations(connStr, absoluteMigrationsPath)

	db := mustConnectDB(connStr)
	return &PostgresTestUtil{
		DB:             db,
		MigrationsPath: absoluteMigrationsPath,
		ConnStr:        connStr,
	}
}

// Tears down the database by rolling back migrations and closing connections.
func (util *PostgresTestUtil) Close() {
	defer util.DB.Close()
	log.Println("Rolling back migrations...")

	migration, err := migrate.New("file://"+util.MigrationsPath, util.ConnStr)
	if err != nil {
		log.Fatalf("Failed to initialize rollback: %v", err)
	}
	if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to rollback migrations: %v", err)
	}
}

// Clears test data from database.
func (utils *PostgresTestUtil) TruncateTables(t *testing.T, tables ...string) {
	for _, table := range tables {
		_, err := utils.DB.Exec(context.Background(), `TRUNCATE `+table+` RESTART IDENTITY CASCADE;`)
		if err != nil {
			t.Fatalf("Failed to truncate table: %s", table)
		}
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in environment", key)
	}
	return value
}

func getAbsolutePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

func mustConnectDB(connStr string) *pgxpool.Pool {
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to create connection pool")
	}
	return db
}

func applyMigrations(connStr, migrationsPath string) {
	migration, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
