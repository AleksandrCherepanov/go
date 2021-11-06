package main

import (
	"log"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-link-api/internal/api"
	"github.com/ozonva/ova-link-api/internal/repo"
	"github.com/rs/zerolog"

	linkAPI "github.com/ozonva/ova-link-api/pkg/ova-link-api"

	"google.golang.org/grpc"
)

func main() {
	const grpcPort = ":82"
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db, err := sqlx.Open("pgx", "postgres://user_links:123456@localhost:5432/user_links?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	linkAPI.RegisterLinkAPIServer(s, api.NewLinkAPI(repo.NewLinkRepo(db), zerolog.New(os.Stdout)))

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return
}
