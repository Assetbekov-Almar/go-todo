package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"package/todohandlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectToSql() (*sql.DB, error) {
	connStr := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
    }

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Set PORT in .env file")
	}

	db, sqlErr := connectToSql()
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}
	defer db.Close()
	r := mux.NewRouter()

	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"}),
		handlers.ExposedHeaders([]string{"Link"}),
	))

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler((http.StripPrefix("/static/", fs)))

	todoRouter := r.PathPrefix("/todo").Subrouter()

	handler := &todohandlers.TodoHandler{ 
		DB: db,
	}

	todoRouter.HandleFunc("/all", handler.ReadHandler).Methods("GET")
	todoRouter.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		// Pre-flight request. Reply successfully:
		if r.Method == http.MethodOptions {
			return
		}
	}).Methods(http.MethodOptions)
	todoRouter.HandleFunc("/add", handler.CreateHandler).Methods("POST")
	todoRouter.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Pre-flight request. Reply successfully:
		if r.Method == http.MethodOptions {
			return
		}
	}).Methods(http.MethodOptions)
	todoRouter.HandleFunc("/delete/{id}", handler.DeleteHandler).Methods("DELETE")
	todoRouter.HandleFunc("/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Pre-flight request. Reply successfully:
		if r.Method == http.MethodOptions {
			return
		}
	}).Methods(http.MethodOptions)
	todoRouter.HandleFunc("/update/{id}", handler.UpdateHandler).Methods("PATCH")

	srv := &http.Server{
		Handler: r,
		Addr: 	 ":" + port,
	}
	log.Printf("Server starting on port %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error %v", err)
	}
}