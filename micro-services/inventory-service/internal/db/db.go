package db

import (
    "context"
    "fmt"
	"log"

    "github.com/jackc/pgx/v4"
	
	// internal
	"inventory/internal/config"
)

type DbService struct {
    db *pgx.Conn
}

func NewDbService() (*DbService, error) {
    // Create a new configuration object for the database connection.
	// Open a connection to the Postgres database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s/%s", config.PostgresConfig.Username, config.PostgresConfig.Password, config.PostgresConfig.Host, config.PostgresConfig.Database))
	if err != nil {
		log.Fatal(err)
	}

    // Test the connection to make sure it's working.
    err = conn.Ping(ctx)
    if err != nil {
        return nil, err
    }

    // Return a new DbService with the open connection.
    return &DbService{conn}, nil
}

func (s *DbService) Close() error {
    // Close the connection to the database.
    return s.db.Close(context.Background())
}


func (s *DbService) PutOnHold(productCode string, quantity int, userID string) (int, error) {
	withLock := quantity > 5
	if withLock {
		log.Printf("[debug] more than 5 products is requsted to be put on hold, using lock.")
	}
	return s.hold(productCode, quantity, userID, withLock)
}

func (s *DbService) hold(productCode string, quantity int, userID string, withLock bool) (int, error) {
	if (withLock) {
		// checks out lock from redis.
		log.Printf("[info] making transaction with lock")
	}


	// makes transaction
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
    if err != nil {
        return -1, err
    }

	// adds the hold-entry to the hold table
	_, err = tx.Exec(ctx, "INSERT INTO hold (created_at, user_id, product_code, quantity) VALUES (NOW(), $1, $2, $3)", userID, productCode, quantity)
	if err != nil {
		log.Printf("[warning] failed instering into hold-table: %s", err)
		tx.Rollback(ctx)
		return -1, err
	}

	// get the aviable product count
	var inventoryCount int
	err = tx.QueryRow(ctx, "SELECT available_quantity FROM available_products WHERE product_code = $1", productCode).Scan(&inventoryCount)
	if err != nil {
		log.Printf("[warning] failed getting quantity from available_products: %s", err)
		tx.Rollback(ctx)
		return -1, err
	}

	// checks so new quantity is not below zero
	if (inventoryCount < 0) {
		log.Printf("[info] transaction resultet in negative inventory count")
		tx.Rollback(ctx)
		return -1, fmt.Errorf("transaction resulted in negative inventory count.")
	}

	// different ways to handle if a lock was checkedout or not
	if (withLock) {
		err = tx.Commit(ctx)
		if err != nil {
			log.Printf("[warning] could not commit transction")
			return -1, err
		}
		
		// returns lock to redis
		log.Printf("[TODO] implement returning lock to redis.")
		
	} else {
		if (inventoryCount < 5) {
			// checks if inventory count is low after transaction, in that case same procedure is tested again
			// but this time checking out mutex lock from redis.
			tx.Rollback(ctx)
			log.Printf("[info] transaction resulted in inventory count below 5, re-trying with lock")
			return s.hold(productCode, quantity, userID, true)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				log.Printf("[warning] could not commit transction")
				return -1, err
			}
		}
	}

	log.Printf("[debug] inventory hold was accepted")
	return inventoryCount, nil 
}


