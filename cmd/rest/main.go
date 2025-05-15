package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Order struct {
	ID           int    `json:"id"`
	CustomerName string `json:"customer_name"`
	CreatedAt    string `json:"created_at"`
}

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, customer_name, created_at FROM orders")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(&o.ID, &o.CustomerName, &o.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			orders = append(orders, o)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	log.Println("REST service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
