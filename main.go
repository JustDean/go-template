package main

import (
	"context"
	"gon/internal/db"
	"gon/internal/s3"
	"gon/internal/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting server")
	ctx := context.Background()

	// set postgres
	func() {
		dbConErr := db.Connect(ctx, db.ConnectionConfig{
			Host:     "localhost",
			Port:     "5432",
			DbName:   "gon",
			Username: "gon",
			Password: "gon",
		})
		if dbConErr != nil {
			log.Fatalf("Failed to connect to the database. %s", dbConErr)
			return
		}
		log.Println("Connected to the database")
	}()

	// set minio
	func() {
		s3.Connect(ctx, s3.MinioConfig{
			PoolSize:        3,
			Endpoint:        "localhost:9000",
			AccessKeyID:     "YiauDxeh5IJ2rlIPvdQr",
			SecretAccessKey: "HulYWV1KpW9VysnDoQ3q1z3jjb6c7RmivOAVXhqj",
			UseSSL:          false,
		})
	}()
	r := web.SetRoutes()
	http.Handle("/", r)
	
	// run in a seperate thread
	go func() {
		ADDR := ":8080"
		log.Printf("Server is starting on %s", ADDR)
		err := http.ListenAndServe(ADDR, r)
		if err != nil {
			log.Fatalf("Failed to start the server %s\n", err)
		}
	}()
	exitSig := make(chan os.Signal, 1)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM)
	<-exitSig
	func() {
		log.Println("Exiting")
		err := db.Disconnect()
		if err != nil {
			log.Fatalf("Error disconnectiong from the databaase. %s", err)
		}
		log.Println("Db connections is closed")
		log.Println("All is closed. Have a nice day!")
	}()
}
