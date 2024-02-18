package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB
var tm *templateManager

type Post struct {
	Text     string
	ImageURL string
}

func main() {
	// Initialize template manager
	tm = newTemplateManager()

	// Open database connection
	db = connectDB()
	defer db.Close()

	// Handle routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate_random_post", generateRandomPostHandler)

	// Serve static files
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type templateManager struct {
	templates *template.Template
}

func newTemplateManager() *templateManager {
	return &templateManager{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func (tm *templateManager) executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := tm.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:seh1iWk2MvRySPWhUHp01m1N@tcp(tai.liara.cloud:30983)/trusting_merkle")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Query random post from random_words table
	var postText, postImageBase64 string
	err := db.QueryRow("SELECT paragraph, screenshot FROM random_words, screenshots ORDER BY RAND() LIMIT 1").Scan(&postText, &postImageBase64)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert Base64 image to PNG
	imageURL, err := base64ToPNG(postImageBase64)
	if err != nil {
		log.Println("Error converting Base64 to PNG:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render post page
	tm.executeTemplate(w, "index.html", struct {
		Text     string
		ImageURL string
	}{
		Text:     postText,
		ImageURL: "/uploads/" + imageURL,
	})
}

func generateRandomPostHandler(w http.ResponseWriter, r *http.Request) {
	// Not used in this version
}

func base64ToPNG(base64String string) (string, error) {
	// Decode Base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}

	// Generate unique filename
	filename := uuid.New().String() + ".png"

	// Save PNG file to uploads directory
	err = os.WriteFile(filepath.Join("uploads", filename), imageData, 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}