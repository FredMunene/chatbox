package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/backend/route"
	utils "forum/backend/util"
)

func main() {
	// load environment variables

	if err := utils.LoadEnvVariables(".env"); err != nil {
		fmt.Printf("Error loading .env file :%v", err)
	}

	//  start db

	utils.Init()
	defer utils.Database.Close()

	router := route.InitRoutes()

	server := &http.Server{
		Addr:         ":9000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server running at: http://localhost:%s\n", "9000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
