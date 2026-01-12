package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := getPort()
	frontendDir := getFrontendDir()

	fileServer := http.FileServer(http.Dir(frontendDir))
	http.Handle("/", fileServer)

	log.Printf("Starting control panel server on :%s", port)
	log.Printf("Serving frontend from: %s", frontendDir)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "3001"
	}
	return port
}

func getFrontendDir() string {
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir != "" {
		return frontendDir
	}

	executablePath, err := os.Executable()
	if err != nil {
		return "./frontend/dist"
	}

	executableDir := filepath.Dir(executablePath)
	return filepath.Join(executableDir, "frontend", "dist")
}

