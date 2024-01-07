package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"package/todohandlers"

	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
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

var mySigningKey = []byte("secret")

func generateRandomKey(length int) string {
    key := make([]byte, length)
    _, err := rand.Read(key)
    if err != nil {
        log.Fatalf("Failed to generate random key: %v", err)
    }
    return base64.StdEncoding.EncodeToString(key)
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "username"
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			fmt.Fprintf(w, "Not Authorized")
		   return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Fprintln(w, err)
		}

		if token.Valid {
			endpoint(w, r)
		}
	}
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

	todoRouter.HandleFunc("/all", isAuthorized(handler.ReadHandler)).Methods("GET")
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