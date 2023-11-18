package stores

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LanPavletic/go-rest-server/internal/responds"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash []byte `json:"password"`
}

type UserStore struct {
	Collection *mongo.Collection
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserStore(client *mongo.Client) *UserStore {
	return &UserStore{
		Collection: client.Database(DATABASE_NAME).Collection("users"),
	}
}

func (us *UserStore) Login(uc *UserCredentials) (string, error) {
	// get user from store
	var user *User
	filter := bson.D{{Key: "username", Value: uc.Username}}
	err := us.Collection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		// if user not found, return unauthorized
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		panic(err)
	}

	// compare password hashes
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(uc.Password))

	// if passwords don't match, return unauthorized
	if err != nil {
		return "", err
	}

	// passwords match, create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uc.Username,
		"iss": "auth-server",
		"exp": time.Now().Add(time.Hour * 1).Unix(), // token expires in 1 hour
	})

	// sign token with secret
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		// In case token signing fails, return internal server error
		panic(err)
	}
	return tokenString, nil
}

func (us *UserStore) Create(uc *UserCredentials) (primitive.ObjectID, error) {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(uc.Password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err
	}

	user := &User{
		Username:     uc.Username,
		PasswordHash: hash,
	}

	// store user
	result, err := us.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (us *UserStore) Get(id primitive.ObjectID) (*User, error) {
	var user *User

	filter := bson.D{{Key: "_id", Value: id}}
	err := us.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
