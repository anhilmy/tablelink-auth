package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	lgrpc "github.com/anhilmy/tablelink-auth/pkg/grpc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = ":2222"
const REDIS_PORT = ":6679"
const REDIS_PASSWORD = ""

const POSTGRES_PORT = 6432
const POSTGRES_USER = "postgres"
const POSTGRES_PASSWORD = "Admin123!!"
const POSTGRES_DB = "postgres"
const POSTGRES_HOST = "localhost"

func main() {

	ctx := context.Background()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go router.Run()

	// postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	// redis
	var rdsClient *redis.Client
	defer func() {
		if rdsClient != nil {
			_ = rdsClient.Close()
		}
	}()
	rdsClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_PORT,
		Password: REDIS_PASSWORD, // no password set
	})
	if err := rdsClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
		os.Exit(-1)
	}

	// grpc
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	implServer := &server{}
	lgrpc.RegisterAuthServer(s, implServer)
	reflection.Register(s)
	log.Println("gRPC server started on port", PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(-1)
	}

}

type server struct {
	lgrpc.UnimplementedAuthServer
}
