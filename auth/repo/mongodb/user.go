package mongodb

import (
	"context"
	e "example-grpc-auth/err"
	"example-grpc-auth/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	talbleUsers = "users"
)

type UserRepo struct {
	db *mongo.Database
}

type user struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	MysqlID  int                `bson:"mysql_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(c context.Context, u string, p string) error {
	cur := r.db.Collection(talbleUsers)

	user := &user{
		Username: u,
		Password: p,
	}

	if _, err := cur.InsertOne(c, user); err != nil {

		if mongo.IsDuplicateKeyError(err) {
			return e.ErrDupKey
		}
	}

	return nil
}

func (r *UserRepo) GetUser(c context.Context, u string, p string) (*models.User, error) {
	cur := r.db.Collection(talbleUsers)

	user := new(user)

	err := cur.FindOne(c, bson.M{"username": u}).Decode(user)
	if err == mongo.ErrNoDocuments {
		return nil, e.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p))
	if err != nil {
		return nil, e.ErrInvalidCred
	}
	return toModelsUser(user), nil
}

func (r *UserRepo) DeleteUser(c context.Context, u *models.User) error {
	user := toDBUser(u)

	cur := r.db.Collection(talbleUsers)

	res, err := cur.DeleteOne(c, user)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return e.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) UpdateUser(c context.Context, filt *models.User, upd *models.User) (*models.User, error) {
	filtDB := toDBUser(filt)
	updDB := toDBUser(upd)

	update := bson.M{
		"$set": updDB,
	}

	cur := r.db.Collection(talbleUsers)

	res, err := cur.UpdateOne(c, filtDB, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, e.ErrUserNotFound
	}

	user := new(user)
	err = cur.FindOne(c, bson.M{"_id": filtDB.ID}).Decode(user)
	if err != nil {
		return nil, err
	}

	return toModelsUser(user), nil

}

func toDBUser(u *models.User) *user {
	id, _ := primitive.ObjectIDFromHex(u.ID)
	return &user{
		ID:       id,
		MysqlID:  u.MysqlID,
		Username: u.Username,
		Password: u.Password,
	}
}

func toModelsUser(u *user) *models.User {
	return &models.User{
		ID:       u.ID.Hex(),
		MysqlID:  u.MysqlID,
		Username: u.Username,
		Password: u.Password,
	}
}
