package middleware

import (
	"database/sql"
	"encoding/json" 
	"api/models" 
	"log"
	"net/http" 
	

	"github.com/gorilla/mux" 
	_ "github.com/lib/pq"      
)


func Getstocks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM stocks")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		stocks := []models.Stock{}
		for rows.Next() {
			var u models.Stock
			if err := rows.Scan(&u.StockID, &u.Name, &u.Company, &u.Price); err != nil {
				log.Fatal(err)
			}
			stocks = append(stocks, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(stocks)
	}
}

func GetStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.Stock
		err := db.QueryRow("SELECT * FROM stocks WHERE id = $1", id).Scan(&u.StockID, &u.Name, &u.Company, &u.Price)
		
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

func CreateStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.Stock
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO stocks (name, company, price) VALUES ($1, $2, $3) RETURNING id", u.Name, u.Company, u.Price).Scan(&u.StockID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

func UpdateStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.Stock
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE stocks SET name = $1, company = $2, price = $3 WHERE id = $4", u.Name, u.Company, u.Price, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

func DeleteStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.Stock
		err := db.QueryRow("SELECT * FROM stocks WHERE id = $1", id).Scan(&u.StockID, &u.Name, &u.Company, &u.Price)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM stocks WHERE id = $1", id)
			if err != nil {
				log.Fatalf("Unable to execute the query. %v", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
	
			json.NewEncoder(w).Encode("Stock deleted")
		}
	}
}