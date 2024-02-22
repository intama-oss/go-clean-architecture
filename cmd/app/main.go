package main

import "go-clean-architecture/internal/infrastructure"

// @title			Article API Documentation
// @version		1.0
// @description	This is a sample server for Go Clean Architecture.
// @host			localhost:3000
// @BasePath		/api
// @schemes		http https
func main() {
	infrastructure.Run()
}
