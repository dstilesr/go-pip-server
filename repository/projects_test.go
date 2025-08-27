package repository

import (
	"context"
	"testing"
)

// TestCreateProject tests the GetOrCreateProject function to ensure it creates
// and retrieves projects correctly
func TestCreateProject(t *testing.T) {
	repo := getTestRepository()
	err := repo.SetUpDB()
	if err != nil {
		t.Fatalf("SetUpDB failed: %v", err)
	}

	ctx := context.Background()
	projectName := "TestProject"

	// First call should create the project
	p1, err := repo.GetOrCreateProject(projectName, ctx)
	if err != nil {
		t.Fatalf("GetOrCreateProject failed: %v", err)
	}
	if p1.Name != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, p1.Name)
	}

	// Second call should retrieve the existing project
	p2, err := repo.GetOrCreateProject(projectName, ctx)
	if err != nil {
		t.Fatalf("GetOrCreateProject failed: %v", err)
	}
	if p2.ID != p1.ID {
		t.Errorf("Expected project ID %d, got %d", p1.ID, p2.ID)
	}
}
