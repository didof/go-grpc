package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/didof/go-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServe() *UserManagementServer {
	db, err := NewUserManagementDatabaseRedis()
	if err != nil {
		log.Fatal(err)
	}

	return &UserManagementServer{
		db: db,
	}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	db UserManagementDatabase
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)

	log.Printf("Server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   rand.Int31(),
	}

	server.db.Add(created_user)

	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUserParams) (*pb.UsersList, error) {
	return server.db.GetAll()
}

func main() {
	var user_mgmt_server *UserManagementServer = NewUserManagementServe()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
