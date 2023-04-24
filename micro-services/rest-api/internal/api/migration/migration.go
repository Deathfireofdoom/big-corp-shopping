package migration

import (
    "encoding/csv"
    "fmt"
    "os"
    "log"
    "context"
    "bufio"
    "strconv"
    "strings"

    "github.com/jackc/pgx/v4"
	"big-corp-shopping/rest-api/internal/config"
)

func initializeInventory() {
	// Open a connection to the Postgres database
    ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s/%s", config.PostgresConfig.Username, config.PostgresConfig.Password, config.PostgresConfig.Host, config.PostgresConfig.Database))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	// Create the "inventory" table
	createTable := `
            DROP TABLE IF EXISTS inventory CASCADE;
            CREATE TABLE inventory (
            id SERIAL PRIMARY KEY,
            product_code VARCHAR(100) NOT NULL,
            name VARCHAR(100) NOT NULL,
            quantity INTEGER NOT NULL
        );
    `
    _, err = conn.Exec(ctx, createTable)
    if err != nil {
        log.Fatal(err)
    }

	// Create the "hold" table
	createTable = `
        DROP TABLE IF EXISTS hold CASCADE;
        CREATE TABLE hold (
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		user_id VARCHAR(50) NOT NULL,
		product_code VARCHAR(100) NOT NULL,
		quantity INTEGER NOT NULL
		);
    `
    _, err = conn.Exec(ctx, createTable)
    if err != nil {
        log.Fatal(err)
    }

	// Create the "aviable-products" view
	createView := `
            DROP VIEW IF EXISTS available_products CASCADE;
			CREATE VIEW available_products AS
			SELECT i.product_code, i.name, i.quantity - COALESCE(SUM(h.quantity), 0) AS available_quantity
			FROM inventory i
			LEFT JOIN hold h ON i.product_code = h.product_code
			GROUP BY i.product_code, i.name, i.quantity;
    `
    _, err = conn.Exec(ctx, createView)
    if err != nil {
        log.Fatal(err)
    }

	// loads the data into the inventory file
    // Open the CSV file
    file, err := os.Open("data/inventory.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(bufio.NewReader(file))

    for {
        record, err := reader.Read()
        if err != nil {
            break
        }

        productCode := record[0]
        name := record[1]
        quantity, _ := strconv.Atoi(strings.TrimSpace(record[2]))

        // Insert each row into the "inventory" table
        insertStatement := fmt.Sprintf(`
            INSERT INTO inventory (product_code, name, quantity)
            VALUES ('%s', '%s', %d);
        `, productCode, name, quantity)

        _, err = conn.Exec(ctx, insertStatement)
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("Inventory data loaded successfully")
}
