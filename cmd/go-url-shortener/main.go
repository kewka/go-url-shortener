package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kewka/go-url-shortener/internal/handler"
	"github.com/kewka/go-url-shortener/internal/postgres"
	"github.com/kewka/go-url-shortener/internal/service"
)

var (
	port      string
	publicUrl string
)

func init() {
	flag.StringVar(&port, "port", "4000", "server port")
	flag.StringVar(&publicUrl, "public", "http://localhost:4000/", "public url prefix")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbpool, err := postgres.Setup(context.Background())
	if err != nil {
		return err
	}
	defer dbpool.Close()
	handler := handler.New(handler.Config{
		Service:   service.New(dbpool),
		PublicUrl: publicUrl,
	})
	addr := fmt.Sprintf(":%v", port)
	log.Printf("server running on %v", addr)
	return http.ListenAndServe(addr, handler)
}
