package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Open database
func OpenDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Print(err)
	}

	return db
}

// Get user id with username
func GetUserIdByUsername(username string) (string, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT id FROM USERS WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
		return "", err
	}

	return userID, nil
}

// Get user id with username
func GetUserIdByEmail(email string) (string, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT id FROM USERS WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
		return "", err
	}

	return userID, nil
}

func GetCategoryById(id int) (string, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT name FROM CATEGORIES WHERE id = ?", id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
		return "", err
	}

	return userID, nil
}

// Get user informations with id
func GetUserData(ID string) (User, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	var user User
	err := db.QueryRow("SELECT id, username, email, password, profile_picture FROM USERS WHERE id = ?", ID).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Profile_picture)
	if err != nil {
		log.Println("Erreur lors de la récupération des données de l'utilisateur:", err)
		return User{}, err
	}

	return user, nil
}

// Set user informations
func SetUserData(ID, username, email, password string) error {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	_, err := db.Exec("UPDATE USERS SET username = ?, email = ?, password = ? WHERE id = ?", username, email, password, ID)
	if err != nil {
		log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
		return err
	}

	return nil
}

// Hash a password
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Check that the password complies with CNIL restrictions
func VerifyPassword(password string) bool {
	if len(password) < 12 || strings.ToUpper(password) == password {
		return false
	}

	if ok, _ := regexp.MatchString(`[!@#$%^&*()_+{}\[\]:;<>,.?/~\-]`, password); !ok {
		return false
	}

	if ok, _ := regexp.MatchString(`[0-9]`, password); !ok {
		return false
	}

	return true
}

// Check that the email is unique
func UniqueEmail(email string) bool {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	rowsUsers := selectValueFromTable("USERS", "email")

	for rowsUsers.Next() {
		var userEmail string
		err := rowsUsers.Scan(&userEmail)
		if err != nil {
			fmt.Println(err)
		}
		if userEmail == email {
			return false
		}
	}

	return true
}

// Check that the username is unique
func UniqueUsername(username string) bool {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	rowsUsers := selectValueFromTable("USERS", "username")

	for rowsUsers.Next() {
		var userPseudo string
		err := rowsUsers.Scan(&userPseudo)
		if err != nil {
			fmt.Println(err)
		}
		if userPseudo == username {
			return false
		}
	}

	return true
}

// Log a user
func AuthenticateUser(username, password string) (bool, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	hashedPassword := HashPassword(password)

	var storedPassword string
	query := "SELECT password FROM USERS WHERE username = ? OR email = ?"
	err := db.QueryRow(query, username, username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return storedPassword == hashedPassword, nil
}

// Select a value from a table
func selectValueFromTable(table string, value string) *sql.Rows {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT " + value + " FROM " + table
	result, _ := db.Query(query)
	return result
}

// Select all value from a table
func SelectAllFromTable(database string, table string) *sql.Rows {
	db := OpenDatabase(database)
	defer db.Close()

	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

// Formatting strings
func ReplaceEmptyString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Generate recovery code
func GenerateCode() int {
	return rand.Intn(900000) + 100000
}

func GetAllCategories() ([]Categorie, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(name, ''), COALESCE(description, '') FROM CATEGORIES"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categorie
	for rows.Next() {
		var categorie Categorie
		err := rows.Scan(&categorie.ID, &categorie.Name, &categorie.Description)
		if err != nil {
			fmt.Println(err)
			continue
		}
		categories = append(categories, categorie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetAllTopics() ([]Topic, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(title, ''), COALESCE(body, ''), COALESCE(created_by, ''), created_at, COALESCE(category, 0) FROM TOPICS"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		var createdAt string
		err := rows.Scan(&topic.ID, &topic.Title, &topic.Body, &topic.CreatedBy, &createdAt, &topic.Category)
		if err != nil {
			fmt.Println(err)
			continue
		}

		topic.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", createdAt)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			topic.CreatedAt = time.Time{}
		}

		topic.CreatedAtFormatted = topic.CreatedAt.Format("2006-01-02 15:04:05")

		topic.CategoryName, err = GetCategoryById(topic.Category)
		if err != nil {
			fmt.Println(err)
			continue
		}

		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}

// Get all registered users
func GetAllPosts() ([]Post, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(topic_id, 0), COALESCE(body, ''), COALESCE(created_by, ''), created_at FROM POSTS"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt string
		err := rows.Scan(&post.ID, &post.TopicID, &post.Body, &post.CreatedBy, &createdAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		post.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", createdAt)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			post.CreatedAt = time.Time{}
		}
		post.CreatedAtFormatted = post.CreatedAt.Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Get all registered users
func GetAllUsers() ([]User, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(username, ''), COALESCE(email, ''), COALESCE(password, ''), COALESCE(recovery_code, 0), COALESCE(profile_picture, ''), COALESCE(biography, '') FROM USERS"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RecoveryCode, &user.Profile_picture, &user.Biography)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Get all registered users
func GetAllLikes() ([]LikeDislikeRequest, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(user_id, 0), COALESCE(post_id, 0), COALESCE(topic_id, 0) FROM LIKES"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []LikeDislikeRequest
	for rows.Next() {
		var like LikeDislikeRequest
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.TopicID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		likes = append(likes, like)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func GetAllDislikes() ([]LikeDislikeRequest, error) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	query := "SELECT id, COALESCE(user_id, 0), COALESCE(post_id, 0), COALESCE(topic_id, 0) FROM DISLIKES"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []LikeDislikeRequest
	for rows.Next() {
		var like LikeDislikeRequest
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.TopicID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		likes = append(likes, like)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

// Register a user
func RegisterUser(username, password, email string) error {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	hashedPassword := HashPassword(password)

	_, err := db.Exec("INSERT INTO USERS (username, email, password) VALUES (?, ?, ?)", username, strings.ToLower(email), hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIError{Error: errorMessage})
}

func SendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Message: message, Data: data})
}
