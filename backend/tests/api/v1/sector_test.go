package v1Test

import (
	v1 "app/internal/api/v1"
	"app/internal/database"
	"context"
	"testing"
)

func TestSector(t *testing.T) {
	// Setup temporary database cache location
	tmpDir, clean := database.TestingTempDir(t, "sectordb_cache_test")
	defer clean()

	v1.NewTestingSector(context.Background(), "log_test.txt", tmpDir, t)

	// TODO: figure out a way to auto-generate test stubs for each endpoint!!
}
