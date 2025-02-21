package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database instance
var db *gorm.DB

// Todo Model (Like @Entity in Spring Boot)
type Todo struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"type:text;not null" json:"title"`
	Done  bool   `json:"done"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (or specify your frontend URL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Initialize database connection
func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=1234 dbname=todos_db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Auto Migrate: Creates the table if it does not exist
	db.AutoMigrate(&Todo{})
	fmt.Println("Connected to the database successfully.")
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos) // SELECT * FROM todos

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Get single todo by ID
func getTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo

	if err := db.First(&todo, params["id"]).Error; err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Create a new todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	json.NewDecoder(r.Body).Decode(&newTodo)

	db.Create(&newTodo) // INSERT INTO todos (title, done) VALUES (...)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// Update a todo
func updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo

	if db.First(&todo, params["id"]).Error != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(&todo)
	db.Save(&todo) // UPDATE todos SET ...

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Delete a todo
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo Todo

	if db.First(&todo, params["id"]).Error != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	db.Delete(&todo) // DELETE FROM todos WHERE id=...
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	initDB() // Connect to database
	// Example usage: Create a new todo
	db.Create(&Todo{Title: "Learn GORM", Done: false})
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", enableCORS(r))
}
