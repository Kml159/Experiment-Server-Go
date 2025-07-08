package parameter

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestSubtractCompleted(t *testing.T) {
    // Create a temporary directory for testing
    tempDir, err := os.MkdirTemp("", "test_received_output")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)

    tests := []struct {
        name           string
        files          []string
        initialMap     map[string]Parameter
        expectedRemain []string
    }{
        {
            name:  "Remove single experiment",
            files: []string{"output_5600100000000x00.txt"},
            initialMap: map[string]Parameter{
                "5600100000000x00": {ID: "5600100000000x00"},
                "5600100000000x01": {ID: "5600100000000x01"},
                "5600100000000x02": {ID: "5600100000000x02"},
            },
            expectedRemain: []string{"5600100000000x01", "5600100000000x02"},
        },
        {
            name:  "Remove multiple experiments",
            files: []string{"output_5600100000000x00.txt", "output_5600100000000x01.csv"},
            initialMap: map[string]Parameter{
                "5600100000000x00": {ID: "5600100000000x00"},
                "5600100000000x01": {ID: "5600100000000x01"},
                "5600100000000x02": {ID: "5600100000000x02"},
                "5600100000000x03": {ID: "5600100000000x03"},
            },
            expectedRemain: []string{"5600100000000x02", "5600100000000x03"},
        },
        {
            name:  "Different file extensions",
            files: []string{"output_5600100000000x10.json", "result_5600100000000x11.xml"},
            initialMap: map[string]Parameter{
                "5600100000000x10": {ID: "5600100000000x10"},
                "5600100000000x11": {ID: "5600100000000x11"},
                "5600100000000x12": {ID: "5600100000000x12"},
            },
            expectedRemain: []string{"5600100000000x12"},
        },
        {
            name:  "Non-matching filenames",
            files: []string{"invalid_file.txt", "no_underscore.csv"},
            initialMap: map[string]Parameter{
                "5600100000000x00": {ID: "5600100000000x00"},
                "5600100000000x01": {ID: "5600100000000x01"},
            },
            expectedRemain: []string{"5600100000000x00", "5600100000000x01"},
        },
        {
            name:       "Empty directory",
            files:      []string{},
            initialMap: map[string]Parameter{
                "5600100000000x00": {ID: "5600100000000x00"},
            },
            expectedRemain: []string{"5600100000000x00"},
        },
        {
            name:  "Files without extension",
            files: []string{"output_5600100000000x20"},
            initialMap: map[string]Parameter{
                "5600100000000x20": {ID: "5600100000000x20"},
                "5600100000000x21": {ID: "5600100000000x21"},
            },
            expectedRemain: []string{"5600100000000x21"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create test files
            for _, fileName := range tt.files {
                filePath := filepath.Join(tempDir, fileName)
                err := os.WriteFile(filePath, []byte("test content"), 0644)
                if err != nil {
                    t.Fatalf("Failed to create test file %s: %v", fileName, err)
                }
            }

            // Copy initial map to avoid modifying the test data
            experimentsMap := make(map[string]Parameter)
            for k, v := range tt.initialMap {
                experimentsMap[k] = v
            }

            // Call the function
            err := SubtractCompleted(&experimentsMap, tempDir, 5600100000000)
            if err != nil {
                t.Errorf("SubtractCompleted() error = %v", err)
                return
            }

            // Check remaining experiments
            if len(experimentsMap) != len(tt.expectedRemain) {
                t.Errorf("Expected %d remaining experiments, got %d", 
                    len(tt.expectedRemain), len(experimentsMap))
            }

            for _, expectedID := range tt.expectedRemain {
                if _, exists := experimentsMap[expectedID]; !exists {
                    t.Errorf("Expected experiment %s to remain, but it was removed", expectedID)
                }
            }

            // Verify removed experiments are actually gone
            for originalID := range tt.initialMap {
                found := false
                for _, remainingID := range tt.expectedRemain {
                    if originalID == remainingID {
                        found = true
                        break
                    }
                }
                if !found {
                    if _, exists := experimentsMap[originalID]; exists {
                        t.Errorf("Expected experiment %s to be removed, but it still exists", originalID)
                    }
                }
            }

            // Clean up test files for next iteration
            for _, fileName := range tt.files {
                os.Remove(filepath.Join(tempDir, fileName))
            }
        })
    }
}

func TestSubtractCompleted_DirectoryNotExists(t *testing.T) {
    experimentsMap := map[string]Parameter{
        "test001": {ID: "test001"},
    }

    err := SubtractCompleted(&experimentsMap, "/nonexistent/directory", 123)
    if err == nil {
        t.Error("Expected error for non-existent directory, got nil")
    }

    // Map should remain unchanged
    if len(experimentsMap) != 1 {
        t.Error("Map should remain unchanged when directory doesn't exist")
    }
}

func TestSubtractCompleted_DirectoryWithSubdirectories(t *testing.T) {
    // Create a temporary directory for testing
    tempDir, err := os.MkdirTemp("", "test_with_subdirs")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)

    // Create subdirectory
    subDir := filepath.Join(tempDir, "subdir")
    err = os.Mkdir(subDir, 0755)
    if err != nil {
        t.Fatalf("Failed to create subdirectory: %v", err)
    }

    // Create files in main directory and subdirectory
    mainFile := filepath.Join(tempDir, "output_5600100000000x30.txt")
    subFile := filepath.Join(subDir, "output_5600100000000x31.txt")

    err = os.WriteFile(mainFile, []byte("test"), 0644)
    if err != nil {
        t.Fatalf("Failed to create main file: %v", err)
    }

    err = os.WriteFile(subFile, []byte("test"), 0644)
    if err != nil {
        t.Fatalf("Failed to create sub file: %v", err)
    }

    experimentsMap := map[string]Parameter{
        "5600100000000x30": {ID: "5600100000000x30"},
        "5600100000000x31": {ID: "5600100000000x31"},
        "5600100000000x32": {ID: "5600100000000x32"},
    }

    err = SubtractCompleted(&experimentsMap, tempDir, 5600100000000)
    if err != nil {
        t.Errorf("SubtractCompleted() error = %v", err)
    }

    // Only the main file should be processed (subdirectories are skipped)
    if len(experimentsMap) != 2 {
        t.Errorf("Expected 2 remaining experiments, got %d", len(experimentsMap))
    }

    if _, exists := experimentsMap["5600100000000x30"]; exists {
        t.Error("Expected experiment from main directory to be removed")
    }

    if _, exists := experimentsMap["5600100000000x31"]; !exists {
        t.Error("Expected experiment from subdirectory to remain (subdirs should be skipped)")
    }
}

func BenchmarkSubtractCompleted(b *testing.B) {
    // Create temp directory
    tempDir, err := os.MkdirTemp("", "bench_test")
    if err != nil {
        b.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)

    // Create 100 test files
    for i := 0; i < 100; i++ {
        fileName := filepath.Join(tempDir, fmt.Sprintf("output_560010000000%02dx%02d.txt", i, i))
        os.WriteFile(fileName, []byte("test"), 0644)
    }

    // Create large experiments map
    experimentsMap := make(map[string]Parameter)
    for i := 0; i < 1000; i++ {
        id := fmt.Sprintf("560010000000%02dx%02d", i, i)
        experimentsMap[id] = Parameter{ID: id}
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Copy map for each iteration
        testMap := make(map[string]Parameter)
        for k, v := range experimentsMap {
            testMap[k] = v
        }

        SubtractCompleted(&testMap, tempDir, 5600100000000)
    }
}