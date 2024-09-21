package main

import (
"database/sql"
"fmt"
"html/template"
"log"
"net/http"

_ "github.com/go-sql-driver/mysql"
)

// DB connection string
const (
DB_USER = "admin"
DB_PASSWORD = "password123ABC***"
DB_HOST = ""
DB_PORT = "3306"
DB_NAME = "db_name"
)

// User struct to hold the form data
type User struct {
Name string
Email string
Age int
}

var tmpl = template.Must(template.ParseFiles("templates/form.html"))
var db *sql.DB

// Initialize database connection
func initDB() {
var err error
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
db, err = sql.Open("mysql", dsn)
if err != nil {
log.Fatal("Error connecting to the database: ", err)
}

// Ping the database to ensure the connection is established
err = db.Ping()
if err != nil {
log.Fatal("Unable to reach the database: ", err)
}
fmt.Println("Connected to the database successfully!")
}

// Handler for displaying the form
func formHandler(w http.ResponseWriter, r *http.Request) {
if r.Method == http.MethodPost {
err := r.ParseForm()
if err != nil {
http.Error(w, "Unable to process form", http.StatusInternalServerError)
return
}

name := r.FormValue("name")
email := r.FormValue("email")
age := r.FormValue("age")

_, err = db.Exec("INSERT INTO users(name, email, age) VALUES(?, ?, ?)", name, email, age)
if err != nil {
http.Error(w, "Unable to save user data", http.StatusInternalServerError)
return
}

fmt.Fprintln(w, "User data saved successfully!")
return
}

tmpl.Execute(w, nil)
}

func main() {
initDB()
defer db.Close()

// Create users table if it doesn't exist
_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
id INT AUTO_INCREMENT PRIMARY KEY,
name VARCHAR(100),
email VARCHAR(100),
age INT
)`)
if err != nil {
log.Fatal("Unable to create table: ", err)
}

http.HandleFunc("/", formHandler)

fmt.Println("Server started at :8080")
log.Fatal(http.ListenAndServe(":8080", nil))
}
