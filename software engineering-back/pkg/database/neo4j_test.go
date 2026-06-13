package database

import (
	"os"
	"testing"
)

func TestConnectNeo4jMarksDriverUnavailableWhenConnectivityFails(t *testing.T) {
	originalURI := os.Getenv("NEO4J_URI")
	originalUser := os.Getenv("NEO4J_USER")
	originalPassword := os.Getenv("NEO4J_PASSWORD")
	originalDriver := Neo4jDriver

	t.Cleanup(func() {
		Neo4jDriver = originalDriver
		_ = os.Setenv("NEO4J_URI", originalURI)
		_ = os.Setenv("NEO4J_USER", originalUser)
		_ = os.Setenv("NEO4J_PASSWORD", originalPassword)
	})

	Neo4jDriver = nil
	_ = os.Setenv("NEO4J_URI", "bolt://127.0.0.1:1")
	_ = os.Setenv("NEO4J_USER", "neo4j")
	_ = os.Setenv("NEO4J_PASSWORD", "password")

	ConnectNeo4j()

	if IsNeo4jAvailable() {
		t.Fatal("expected neo4j to be unavailable after connectivity check failure")
	}
}
