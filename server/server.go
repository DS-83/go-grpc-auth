package server

import (
	"context"
	pb "example-grpc-auth/api"
	"example-grpc-auth/auth/repo/mongodb"
	"example-grpc-auth/auth/usecase"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type key string

// Context config keys
const (
	mongoHost           = key("mongoHost")
	mongoUsername       = key("mongoUsername")
	mongoPwd            = key("mongoPwd")
	mongoDB             = key("mongoDB")
	mongoCredAuthMech   = key("mongoCredAuthMech")
	mongoCredAuthSource = key("mongoCredAuthSource")
	jwtKey              = key("jwtKey")
)

type App struct {
	authServer *usecase.AuthServer
}

func NewApp() *App {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, mongoHost, os.Getenv("MONGO_HOST"))
	ctx = context.WithValue(ctx, mongoUsername, os.Getenv("MONGO_CRED_USER"))
	ctx = context.WithValue(ctx, mongoPwd, os.Getenv("MONGO_CRED__PASSWORD"))
	ctx = context.WithValue(ctx, mongoDB, os.Getenv("MONGO_DATABASE"))
	ctx = context.WithValue(ctx, mongoCredAuthMech, os.Getenv("MONGO_CRED_AUTH_MECH"))
	ctx = context.WithValue(ctx, mongoCredAuthSource, os.Getenv("MONGO_CRED_AUTH_SOURCE"))
	ctx = context.WithValue(ctx, jwtKey, os.Getenv("JWT_SECRET"))

	mongoDB := initMongoDB(ctx)

	userRepo := mongodb.NewUserRepo(mongoDB)
	tokenRepo := mongodb.NewTokenRepo(mongoDB)

	return &App{
		authServer: usecase.NewAuthServer(
			userRepo,
			tokenRepo,
			[]byte(ctx.Value(jwtKey).(string))),
	}
}

func initMongoDB(ctx context.Context) *mongo.Database {
	uri := fmt.Sprintf(
		"mongodb://%s/?maxPoolSize=20&w=majority",
		ctx.Value(mongoHost).(string))

	log.Println(uri)

	clientCred := options.Credential{
		AuthMechanism: ctx.Value(mongoCredAuthMech).(string),
		AuthSource:    ctx.Value(mongoCredAuthSource).(string),
		Username:      ctx.Value(mongoUsername).(string),
		Password:      ctx.Value(mongoPwd).(string),
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(clientCred)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to MongoDB: host:%s db:%s", ctx.Value(mongoHost), ctx.Value(mongoDB))
	return client.Database(ctx.Value(mongoDB).(string))
}

func (a *App) Run(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, a.authServer)

	// Register response service
	reflection.Register(s)

	return s.Serve(lis)
}
