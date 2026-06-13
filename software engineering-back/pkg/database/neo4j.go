package database

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Neo4jDriver neo4j.DriverWithContext

func ConnectNeo4j() {
	uri := os.Getenv("NEO4J_URI")
	user := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	if uri == "" {
		Neo4jDriver = nil
		log.Println("warning: NEO4J_URI not set, skipping Neo4j connection")
		return
	}

	var err error
	Neo4jDriver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		Neo4jDriver = nil
		log.Printf("warning: failed to create neo4j driver: %v", err)
		return
	}

	ctx := context.Background()
	err = Neo4jDriver.VerifyConnectivity(ctx)
	if err != nil {
		_ = Neo4jDriver.Close(ctx)
		Neo4jDriver = nil
		log.Printf("warning: neo4j connectivity check failed: %v", err)
		return
	}
	log.Println("neo4j connected successfully")
}

func CloseNeo4j() {
	if Neo4jDriver != nil {
		Neo4jDriver.Close(context.Background())
	}
}

func IsNeo4jAvailable() bool {
	return Neo4jDriver != nil
}
