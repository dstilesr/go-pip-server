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

// TestCreateProjectVersion tests the basic functionality of CreateProjectVersion
func TestCreateProjectVersion(t *testing.T) {
	repo := getTestRepository()
	err := repo.SetUpDB()
	if err != nil {
		t.Fatalf("SetUpDB failed: %v", err)
	}

	ctx := context.Background()
	pvi := &ProjectVersionInsert{
		ProjectName: "test-project",
		Version:     "1.0.0",
		Digest:      "abc123",
		DigestType:  "sha256",
		FilePath:    "/path/to/file.whl",
		FileType:    "bdist_wheel",
	}

	// Create the first version
	err = repo.CreateProjectVersion(pvi, ctx)
	if err != nil {
		t.Fatalf("CreateProjectVersion failed: %v", err)
	}

	// Verify version was created by getting its ID
	versionId, err := repo.GetLatestProjectVersionId("test-project", ctx, nil)
	if err != nil {
		t.Fatalf("GetLatestProjectVersionId failed: %v", err)
	}
	if versionId <= 0 {
		t.Errorf("Expected positive version ID, got %d", versionId)
	}

	// Verify version data in database
	var digest, digestType, filepath string
	err = repo.DB.QueryRow(
		"SELECT digest, digest_type, filepath FROM versions WHERE id = ?",
		versionId,
	).Scan(&digest, &digestType, &filepath)
	if err != nil {
		t.Fatalf("Error querying version: %v", err)
	}

	if digest != pvi.Digest {
		t.Errorf("Expected digest %s, got %s", pvi.Digest, digest)
	}
	if digestType != pvi.DigestType {
		t.Errorf("Expected digest type %s, got %s", pvi.DigestType, digestType)
	}
	if filepath != pvi.FilePath {
		t.Errorf("Expected filepath %s, got %s", pvi.FilePath, filepath)
	}
}

// TestCreateProjectVersionWithMetadata tests creating a version with metadata
func TestCreateProjectVersionWithMetadata(t *testing.T) {
	repo := getTestRepository()
	err := repo.SetUpDB()
	if err != nil {
		t.Fatalf("SetUpDB failed: %v", err)
	}

	ctx := context.Background()
	metadata := []*KeyVal{
		{Key: "requires-python", Val: ">=3.7"},
		{Key: "author", Val: "Test Author"},
	}

	pvi := &ProjectVersionInsert{
		ProjectName: "test-metadata-project",
		Version:     "2.0.0",
		Digest:      "def456",
		DigestType:  "sha256",
		FilePath:    "/path/to/metadata-file.whl",
		Metadata:    metadata,
		FileType:    "bdist_wheel",
	}

	// Create version with metadata
	err = repo.CreateProjectVersion(pvi, ctx)
	if err != nil {
		t.Fatalf("CreateProjectVersion with metadata failed: %v", err)
	}

	// Verify version was created
	versionId, err := repo.GetLatestProjectVersionId("test-metadata-project", ctx, nil)
	if err != nil {
		t.Fatalf("GetLatestProjectVersionId failed: %v", err)
	}

	// Verify metadata was stored
	rows, err := repo.DB.Query(
		"SELECT key, value FROM version_metadata_fields WHERE version_id = ?",
		versionId,
	)
	if err != nil {
		t.Fatalf("Error querying metadata: %v", err)
	}
	defer rows.Close()

	foundMetadata := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			t.Fatalf("Error scanning metadata row: %v", err)
		}
		foundMetadata[key] = value
	}

	// Check that all metadata was stored correctly
	for _, kv := range metadata {
		value, exists := foundMetadata[kv.Key]
		if !exists {
			t.Errorf("Metadata key %s not found", kv.Key)
			continue
		}
		if value != kv.Val {
			t.Errorf("Expected metadata value %s for key %s, got %s", kv.Val, kv.Key, value)
		}
	}
}
