package service

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*FileNode `json:"children,omitempty"`
}

// GenerateFileTree recursively generates a file tree for the given directory
func GenerateFileTree(root string) (*FileNode, error) {
	// Make sure the path is normalized
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	// Check if the root exists and is a directory
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err // root does not exist
		}
		return nil, err // other error
	}
	if !info.IsDir() {
		return nil, os.ErrNotExist // root is not a directory
	}

	// Create the root node
	rootNode := &FileNode{
		Name:  filepath.Base(root),
		Path:  root,
		IsDir: true,
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Use a buffered channel to control the number of goroutines
	sema := make(chan struct{}, runtime.NumCPU()) // Use the number of CPUs for better concurrency control

	// A recursive function to fill the file tree nodes
	var walkDir func(string, *FileNode)
	walkDir = func(path string, node *FileNode) {
		defer wg.Done()

		// List entries under the directory
		entries, err := os.ReadDir(path)
		if err != nil {
			return // Ignore directories that cannot be read
		}

		for _, entry := range entries {
			// Skip hidden files and directories
			if entry.Name()[0] == '.' {
				continue
			}
			// Get the full path of the entry
			fullPath := filepath.Join(path, entry.Name())
			// Node initialization
			childNode := &FileNode{
				Name:  entry.Name(),
				Path:  fullPath,
				IsDir: entry.IsDir(),
			}
			// If it is a directory, recursively traverse the directory
			if entry.IsDir() {
				// Use WaitGroup to add a count before recursion
				wg.Add(1)

				// Block until there is space to put a new goroutine
				sema <- struct{}{}

				// Recursively traverse the directory in parallel
				go func() {
					walkDir(fullPath, childNode)
					<-sema
				}()
			}
			// If it is a file, add it to the children list
			node.Children = append(node.Children, childNode)
		}
	}

	// Set the root node
	wg.Add(1)
	walkDir(root, rootNode)

	// Wait for all goroutines to finish
	wg.Wait()

	return rootNode, nil
}
