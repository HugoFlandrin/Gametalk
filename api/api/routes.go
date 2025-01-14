package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	id := queryValues.Get("id")
	user, err := GetUserData(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get all users
func GetPosts(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetLikes(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllLikes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUnLikes(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllDislikes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllTopics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Register a user
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) string {
	username := r.FormValue("username")
	email := strings.ReplaceAll(r.FormValue("email"), " ", "")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if username == "0" && password == "0" { //admin account
		RegisterUser(username, password, email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return ""
	}

	if !UniqueUsername(username) { //username already used
		return "Ce nom d'utilisateur est déjà utilisé, veuillez en choisir un autre."

	} else if !UniqueEmail(email) { //email already used
		return "Cet email est déjà utilisé, veuillez choisir un autre email."

	} else if !VerifyPassword(password) { //password doesn't follow CNIL recommendations
		return "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux."

	} else if confirmPassword != password { //password and password confirmation don't match
		return "Le mot de passe et la confirmation du mot de passe ne correspondent pas."

	} else { //Account is valid, we can create it
		err := RegisterUser(username, password, email)
		// if err != nil {
		// 	log.Print(err)
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return ""
		// } else {
		sessionID, err := GetUserIdByUsername(username)

		if err != nil {
			http.Error(w, "Erreur de récupération de l'ID de session", http.StatusInternalServerError)
			return ""
		}

		return sessionID
		// }
	}
}

// Log a user
func LoginUser(w http.ResponseWriter, r *http.Request) string {
	username := r.FormValue("username")
	password := r.FormValue("password")
	canConnect, err := AuthenticateUser(username, password)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return ""
	}

	if canConnect {
		sessionID, err := GetUserIdByUsername(username)

		if sessionID == "" {
			sessionID, err = GetUserIdByEmail(username)
		}

		if err != nil {
			http.Error(w, "Erreur de récupération de l'ID de session", http.StatusInternalServerError)
			return ""
		}
		// SetCookie(w, sessionID)
		return sessionID
	} else {
		return "Connexion impossible : nom d'utilisateur ou mot de passe incorrect."
	}
}

// Send an email to the user
func PasswordForgottenHandler(w http.ResponseWriter, r *http.Request) string {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	email := strings.ReplaceAll(r.FormValue("email"), " ", "")

	userExists := false
	query := "SELECT email FROM USERS WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&email)
	if err == nil {
		userExists = true
	} else if err != sql.ErrNoRows {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return ""
	}

	if !userExists {
		return "Cet email n'est associé à aucun compte."
	} else {
		return "Ce service est désactivé."
		// 	fmt.Println("user exists, email valid")
		// 	//create cookie
		// 	sessionID, err := GetUserIdByEmail(email)

		// 	if err != nil {
		// 		http.Error(w, "Erreur de récupération de l'ID de session", http.StatusInternalServerError)
		// 		return ""
		// 	}

		// 	// SetCookie(w, sessionID)

		// 	//send email
		// 	apiKey := "XXX"
		// 	sg := sendgrid.NewSendClient(apiKey)

		// 	code := GenerateCode()
		// 	_, err = db.Exec("UPDATE USERS SET recovery_code = ? WHERE id = ?", code, sessionID)
		// 	if err != nil {
		// 		log.Print(err)
		// 	}

		// 	from := mail.NewEmail("Game Talk", "forumgametalk@gmail.com")
		// 	subject := "Retrouvez votre compte"
		// 	to := mail.NewEmail("Client", email)
		// 	content := mail.NewContent("text/plain", "Bonjour, \n voici le code qui vous permettra de réinitialiser votre mot de passe : "+strconv.Itoa(code)+"\n À bientôt sur Groupie Tracker !")
		// 	message := mail.NewV3MailInit(from, subject, to, content)

		// 	response, err := sg.Send(message)
		// 	if err != nil {
		// 		log.Println("Erreur lors de l'envoi de l'e-mail:", err)
		// 	} else {
		// 		log.Println("Code de statut de l'envoi de l'e-mail:", response.StatusCode)
		// 		log.Println("Réponse de l'API SendGrid:", response.Body)
		// 	}
		// 	return sessionID
	}
}

// Verify the code send
func AccountRecoveryHandler(w http.ResponseWriter, r *http.Request) string {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	recoveryCode := r.FormValue("recoveryCode")
	// userID, err := GetCoockie(w, r, "session_id")
	queryValues := r.URL.Query()
	userID := queryValues.Get("id")
	var userCode int
	_ = db.QueryRow("SELECT recovery_code FROM USERS WHERE id = ?", userID).Scan(&userCode)

	if recoveryCode == strconv.Itoa(userCode) {
		_, _ = db.Exec("UPDATE USERS SET recovery_code = NULL WHERE id = ?", userID)
		return "Succès"
	} else {
		return "Code incorrect."
	}
}

// Reset user password
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) string {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		return "Le mot de passe et la confirmation du mot de passe ne correspondent pas."
	} else if !VerifyPassword(password) {
		return "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux."
	} else {
		queryValues := r.URL.Query()
		userID := queryValues.Get("id")
		db.Exec("UPDATE USERS SET password = ? WHERE id = ?", HashPassword(password), userID)
		return "Succès"
	}
}

// Delete a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	// password := r.FormValue("password")
	queryValues := r.URL.Query()
	userID := queryValues.Get("id")
	// userData, _ := GetUserData(strconv.Itoa(userID))
	db.Exec("DELETE FROM USERS WHERE id = ?", userID)
	// if HashPassword(password) == userData.Password {
	// } else {
	// 	return
	// }
	return
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	db := OpenDatabase("BDD.db")
	defer db.Close()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing form: %v", err)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error getting form file: %v", err)
		return
	}
	defer file.Close()

	dstDir := "static/img/"
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error creating directory: %v", err)
			return
		}
	}

	dstPath := fmt.Sprintf("%s%s", dstDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating file: %v", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error copying file: %v", err)
		return
	}

	queryValues := r.URL.Query()
	userId := queryValues.Get("id")

	// userId, err := GetCoockie(w, r, "session_id")
	// if err != nil {
	// 	http.Error(w, "Session ID cookie not found", http.StatusUnauthorized)
	// 	log.Printf("Session ID cookie not found: %v", err)
	// 	return
	// }

	_, err = db.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", dstPath, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error updating database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"profile_picture": "%s"}`, dstPath)))
	log.Printf("Successfully uploaded image: %s", dstPath)
}

