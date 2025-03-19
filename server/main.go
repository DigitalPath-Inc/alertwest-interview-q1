package main

import (
	"alertwest-interview-q1/lib"
)

func main() {
	db := lib.NewDB()
	server := NewServer(db)

	// Start components
	db.Run()
	server.Start(":8080")
}
