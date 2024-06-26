package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mrpineapples/logger/data"
	"github.com/mrpineapples/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, r *logs.LogRequest) (*logs.LogResponse, error) {
	input := r.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

func (app *Config) grpcListen() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("grpc server started on port %s", grpcPort)

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