// Modify a user
func ModifyUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request method:", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method"+r.Method, http.StatusMethodNotAllowed)
		return
	}

	queryValues := r.URL.Query()
	sessionID := queryValues.Get("id")
	// id := GetCoockie(w, r, "session_id")
	// sessionID := strconv.Itoa(id)
	userData, err := GetUserData(sessionID)
	if err != nil {
		log.Println("Erreur lors de la récupération des données de l'utilisateur:", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des données de l'utilisateur")
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de l'analyse du formulaire")
		return
	}

	username := ReplaceEmptyString(r.FormValue("usernameInput"), userData.Username)
	email := ReplaceEmptyString(r.FormValue("emailInput"), userData.Email)
	currentPassword := r.FormValue("currentPasswordInput")
	newPassword := r.FormValue("newPasswordInput")
	confirmNewPassword := r.FormValue("confirmNewPasswordInput")

	if currentPassword == "" {
		err := SetUserData(sessionID, username, email, userData.Password)
		if err != nil {
			log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
			SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la mise à jour des données de l'utilisateur")
			return
		}
		SendSuccessResponse(w, "Données de l'utilisateur mises à jour avec succès.", nil)
	} else if HashPassword(currentPassword) != userData.Password {
		SendErrorResponse(w, http.StatusBadRequest, "Le mot de passe actuel saisit est incorrect.")
	} else if newPassword != confirmNewPassword {
		SendErrorResponse(w, http.StatusBadRequest, "Le mot de passe et la confirmation du mot de passe ne correspondent pas.")
	} else if !VerifyPassword(newPassword) {
		SendErrorResponse(w, http.StatusBadRequest, "Votre nouveau mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux.")
	} else {
		err := SetUserData(sessionID, username, email, HashPassword(newPassword))
		if err != nil {
			log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
			SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la mise à jour des données de l'utilisateur")
			return
		}
		SendSuccessResponse(w, "Mot de passe mis à jour avec succès.", nil)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.TopicID <= 0 || post.CreatedBy <= 0 || len(post.Body) == 0 {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	post.CreatedAt = time.Now()

	db := OpenDatabase("BDD.db")
	stmt, err := db.Prepare("INSERT INTO POSTS (topic_id, body, created_by, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	res, err := stmt.Exec(post.TopicID, post.Body, post.CreatedBy, post.CreatedAt)
	if err != nil {
		log.Printf("Error executing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert id: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "post_id": postID})
	if err != nil {
		return
	}
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var topic Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	log.Println(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	topic.CreatedAt = time.Now()

	db := OpenDatabase("BDD.db")
	stmt, err := db.Prepare("INSERT INTO TOPICS (title, body, created_by, created_at, category) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		http.Error(w, "Database error 1", http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Error closing statement: %v\n", err)
		}
	}(stmt)

	res, err := stmt.Exec(topic.Title, topic.Body, topic.CreatedBy, topic.CreatedAt, topic.Category)
	if err != nil {
		log.Printf("Error executing statement: %v\n", err)
		http.Error(w, "Database error 2", http.StatusInternalServerError)
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert id: %v\n", err)
		http.Error(w, "Database error 3", http.StatusInternalServerError)
		return
	}

	category, err := GetCategoryById(topic.Category)
	topic.CreatedAtFormatted = topic.CreatedAt.Format("2006-01-02 15:04:05")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "post_id": postID, "category": category, "time": topic.CreatedAtFormatted})
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}

	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Error unmarshalling body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.ID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("UPDATE POSTS SET body = ? WHERE id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Body, post.ID)
	if err != nil {
		log.Printf("Error executing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}

	var topic Topic
	err = json.Unmarshal(body, &topic)
	if err != nil {
		log.Printf("Error unmarshalling body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if topic.ID <= 0 {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("UPDATE TOPICS SET body = ? WHERE id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(topic.Body, topic.ID)
	if err != nil {
		log.Printf("Error executing statement: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")

	log.Println(postID)
	log.Println(userID)

	if postID == "" || userID == "" {
		http.Error(w, "Missing id or user_id", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	_, err := db.Exec("DELETE FROM POSTS WHERE id = ? AND created_by = ?", postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")

	log.Println(postID)
	log.Println(userID)

	if postID == "" || userID == "" {
		http.Error(w, "Missing id or user_id", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	_, err := db.Exec("DELETE FROM TOPICS WHERE id = ? AND created_by = ?", postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func CleanDatabase() error {
	db, err := sql.Open("sqlite3", "./BDD.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM POSTS WHERE id > 0")
	if err != nil {
		return err
	}

	return nil
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LikeDislikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO LIKES (user_id, post_id, topic_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.UserID, req.PostID, req.TopicID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LikeDislikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM LIKES WHERE user_id = ? AND post_id = ? AND topic_id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.UserID, req.PostID, req.TopicID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LikeDislikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO DISLIKES (user_id, post_id, topic_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.UserID, req.PostID, req.TopicID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func UndislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LikeDislikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := OpenDatabase("BDD.db")
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM DISLIKES WHERE user_id = ? AND post_id = ? AND topic_id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.UserID, req.PostID, req.TopicID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

var db *sql.DB

var initialGames = []Game{
	{
		Title:       "Minecraft",
		Description: "Un jeu qui consiste à placer des blocs et à partir à l'aventure.",
		ImageLink:   "https://image.api.playstation.com/vulcan/img/cfn/11307x4B5WLoVoIUtdewG4uJ_YuDRTwBxQy0qP8ylgazLLc01PBxbsFG1pGOWmqhZsxnNkrU3GXbdXIowBAstzlrhtQ4LCI4.png",
	},
	{
		Title:       "League of Legends",
		Description: "Jeu vidéo d'arène de combat en ligne multijoueurs.",
		ImageLink:   "https://blog.king-jouet.com/wp-content/uploads/2021/08/thumbnail_League-Of-Legends-e1634898736732.jpg",
	},
	{
		Title:       "Valorant",
		Description: "Un jeu de tir à la première personne développé et publié par Riot Games.",
		ImageLink:   "https://e.sport.fr/wp-content/uploads/2020/07/valorant.jpeg",
	},
	{
		Title:       "Warzone",
		Description: "Un jeu vidéo de type battle royale free-to-play.",
		ImageLink:   "https://assets.xboxservices.com/assets/db/88/db8834a9-115d-45e7-a9b5-fa4216b2aac2.jpg",
	},
	{
		Title:       "Overwatch",
		Description: "Un jeu de tir à la première personne multijoueur en équipe développé et publié par Blizzard Entertainment.",
		ImageLink:   "https://www.nintendo.com/eu/media/images/10_share_images/games_15/nintendo_switch_download_software_1/2x1_NSwitchDS_Overwatch2_Season6_image1280w.png",
	},
	{
		Title:       "Fortnite",
		Description: "Un jeu de bataille royale développé et édité par Epic Games.",
		ImageLink:   "https://cdn-0001.qstv.on.epicgames.com/BVJAxewzMErFbvlMLx/image/screen_comp.jpeg",
	},
	{
		Title:       "GTA V",
		Description: "Un jeu d'action-aventure développé par Rockstar North et publié par Rockstar Games.",
		ImageLink:   "https://psblog.fr/wp-content/uploads/2021/05/gta-5.jpg",
	},
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	initDB()
}

func initDB() {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS games (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        image_link TEXT
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	for _, game := range initialGames {
		err := insertGame(game)
		if err != nil {
			log.Printf("Failed to insert game %s: %v", game.Title, err)
		}
	}
}

func insertGame(game Game) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM games WHERE title = ?)", game.Title).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if game exists: %v", err)
	}
	if !exists {
		_, err := db.Exec("INSERT INTO games (title, description, image_link) VALUES (?, ?, ?)",
			game.Title, game.Description, game.ImageLink)
		if err != nil {
			return fmt.Errorf("error inserting game: %v", err)
		}
	}
	return nil
}

func Categories(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ThemeSelection(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, image_link FROM games")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		err := rows.Scan(&game.ID, &game.Title, &game.Description, &game.ImageLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		games = append(games, game)
	}

	response := map[string]interface{}{
		"status": "success",
		"games":  games,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GameInfo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var game Game
	err := db.QueryRow("SELECT id, title, description, image_link FROM games WHERE id = ?", id).Scan(&game.ID, &game.Title, &game.Description, &game.ImageLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "game": game})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var user User
	err := db.QueryRow("SELECT id, username, email, password, profile_picture, biography FROM USERS WHERE id = ?", id).Scan(
		&user.Id, &user.Username, &user.Email, &user.Password, &user.Profile_picture, &user.Biography)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profilePicture := ""
	if user.Profile_picture.Valid {
		profilePicture = user.Profile_picture.String
	}

	biography := ""
	if user.Biography.Valid {
		biography = user.Biography.String
	}

	response := map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"id":              user.Id,
			"username":        user.Username,
			"email":           user.Email,
			"password":        user.Password,
			"profile_picture": profilePicture,
			"biography":       biography,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateBiographyHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Id        string `json:"id"`
		Biography string `json:"biography"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating biography for user ID: %s", requestData.Id)
	log.Printf("New biography: %s", requestData.Biography)

	_, err = db.Exec("UPDATE USERS SET biography = ? WHERE id = ?", requestData.Biography, requestData.Id)
	if err != nil {
		log.Printf("Error updating biography: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRumors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetRumors(w, r)
	case http.MethodPost:
		PostRumor(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetRumors(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT text, image FROM rumors")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rumors []Rumor
	for rows.Next() {
		var rumor Rumor
		if err := rows.Scan(&rumor.Text, &rumor.Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rumors = append(rumors, rumor)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rumors)
}

func PostRumor(w http.ResponseWriter, r *http.Request) {
	var rumor Rumor
	if err := json.NewDecoder(r.Body).Decode(&rumor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO rumors (text, image) VALUES ($1, $2)", rumor.Text, rumor.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
