package repository

import "testing"

// TestSetup verifies that the database setup function works correctly
func TestSetup(t *testing.T) {
	repo := getTestRepository()
	err := repo.SetUpDB()
	if err != nil {
		t.Fatalf("SetUpDB failed: %v", err)
	}

	// Test that required tables exist
	reqNames := []string{"projects", "versions", "version_metadata_fields"}
	for _, tableName := range reqNames {
		var name string
		err = repo.DB.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			tableName,
		).Scan(&name)
		if err != nil {
			t.Fatalf("Error checking for table %s: %v", tableName, err)
		}
		if name != tableName {
			t.Errorf("Expected table %s to exist, but it does not", tableName)
		}
	}
}
