package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db *mongo.Client
}

const (
	databaseName   = "posts"
	collectionName = "posts"
)

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}
	return data, cur.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), post)
	return err
}

func (s *Store) UpdatePost(post storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.UpdateOne(context.Background(),
		bson.M{"id": post.ID},
		bson.D{
			{"$set", bson.D{
				{"title", post.Title},
				{"content", post.Content},
				{"author_name", post.AuthorName}}},
		})
	return err
}

func (s *Store) DeletePost(post storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteOne(context.Background(),
		bson.M{"id": post.ID})
	return err
}
