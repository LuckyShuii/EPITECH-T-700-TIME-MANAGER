package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedIfEmptyUsersDevOnly(
	ctx context.Context,
	pool *pgxpool.Pool,
	fixturesPath string,
	minUsers int64,
	projectStatus string,
	fixturesPassword string,
) error {

	// DEV ONLY
	if strings.ToLower(projectStatus) != "dev" {
		return nil
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Advisory lock (pas de double seed)
	_, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(424242);`)
	if err != nil {
		return fmt.Errorf("advisory lock: %w", err)
	}

	// Check users
	var count int64
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return fmt.Errorf("count users: %w", err)
	}

	if count >= minUsers {
		return tx.Commit(ctx)
	}

	// Lire fixtures
	b, err := os.ReadFile(fixturesPath)
	if err != nil {
		return fmt.Errorf("read fixtures: %w", err)
	}

	// Remplacer le placeholder du mot de passe
	sqlContent := strings.ReplaceAll(string(b), "${FIXTURES_PASSWORD}", fixturesPassword)

	// Exec fixtures
	_, err = tx.Exec(ctx, sqlContent)
	if err != nil {
		return fmt.Errorf("exec fixtures: %w", err)
	}

	return tx.Commit(ctx)
}
