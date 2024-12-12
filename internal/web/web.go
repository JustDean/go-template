package web

import (
	"context"
	"encoding/json"
	"gon/internal/db"
	"gon/internal/s3"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func postgresHandler(w http.ResponseWriter, r *http.Request) {
	dbpool, err := db.GetPool()
	if err != nil {
		log.Fatalf("Error occured handling db request. %s", err)
	}
	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Postgres seem to work fine!'").Scan(&greeting)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": greeting})
}

func minioHandler(w http.ResponseWriter, r *http.Request) {
	minioClient, err := s3.GetClient()
	defer s3.ReleseClient(minioClient)
	if err != nil {
		log.Fatalf("Error occured trying get minio client. %s\n", err)
	}
	bucketsInfo, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Error occured trying list buckets. %s\n", err)
	}
	for i, buc := range bucketsInfo {
		log.Printf("Bucket %d info is: %v", i, buc)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Minio seems to work fine"})
}

func SetRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/postgres", postgresHandler)
	r.HandleFunc("/minio", minioHandler)
	r.Use(loggingMiddleware)
	return r
}
