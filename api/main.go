package main

import (
	api "api/api"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Origin, X-Requested-With, Credentials, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/getusers", corsMiddleware(api.GetUsers))
	mux.HandleFunc("/api/getuser", corsMiddleware(api.GetUser))
	mux.HandleFunc("/api/getposts", corsMiddleware(api.GetPosts))
	mux.HandleFunc("/api/gettopics", corsMiddleware(api.GetTopics))
	mux.HandleFunc("/api/getLikes", corsMiddleware(api.GetLikes))
	mux.HandleFunc("/api/getUnLikes", corsMiddleware(api.GetUnLikes))

	mux.HandleFunc("/api/login", corsMiddleware(handleLoginUser))
	mux.HandleFunc("/api/register", corsMiddleware(handleRegisterUser))
	mux.HandleFunc("/api/passwordForgotten", corsMiddleware(handlePasswordForgotten))
	mux.HandleFunc("/api/accountRecovery", corsMiddleware(handleAccountRecovery))
	mux.HandleFunc("/api/resetPassword", corsMiddleware(handleResetPassword))

	mux.HandleFunc("/api/deleteUser", corsMiddleware(api.DeleteUserHandler))
	mux.HandleFunc("/api/modifyUser", corsMiddleware(api.ModifyUserHandler))
	mux.HandleFunc("/api/uploadImg", corsMiddleware(api.UploadImg))

	mux.HandleFunc("/api/post", corsMiddleware(api.CreatePost))
	mux.HandleFunc("/api/update_post", corsMiddleware(api.UpdatePost))
	mux.HandleFunc("/api/delete_post", corsMiddleware(api.DeletePost))

	mux.HandleFunc("/api/like", corsMiddleware(api.LikePost))
	mux.HandleFunc("/api/unlike", corsMiddleware(api.UnlikePost))
	mux.HandleFunc("/api/dislike", corsMiddleware(api.DislikePost))
	mux.HandleFunc("/api/undislike", corsMiddleware(api.UndislikePost))

	mux.HandleFunc("/api/GameInfo", corsMiddleware(api.GameInfo))
	mux.HandleFunc("/api/ThemeSelection", corsMiddleware(api.ThemeSelection))
	mux.HandleFunc("/api/getCategories", corsMiddleware(api.Categories))

	mux.HandleFunc("/api/topic", corsMiddleware(api.CreateTopic))
	mux.HandleFunc("/api/delete_topic", corsMiddleware(api.DeleteTopic))
	mux.HandleFunc("/api/update_topic", corsMiddleware(api.UpdateTopic))
	mux.HandleFunc("/api/ProfilePageHandler", corsMiddleware(api.ProfilePageHandler))
	mux.HandleFunc("/api/UpdateBiographyHandler", corsMiddleware(api.UpdateBiographyHandler))
	mux.HandleFunc("/api/Rumors", corsMiddleware(api.HandleRumors))

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles("./pages/"+templatePath, "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleLoginUser(w http.ResponseWriter, r *http.Request) {
	errorMessage := api.LoginUser(w, r)

	if errorMessage != "" {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := struct {
			Error string `json:"error"`
		}{errorMessage}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

func handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	errorMessage := api.RegisterUserHandler(w, r)
	if errorMessage != "" {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := struct {
			Error string `json:"error"`
		}{errorMessage}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

func handlePasswordForgotten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	errorMessage := api.PasswordForgottenHandler(w, r)

	if errorMessage != "" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := struct {
			Error string `json:"error"`
		}{errorMessage}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

func handleAccountRecovery(w http.ResponseWriter, r *http.Request) {
	errorMessage := api.AccountRecoveryHandler(w, r)

	if errorMessage != "" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := struct {
			Error string `json:"error"`
		}{errorMessage}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

func handleResetPassword(w http.ResponseWriter, r *http.Request) {
	errorMessage := api.ResetPasswordHandler(w, r)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Error string `json:"error"`
	}{errorMessage}
	json.NewEncoder(w).Encode(jsonResponse)
	return
}
