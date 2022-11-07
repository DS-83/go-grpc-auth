package mongodb

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	rTokensT string = "revokedTokens"
)

type TokenRepo struct {
	db *mongo.Database
}
type tokenDB struct {
	ID   string    `bson:"_id"`
	Date time.Time `bson:"revoketion_date"`
}

func NewTokenRepo(db *mongo.Database) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (t TokenRepo) RevokeToken(c context.Context, tokenString string) error {
	cur := t.db.Collection(rTokensT)

	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(tokenString)); err != nil {
		log.Println(err)
		return err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	token := tokenDB{
		ID:   hash,
		Date: time.Now(),
	}

	if _, err := cur.InsertOne(c, token); err != nil {
		return err
	}
	return nil
}

func (t TokenRepo) IsRevoked(c context.Context, token string) (bool, error) {
	cur := t.db.Collection(rTokensT)
	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(token)); err != nil {
		log.Println(err)
		return true, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	res := cur.FindOne(c, bson.M{"_id": hash})
	if res.Err() == mongo.ErrNoDocuments {
		return false, nil
	}

	return true, nil
}
