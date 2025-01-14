package api

import (
	"database/sql"
	"time"
)

type LikeDislikeRequest struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	PostID  int `json:"post_id"`
	TopicID int `json:"topic_id"`
}

type User struct {
	Id              int            `json:"id"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	Created_at      string         `json:"created_at"`
	RecoveryCode    *int           `json:"recoveryCode"`
	Profile_picture sql.NullString `json:"profile_picture"`
	Biography       sql.NullString `json:"biography"`
}

type Post struct {
	ID                 int       `json:"id"`
	TopicID            int       `json:"topic_id"`
	Body               string    `json:"body"`
	CreatedBy          int       `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedAtFormatted string
}

type Topic struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Body               string    `json:"body"`
	CreatedBy          int       `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedAtFormatted string
	Category           int `json:"category"`
	CategoryName       string
}

type Game struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"imageLink"`
}

type Categorie struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type APIError struct {
	Error string `json:"error"`
}

type Rumor struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}
