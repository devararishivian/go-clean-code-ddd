package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	appConfig "github.com/devararishivian/go-clean-code-ddd/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Conn *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
	// Create a new context for the connection
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.Name,
	)

	// Define the database configuration
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Connect to the database
	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create a new database instance
	db := &Database{
		Conn: dbpool,
	}

	// Register a signal handler to gracefully close the database connection when the application shuts down
	go db.gracefulShutdown(ctx)

	return db, nil
}

func (db *Database) gracefulShutdown(ctx context.Context) {
	// Create a channel to listen for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal
	<-c

	log.Println("Shutting down database connection...")

	// Close the database connection
	db.Conn.Close()

	// Cancel the context to release any idle connections held by the pool
	cancelFunc := func() {}
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		ctx, cancelFunc = context.WithTimeout(ctx, timeout)
	}
	cancelFunc()

	log.Println("Database connection closed")

	os.Exit(0)
}
