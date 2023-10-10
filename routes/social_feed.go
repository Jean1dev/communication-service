package routes

import (
	"communication-service/application"
	"communication-service/infra/config"
	"encoding/json"
	"net/http"
	"strings"
)

type SocialPost struct {
	AuthorName   string `json:"authorName"`
	AuthorAvatar string `json:"authorAvatar"`
	Message      string `json:"message"`
	Media        string `json:"media"`
}

type LikePayload struct {
	PostId string `json:"postId"`
	Like   bool   `json:"like"`
	Unlike bool   `json:"unlike"`
}

type CommentPayload struct {
	PostId       string `json:"postId"`
	AuthorName   string `json:"authorName"`
	AuthorAvatar string `json:"authorAvatar"`
	Message      string `json:"message"`
}

func SocialFeedHandler(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	method := r.Method
	if method == "POST" {
		if strings.HasPrefix(r.URL.Path, "/social-feed/like") {
			handleLike(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/social-feed/comment") {
			handleComment(w, r)
			return
		}
		handlePost(w, r)
		return
	}

	if method == "GET" {
		id := strings.TrimPrefix(r.URL.Path, "/social-feed/")
		if id == "" {
			handleGet(w, r)
			return
		}

		handleGetOne(w, r, id)

	}
}

func handleGetOne(w http.ResponseWriter, r *http.Request, id string) {
	post, err := application.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResult, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	var payload CommentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := application.AddComment(payload.Message, payload.PostId, payload.AuthorName, payload.AuthorAvatar); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleLike(w http.ResponseWriter, r *http.Request) {
	var payload LikePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Like {
		err := application.AddLike(payload.PostId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	err, result := application.MyFeed(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var payload SocialPost
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := application.InsertNewPost(payload.AuthorName, payload.AuthorAvatar, payload.Message, payload.Media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
