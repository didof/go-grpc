package main

import (
	"encoding/json"
	"fmt"

	pb "github.com/didof/go-grpc/usermgmt"
	"github.com/go-redis/redis"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserManagementDatabase interface {
	Add(*pb.User) error
	GetAll() (*pb.UsersList, error)
}

const (
	redisAddr = "localhost:6379"
)

func NewUserManagementDatabaseRedis() (*UserManagementDatabaseRedis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("could not establish connection to %s: %v", redisAddr, err)
	}

	return &UserManagementDatabaseRedis{
		client: client,
	}, nil

}

type UserManagementDatabaseRedis struct {
	client *redis.Client
}

func (db *UserManagementDatabaseRedis) Add(user *pb.User) error {
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := db.client.HSet("users", fmt.Sprint(user.Id), json).Err(); err != nil {
		return fmt.Errorf("could not SET the user: %v", err)
	}

	return nil
}

func (db *UserManagementDatabaseRedis) GetAll() (*pb.UsersList, error) {
	users, err := db.client.HGetAll("users").Result()
	if err != nil {
		return nil, fmt.Errorf("could not HGETALL the users: %v", err)
	}

	users_list := new(pb.UsersList)
	for _, jsonUser := range users {
		user := new(pb.User)
		if err := protojson.Unmarshal([]byte(jsonUser), user); err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, user)
	}

	return users_list, nil
}
