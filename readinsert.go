package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func insert() {
	// create the new connection of database
	connStr := "user=nilesh dbname=nileshdb password= nilesh host=localhost sslmode=disable "
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	// Create the new context and begin the transction

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 'tx' is an instance  of '*sql.Tx' through which we can execute our queries

	// Here the queury is executed on the transaction instance , and not applied to database yet
	_, err = tx.ExecContext(ctx, "INSERT INTO pets(name,species) VALUES('fido','dog'),('albert','cat')")

	if err != nil {
		// incase we find the  any error in the query execution,rollback the transction
		tx.Rollback()
		return
	}

	// the next query is handeled similarly
	_, err = tx.ExecContext(ctx, "INSERT INTO food(name,quantity) VALUES ('Dog Biscuit',3),('cat Food',5)")
	if err != nil {
		tx.Rollback()
		return
	}

	// Finally, if no errors are recieved from the queires,commit the transaction
	// this applies the above changes to our database
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

func readAndUpdate() {
	connStr := "user=nilesh dbname=nileshdb password= nilesh host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Run the query to get a count of all cats
	row := tx.QueryRow("SELECT count(*) FROM pets WHERE species='cat'")
	var catCount int

	// store the count in the 'catCount' varibale

	err = row.Scan(&catCount)
	if err != nil {
		tx.Rollback()
		return
	}

	// Now update the food table increasing the quantity of cat food by 10x number of cats
	_, err = tx.ExecContext(ctx, "UPDATE food SET quantity=quantity +$1 WHERE name='Cat Food '", 10*catCount)
	if err != nil {
		tx.Rollback()
		return

	}
	// commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
