package api

import (
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type Post struct {
	Id int
	Title string
	Content string
	Post_date int
	UserId int
	Name string
	Birthday int
	Avatar string
	Num_likes int
}

type UserLiked struct {
	PostID string
	LikedUserID string
}

type PostLikes struct {
	LikeDate int
	UserId int
	UserName string
	Birthday string
	Avatar string
}


func Mount(db *sql.DB) *chi.Mux {
    r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Helloooooo"))
	})
	r.Get("/posts", getPosts(db))
	r.Get("/posts/{id}", getPost(db))
	r.Get("/posts/user={user}", getPostUser(db))
	r.Get("/posts/{id}-{user}", getPostIDUser(db))
	r.Get("/posts/{id}/likes", getLikes(db))
	r.Get("/users/{id}", getUser(db))

    return r
}

func getPosts(db *sql.DB) http.HandlerFunc{
	query := `
			SELECT posts.id, posts.title, posts.content,
			   posts.post_date, users.id, users.name,
			   users.birthday, users.avatar, COUNT(likes.post_id)
			FROM posts 
			INNER JOIN users
				ON posts.author_id = users.id
			INNER JOIN likes
				ON posts.id = likes.post_id
			GROUP BY posts.id, users.id
			`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
    defer rows.Close()

    var posts []Post

    for rows.Next() {
        var p Post
        rows.Scan(&p.Id, &p.Title, &p.Content, &p.Post_date, &p.UserId, &p.Name, &p.Birthday, &p.Avatar, &p.Num_likes)
        posts = append(posts, p)
    }

	postsBytes, _ := json.MarshalIndent(posts, "", "\t")
    
    defer rows.Close()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(postsBytes)
	}
}

func getPost(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		query := fmt.Sprintf(`
		SELECT posts.id, posts.title, posts.content,
		   posts.post_date, users.id, users.name,
		   users.birthday, users.avatar, COUNT(likes.post_id)
		FROM posts 
		INNER JOIN users
			ON posts.author_id = users.id
		INNER JOIN likes
			ON posts.id = likes.post_id
		GROUP BY posts.id, users.id
		HAVING posts.id = ` + idParam)

		row := db.QueryRow(query)
    	var post Post
    	row.Scan(&post.Id, &post.Title, &post.Content, &post.Post_date, &post.UserId, &post.Name, &post.Birthday, &post.Avatar, &post.Num_likes)
		postBytes, _ := json.MarshalIndent(post, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(postBytes)
	}
}

func getPostUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userParam := chi.URLParam(r, "user")
		query := fmt.Sprintf(`
		SELECT posts.id, likes.user_id
		FROM posts 
		LEFT JOIN likes
			ON (posts.id = likes.post_id) AND (likes.user_id = ` + userParam + `)`)

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err)
		}

		var userLikes []UserLiked

		for rows.Next() {
			var likes UserLiked
			rows.Scan(&likes.PostID, &likes.LikedUserID)
			userLikes = append(userLikes, likes)
		}

		userLikesBytes, _ := json.MarshalIndent(userLikes, "", "\t")

		defer rows.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Write(userLikesBytes)
	}
}

func getPostIDUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userParam := chi.URLParam(r, "user")
		idParam := chi.URLParam(r, "id")
		query := fmt.Sprintf(`
		SELECT posts.id, likes.user_id
		FROM posts 
		LEFT JOIN likes
			ON (posts.id = likes.post_id) AND (likes.user_id = ` + userParam + `)
		WHERE posts.id = ` + idParam)

		row := db.QueryRow(query)
		var like UserLiked
		row.Scan(&like.PostID, &like.LikedUserID)
		userLikeBytes, _ := json.MarshalIndent(like, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(userLikeBytes)
	}
}

func getLikes(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		query := fmt.Sprintf(`
		SELECT likes.like_date, likes.user_id, users.name,
			   users.birthday, users.avatar
		FROM likes
		INNER JOIN users
			   ON likes.user_id = users.id
		WHERE likes.post_id = ` + idParam)

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err)
		}

		var likeDataList []PostLikes

		for rows.Next() {
			var likeData PostLikes
			rows.Scan(&likeData.LikeDate, &likeData.UserId, &likeData.UserName, &likeData.Birthday, &likeData.Avatar)
			likeDataList = append(likeDataList, likeData)
		}

		defer rows.Close()

		likeDataByte, _ := json.MarshalIndent(likeDataList, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(likeDataByte)

	}
}

func getUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userParam := chi.URLParam(r, "id")
		query := fmt.Sprintf(`
		SELECT posts.id, posts.title, posts.content,
		   posts.post_date, users.id, users.name,
		   users.birthday, users.avatar, COUNT(likes.post_id)
		FROM posts 
		INNER JOIN users
			ON posts.author_id = users.id
		INNER JOIN likes
			ON posts.id = likes.post_id
		GROUP BY posts.id, users.id
		HAVING posts.author_id = ` + userParam)

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err)
		}

		var latestPosts []Post

		for rows.Next() {
			var p Post
        	rows.Scan(&p.Id, &p.Title, &p.Content, &p.Post_date, &p.UserId, &p.Name, &p.Birthday, &p.Avatar, &p.Num_likes)
        	latestPosts = append(latestPosts, p)
		}

		defer rows.Close()

		latestPostBytes, _ := json.MarshalIndent(latestPosts, "", "\t")

		w.Header().Set("Content_type", "application/json")
		w.Write(latestPostBytes)
	}
}