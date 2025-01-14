package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", LandingPage)
	mux.HandleFunc("/Login", LoginPage)
	mux.HandleFunc("/Register", RegisterPage)
	mux.HandleFunc("/PasswordForgotten", PasswordForgottenPage)
	mux.HandleFunc("/AccountRecovery", AccountRecoveryPage)
	mux.HandleFunc("/ResetPassword", ResetPasswordPage)
	mux.HandleFunc("/DeleteUser", DeleteUserPage)
	mux.HandleFunc("/ProfileManager", ProfileManagerPage)
	mux.HandleFunc("/Profile", ProfilePage)
	mux.HandleFunc("/PostPage", PostPage)
	mux.HandleFunc("/Topics", TopicsPage)
	mux.HandleFunc("/User", UserPage)
	mux.HandleFunc("/MyPosts", MyPostsPage)
	mux.HandleFunc("/MyLikes", MyLikesPage)
	mux.HandleFunc("/ThemeSelection", ThemeSelection)
	mux.HandleFunc("/GameInfo", GameInfo)
	mux.HandleFunc("/Comments", Comments)
	mux.HandleFunc("/WikiPage", WikiPage)
	mux.HandleFunc("/TournoiPage", TournoiPage)
	mux.HandleFunc("/AvenirPage", AvenirPage)
	mux.HandleFunc("/Rumor", Rumor)
	mux.HandleFunc("/Message", Message)

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":8080", mux)
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

func LandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "LandingPage.html", nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/Login.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/Register.html", nil)
}

func PasswordForgottenPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/PasswordForgotten.html", nil)
}

func AccountRecoveryPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/AccountRecovery.html", nil)
}

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/ResetPassword.html", nil)
}

func DeleteUserPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Authentification/DeleteUser.html", nil)
}

func ProfileManagerPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "ProfileManagerPage.html", nil)
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "ProfilePage.html", nil)
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PostPage.html", nil)
}

func ThemeSelection(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "ThemeSelection.html", nil)
}

func GameInfo(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "GameInfo.html", nil)
}

func Comments(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Comments.html", nil)
}

func TopicsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Topics.html", nil)
}

func UserPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "User.html", nil)
}

func MyPostsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "MyPosts.html", nil)
}

func MyLikesPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "MyLikes.html", nil)
}

func WikiPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "WikiPage.html", nil)
}

func TournoiPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "TournoiPage.html", nil)
}

func AvenirPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "AvenirPage.html", nil)
}

func Rumor(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Rumor.html", nil)
}

func Message(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Message.html", nil)
}
