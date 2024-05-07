package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Person struct {
	ID   int
	Name string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	uname := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DB")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", uname, pass, dbname))
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var people []Person

	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		people = append(people, p)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("index").Parse(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>People</title>
        </head>
        <body>
            <h1>People</h1>
            <ul>
                {{range .}}
                    <li>{{.ID}}: {{.Name}}</li>
                {{end}}
            </ul>
        </body>
        </html>
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and send the response
	err = tmpl.Execute(w, people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
