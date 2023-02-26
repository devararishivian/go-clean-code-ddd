package infrastructure

import (
	"context"
	"fmt"
	appConfig "github.com/devararishivian/antrekuy/internal/config"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Database struct {
	Conn *pgx.Conn
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
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Connect to the database
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create a new database instance
	db := &Database{
		Conn: conn,
	}

	// Register a signal handler to gracefully close the database connection when the application shuts down
	go db.gracefulShutdown()

	return db, nil
}

func (db *Database) gracefulShutdown() {
	// Create a channel to listen for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal
	<-c

	// Start a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close the database connection
	log.Println(db.Conn.Close(ctx))

	os.Exit(0)
}
