package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/anhilmy/tablelink-auth/internal/auth"
	lgrpc "github.com/anhilmy/tablelink-auth/pkg/grpc"
	"github.com/anhilmy/tablelink-auth/pkg/interceptor"
	"github.com/anhilmy/tablelink-auth/repository"
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

	// repo
	repo := repository.NewRepository(db, rdsClient)

	// service
	authServ := auth.NewService(repo, rdsClient)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.IncomingRequest()),
	)

	implServer := NewServer(authServ)
	lgrpc.RegisterAuthServer(s, implServer)
	reflection.Register(s)
	log.Println("gRPC server started on port", PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(-1)
	}

}

type server struct {
	authServ auth.Service
	lgrpc.UnimplementedAuthServer
}

func NewServer(auth auth.Service) *server {
	return &server{
		authServ: auth,
	}
}

func (s server) Login(ctx context.Context, req *lgrpc.LoginRequest) (*lgrpc.LoginResponse, error) {
	return s.authServ.Login(ctx, req)
}

func (s server) GetAllUser(context.Context, *lgrpc.Empty) (*lgrpc.GetAllUserResponse, error) {
	panic("implement me")

}
func (s server) CreateUser(context.Context, *lgrpc.CreateUserRequest) (*lgrpc.Response, error) {
	panic("implement me")

}
func (s server) UpdateUser(context.Context, *lgrpc.UpdateUserRequest) (*lgrpc.Response, error) {
	panic("implement me")

}
func (s server) DeleteUser(context.Context, *lgrpc.DeleteUserRequest) (*lgrpc.Response, error) {
	panic("implement me")

}
