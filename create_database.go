package main

import (
	"database/sql"
	"fmt"
	"strings"
)

func execCmdCreateDB() error {
	var err error

	if err = createDB(); err != nil {
		return err
	}

	if err = emptyDB(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	return nil
}

func createDB() error {
	dbinfo := fmt.Sprintf("user=%s password=%s sslmode=disable", DB_USER, DB_PASSWORD)

	database, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}

	_, err = database.Query("CREATE DATABASE " + DB_NAME)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Printf("database %s already exists\n", DB_NAME)
		}
		return err
	} else {
		fmt.Printf("database %s is created\n", DB_NAME)
	}
	return database.Close()
}

func emptyDB() error {
	_, err := db.Query("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	return err
}

func createTables() error {
	fmt.Println("createTables")
	query := `
		CREATE TABLE customers (
			id SERIAL,
			customer_id character varying(255),
			created_at timestamp without time zone
		);
	`
	_, err := db.Query(query)
	return err
}
