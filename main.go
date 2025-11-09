package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"petclinic/logger"
)

func main() {
	// Initialize logger with INFO level by default
	logger.SetLevel(os.Getenv("LOG_LEVEL"))
	logger.Info("Starting Pet Clinic application...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file: %v", err)
	}

	// Initialize database
	var err error
	DB, err = InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	// Initialize database logging
	logger.SetDB(DB)

	logger.Info("Database connection established")

	http.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			Register(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			Login(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	// Owners
	http.HandleFunc("/owners", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetOwners(w, r)
		} else if r.Method == http.MethodPost {
			CreateOwner(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	}))

	// Files (local storage)
	http.HandleFunc("/files", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugCtx(r.Context(), "%s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			UploadFile(w, r)
		case http.MethodGet:
			DownloadFile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/owners/id", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetOwnerByID(w, r)
		case http.MethodPut:
			UpdateOwner(w, r)
		case http.MethodDelete:
			DeleteOwner(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	}))

	// Pets
	http.HandleFunc("/pets", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugCtx(r.Context(), "%s %s", r.Method, r.URL.Path)
		if r.Method == http.MethodGet {
			GetPets(w, r)
		} else if r.Method == http.MethodPost {
			CreatePet(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	}))

	http.HandleFunc("/pets/id", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugCtx(r.Context(), "%s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			GetPetByID(w, r)
		case http.MethodPut:
			UpdatePet(w, r)
		case http.MethodDelete:
			DeletePet(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	}))

	// Visits
	http.HandleFunc("/visits", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugCtx(r.Context(), "%s %s", r.Method, r.URL.Path)
		if r.Method == http.MethodGet {
			GetVisits(w, r)
		} else if r.Method == http.MethodPost {
			CreateVisit(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	}))

	http.HandleFunc("/visits/id", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetVisitByID(w, r)
		case http.MethodPut:
			UpdateVisit(w, r)
		case http.MethodDelete:
			DeleteVisit(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	}))

	logger.Info("Server starting on :8080")
	logger.Fatal("Server stopped: %v", http.ListenAndServe(":8080", nil))
}
