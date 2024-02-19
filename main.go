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
    // Check if the random_posts table exists
    rows, err := db.Query("SHOW TABLES LIKE 'random_posts'")
    if err != nil {
        log.Println("Error checking for table existence:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // If the table doesn't exist or is empty, display a message
    if !rows.Next() {
        fmt.Fprint(w, "No posts available.")
        return
    }

    // Query posts from random_posts table in descending order
    rows, err = db.Query("SELECT text, image_url FROM random_posts ORDER BY id DESC")
    if err != nil {
        log.Println("Error querying database:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Iterate over the rows and collect posts
    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.Text, &post.ImageURL); err != nil {
            log.Println("Error scanning row:", err)
            continue
        }
        posts = append(posts, post)
    }
    if err := rows.Err(); err != nil {
        log.Println("Error iterating over rows:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // If there are no posts, display a message
    if len(posts) == 0 {
        fmt.Fprint(w, "No posts available.")

		
        return
    }

    // Render posts in reverse order
    for i := len(posts) - 1; i >= 0; i-- {
        tm.executeTemplate(w, "index.html", posts[i])
    }

    // Add a button to generate a random post
    fmt.Fprintf(w, `<form action="/generate_random_post" method="post">
                    <input type="submit" value="Generate Random Post">
                </form>`)
}

func generateRandomPostHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the table exists, if not, create it
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS random_posts (
            id INT AUTO_INCREMENT PRIMARY KEY,
            text TEXT,
            image_url VARCHAR(255)
        ) ENGINE=InnoDB;
    `)
    if err != nil {
        log.Println("Error creating table:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Query random post from random_words table
    var postText, postImageBase64 string
    err = db.QueryRow("SELECT paragraph, screenshot FROM random_words, screenshots ORDER BY RAND() LIMIT 1").Scan(&postText, &postImageBase64)
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

    // Insert the post into the random_posts table
    _, err = db.Exec("INSERT INTO random_posts (text, image_url) VALUES (?, ?)", postText, imageURL)
    if err != nil {
        log.Println("Error inserting post into database:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Redirect to the main page to display the new post
    http.Redirect(w, r, "/", http.StatusSeeOther)
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
