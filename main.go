package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "stripe_importer"
)

var db *sql.DB

var envStripeKey string

func main() {
	envStripeKey = os.Getenv("SI_STRIPE_SECRET")
	var err error
	db, err = openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Please provide a subcommand")
		os.Exit(1)
	}
	subCommand := os.Args[1]

	switch subCommand {
	case "create_db":
		err = execCmdCreateDB()
	case "import_stripe":
		err = execCmdImportStripe()
	default:
		flag.PrintDefaults()
		fmt.Printf("Subcommand %s is not in the list", subCommand)
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}

	return
}

func execCmdImportStripe() error {
	stripe.Key = envStripeKey

	customers, err := fetchCustomers()

	if err != nil {
		return err
	}

	fmt.Printf("Number of Customers: %d\n", len(customers))

	return saveCustomers(customers)
}

func fetchCustomers() ([]*stripe.Customer, error) {
	var customers []*stripe.Customer

	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("limit", "", "20")

	i := customer.List(params)
	for i.Next() {
		c := i.Customer()
		customers = append(customers, c)
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func saveCustomers(customers []*stripe.Customer) error {
	fmt.Println("saving customers")

	for _, c := range customers {
		parsedTime := time.Unix(c.Created, 0)
		_, err := db.Exec("INSERT INTO customers(customer_id, created_at) VALUES($1, $2)", c.ID, parsedTime)
		if err != nil {
			return err
		}
	}
	fmt.Println("done saving customers")
	return nil
}

func openDB() (*sql.DB, error) {
	fmt.Println("opening database connection")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	return sql.Open("postgres", dbinfo)
}
