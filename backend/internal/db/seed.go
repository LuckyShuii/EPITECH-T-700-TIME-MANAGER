package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func CreateRootUserIfNotExists(
	ctx context.Context,
	pool *pgxpool.Pool,
	rootUsername string,
	rootPassword string,
) error {
	if rootUsername == "" || rootPassword == "" {
		return fmt.Errorf("ROOT_USERNAME and ROOT_PASSWORD must be set")
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	// Vérifier si l'utilisateur root existe déjà
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", rootUsername).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check root user: %w", err)
	}

	if exists {
		return nil
	}

	// Hasher le mot de passe
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rootPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Créer l'utilisateur root avec les rôles admin et employee
	userUUID := uuid.New().String()
	_, err = conn.Exec(ctx, `
		INSERT INTO users (
			uuid, username, email, password_hash, 
			first_name, last_name, roles, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, userUUID, rootUsername, rootUsername+"@root.local", string(passwordHash),
		"Root", "Administrator", []string{"employee", "admin"}, "active")

	if err != nil {
		return fmt.Errorf("create root user: %w", err)
	}

	return nil
}
