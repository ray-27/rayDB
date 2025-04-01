package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// Function to start the web server using Gin
func startServer() {
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the CLI tool web server!",
		})
	})

	// Another example endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PONG",
		})
	})

	// Start the server on port 8080
	go func() {
		fmt.Println("Web server started on http://localhost:8080")
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	}()

	// Gracefully shut down the server on interrupt signals
	gracefulShutdown()
}

// Function to handle graceful shutdown of the server
func gracefulShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	os.Exit(0)
}

// Entry point of the CLI tool
func main() {
	// Define CLI commands and flags
	serverCommand := flag.String("server", "", "Start or stop the server (start/stop)")
	flag.Parse()

	// If the --server flag is passed with "start", start the web server
	if *serverCommand == "start" {
		startServer()
		fmt.Println("Press Ctrl+C to stop the server.")
		select {} // Block the main thread to keep the server running
	} else if *serverCommand == "stop" {
		fmt.Println("Stopping the server...")
		os.Exit(0)
	} else {
		fmt.Println("Usage: --server <start|stop>")
		os.Exit(1)
	}
}
