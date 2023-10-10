package application

import (
	"communication-service/infra/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	postCollection = "feed-posts"
	dateFormat     = "02/01/2006 15:04"
)

type PostEntity struct {
	Id        string          `json:"id"`
	Author    AuthorEntity    `json:"author"`
	CreatedAt string          `json:"createdAt"`
	IsLiked   bool            `json:"isLiked"`
	Likes     int             `json:"likes"`
	Message   string          `json:"message"`
	Media     string          `json:"media"`
	Comments  []CommentEntity `json:"comments"`
}

type CommentEntity struct {
	Id        string       `json:"id"`
	Author    AuthorEntity `json:"author"`
	CreatedAt string       `json:"createdAt"`
	Message   string       `json:"message"`
}

type AuthorEntity struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

func NewComment(user string, avatar string, comment string) CommentEntity {
	author := &AuthorEntity{
		Avatar: avatar,
		Name:   user,
	}

	commentEntity := &CommentEntity{
		Id:        strconv.Itoa(rand.Int()),
		Author:    *author,
		CreatedAt: time.Now().Format(dateFormat),
		Message:   comment,
	}

	return *commentEntity
}

func NewPostEntityFromInputBody(
	authorName string,
	authorAvatar string,
	message string,
	media string) PostEntity {
	author := &AuthorEntity{
		Avatar: authorAvatar,
		Name:   authorName,
	}

	post := &PostEntity{
		Author:    *author,
		IsLiked:   false,
		Likes:     0,
		Message:   message,
		Media:     media,
		CreatedAt: time.Now().Format(dateFormat),
		Comments:  make([]CommentEntity, 0),
	}
	return *post
}

func (p *PostEntity) Validate() error {
	if p.Author == (AuthorEntity{}) {
		return errors.New("Author is required")
	}

	if p.Message == "" {
		return errors.New("Message is required")
	}

	if p.Author.Name == "" {
		return errors.New("Author name is required")
	}

	return nil
}

func FindById(id string) (*PostEntity, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	db := database.GetDB()
	_, data := db.FindOne(postCollection, filter)

	var result PostEntity
	if err := data.Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Entity not found")
		}

		return nil, err
	}

	return &result, nil
}

func InsertNewPost(authorName string,
	authorAvatar string,
	message string,
	media string) error {
	post := NewPostEntityFromInputBody(
		authorName,
		authorAvatar,
		message,
		media)

	if err := post.Validate(); err != nil {
		return err
	}

	db := database.GetDB()
	err := db.Insert(post, postCollection)

	if err != nil {
		return err
	}

	return nil
}

func AddComment(comment string, postId string, user string, avatar string) error {
	db := database.GetDB()
	id, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return err
	}

	commentEntity := NewComment(user, avatar, comment)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{{Key: "comments", Value: commentEntity}}},
	}

	if err := db.UpdateOne(postCollection, filter, update); err != nil {
		return err
	}

	return nil
}

func AddLike(postId string) error {
	db := database.GetDB()
	id, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "likes", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "isLiked", Value: true}}},
	}

	if err := db.UpdateOne(postCollection, filter, update); err != nil {
		return err
	}

	return nil
}

func MyFeed(username string) (error, []PostEntity) {
	db := database.GetDB()

	log.Printf("buscando posts de %s", username)

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	err, cursor := db.FindAll(postCollection, bson.D{}, opts)
	if err != nil {
		log.Panic(err)
	}

	var results []PostEntity
	for cursor.Next(context.Background()) {
		var doc bson.M
		err := cursor.Decode(&doc)
		if err != nil {
			log.Print(err)
			return err, nil
		}

		jsonData, err := bson.MarshalExtJSON(doc, false, false)
		if err != nil {
			log.Print(err)
			return err, nil
		}

		var post PostEntity
		err = json.Unmarshal(jsonData, &post)
		if err != nil {
			log.Print(err)
			return err, nil
		}

		post.Id = doc["_id"].(primitive.ObjectID).Hex()
		post.CreatedAt = calculateDuration(post.CreatedAt)
		results = append(results, post)
	}

	return nil, results
}

func calculateDuration(dateString string) string {
	t, err := time.Parse(dateFormat, dateString)
	if err != nil {
		log.Panic(err)
	}

	now := time.Now()
	duration := now.Sub(t)

	years := int(duration.Hours() / 24 / 365)
	months := int(duration.Hours()/24/30) % 12
	days := int(duration.Hours()/24) % 30
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if years > 0 {
		return fmt.Sprintf("%d ano(s) atrás", years)
	} else if months > 0 {
		return fmt.Sprintf("%d mês(es) atrás", months)
	} else if days > 0 {
		return fmt.Sprintf("%d dia(s) atrás", days)
	} else if hours > 0 {
		return fmt.Sprintf("%d hora(s) e %d minuto(s) atrás", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minuto(s) atrás", minutes)
	}

	return "Agora"
}
