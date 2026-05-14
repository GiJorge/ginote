package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaytaylor/html2text"
	"github.com/russross/blackfriday/v2"
	_ "modernc.org/sqlite"
)

// --- DATA STRUCTURES ---

type Note struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	AmharicContent string    `json:"amharic_content"`
	Category       string    `json:"category"` // JSON key must be lowercase
	CreatedAt      time.Time `json:"created_at"`
}
var db *sql.DB

// --- DATABASE ---

func initDB() {
	var err error

// Check if we are in Docker or Local
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./data/notes.db" // Local dev
    }

	db, err = sql.Open("sqlite", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		amharic_content TEXT,
		category TEXT,
		created_at DATETIME
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

// --- TRANSLATION LOGIC (NO EXTERNAL LIB NEEDED) ---

func quickTranslate(text string) (string, error) {
	// Google's free web API endpoint
	baseUrl := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=am&dt=t&q="
	resp, err := http.Get(baseUrl + url.QueryEscape(text))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// The response format is a deeply nested JSON array: [[[ "Translated", "Original", ... ]]]
	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil || len(data) == 0 {
		return "", fmt.Errorf("failed to parse translation")
	}

	// Extracting the first translated segment
	inner := data[0].([]interface{})
	translatedText := ""
	for _, segment := range inner {
		s := segment.([]interface{})
		translatedText += fmt.Sprintf("%v", s[0])
	}

	return translatedText, nil
}

// --- HANDLERS ---

func saveNote(c *gin.Context) {
	var input Note
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.ID > 0 {
		// UPDATED: Added amharic_content to the update query
		query := `UPDATE notes SET title=?, content=?, category=?, amharic_content=? WHERE id=?`
		_, err := db.Exec(query, input.Title, input.Content, input.Category, input.AmharicContent, input.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Update failed"})
			return
		}
	} else {
		// New notes usually don't have translation yet, but we include it for consistency
		query := `INSERT INTO notes (title, content, category, amharic_content, created_at) VALUES (?, ?, ?, ?, ?)`
		res, err := db.Exec(query, input.Title, input.Content, input.Category, input.AmharicContent, time.Now())
		if err != nil {
			c.JSON(500, gin.H{"error": "Insert failed"})
			return
		}
		id, _ := res.LastInsertId()
		input.ID = int(id)
	}
	c.JSON(200, input)
}


func translateNote(c *gin.Context) {
	id := c.Param("id")
	var markdown string

	err := db.QueryRow("SELECT content FROM notes WHERE id = ?", id).Scan(&markdown)
	if err != nil {
		c.JSON(404, gin.H{"error": "Note not found"})
		return
	}

	// 1. Strip Markdown to get clean text
	unsafeHTML := blackfriday.Run([]byte(markdown))
	plainText, _ := html2text.FromString(string(unsafeHTML), html2text.Options{})

	// 2. Translate
	amharic, err := quickTranslate(plainText)
	if err != nil {
		c.JSON(500, gin.H{"error": "Translation failed"})
		return
	}

	// 3. Save
	db.Exec("UPDATE notes SET amharic_content = ? WHERE id = ?", amharic, id)

	c.JSON(200, gin.H{"amharic_content": amharic})
}

func getNotes(c *gin.Context) {
	rows, _ := db.Query("SELECT id, title, content, amharic_content, category, created_at FROM notes ORDER BY created_at DESC")
	defer rows.Close()

	var notes []Note = []Note{}
	for rows.Next() {
		var n Note
		rows.Scan(&n.ID, &n.Title, &n.Content, &n.AmharicContent, &n.Category, &n.CreatedAt)
		notes = append(notes, n)
	}
	c.JSON(200, notes)
}

func deleteNote(c *gin.Context) {
	db.Exec("DELETE FROM notes WHERE id = ?", c.Param("id"))
	c.JSON(200, gin.H{"status": "deleted"})
}



// --- MAIN ---

func main() {
	initDB()
	r := gin.Default()

    // 1. Setup CORS (important if you still use different ports in dev)
    r.Use(corsMiddleware())

	
	// Simple CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/notes", saveNote)
		api.GET("/notes", getNotes)
		api.DELETE("/notes/:id", deleteNote)
		api.POST("/notes/:id/translate", translateNote)
	}


	// 3. Serve Frontend Static Files
    // This tells Go: "If a request starts with /assets, look in ./dist/assets"
    r.Static("/assets", "./dist/assets")
    
    // 4. Handle Vue Router (The "Catch-All")
    // If the user hits any other route, send them index.html
    r.NoRoute(func(c *gin.Context) {
        c.File("./dist/index.html")
    })

	fmt.Println("Server running on :8389")
	r.Run(":8389")
}


func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}