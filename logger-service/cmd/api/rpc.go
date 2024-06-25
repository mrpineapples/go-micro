package main

import (
	"context"
	"log"
	"time"

	"github.com/mrpineapples/logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(p RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      p.Name,
		Data:      p.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error writing to mongo: ", err)
		return err
	}

	*resp = "Processed payload via RPC: " + p.Name
	return nil
}
