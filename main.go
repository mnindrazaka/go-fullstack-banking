package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Id      int
	Code    string
	Balance int
}

func main() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dataSourceName := fmt.Sprintf("%s:%s@/%s", dbUsername, dbPassword, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM account")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var accounts []Account

		for rows.Next() {
			account := Account{}
			if err := rows.Scan(&account.Id, &account.Code, &account.Balance); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			accounts = append(accounts, account)
		}

		funcMap := template.FuncMap{
			"sum": func(a int, b int) int {
				return a + b
			},
		}

		tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("views/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := map[string]interface{}{
			"accounts": accounts,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("views/create.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if err := tmpl.Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			code := r.Form.Get("code")
			balance := r.Form.Get("balance")

			if _, err := db.Exec("INSERT INTO account SET code=?, balance=?", code, balance); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/withdrawal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("views/withdrawal.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			data := map[string]interface{}{
				"id": r.URL.Query().Get("id"),
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			amount := r.FormValue("amount")
			id := r.URL.Query().Get("id")

			if _, err := db.Exec("UPDATE account SET balance=balance-? WHERE id=?", amount, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/deposit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("views/deposit.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			data := map[string]interface{}{
				"id": r.URL.Query().Get("id"),
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			amount := r.FormValue("amount")
			id := r.URL.Query().Get("id")

			if _, err := db.Exec("UPDATE account SET balance=balance+? WHERE id=?", amount, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			id := r.URL.Query().Get("id")
			if _, err := db.Exec("DELETE FROM account WHERE id=?", id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
