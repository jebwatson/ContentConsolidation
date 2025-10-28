package main

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// cc add "site" "description"
// cc list
// cc update "site" "description" "id"
// cc delete "id"
func main() {
	// Check if enough arguments were provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: cc <command> [arguments]\nCommands:\n  add <site> <description>\n  list\n  update <site> <description> <id>\n  delete <id>")
	}

	// Figure out what was asked for
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 4 {
			log.Fatal("Usage: cc add <site> <description>")
		}
		CreateContentEntry(os.Args[2], os.Args[3])
	case "list":
		ReadContentEntry()
	case "update":
		if len(os.Args) < 5 {
			log.Fatal("Usage: cc update <site> <description> <id>")
		}
		id, err := bson.ObjectIDFromHex(os.Args[4])
		if err != nil {
			log.Fatal("Could not convert id argument from string to int")
		}

		website := Website{ID: id, Site: os.Args[2], Description: os.Args[3]}
		UpdateContentEntry(website)
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Usage: cc delete <id>")
		}
		id, err := bson.ObjectIDFromHex(os.Args[2])
		if err != nil {
			log.Fatal("Could not convert id argument from string to int")
		}

		DeleteContentEntry(id)
	default:
		log.Fatal("Unknown command: ", command)
	}
}
