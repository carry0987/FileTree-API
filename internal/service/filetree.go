package service

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/carry0987/FileTree-API/internal/utils"
)

type FileNode struct {
	Name         string      `json:"name"`
	Size         int64       `json:"size,omitempty"`
	FileType     string      `json:"fileType,omitempty"`
	Path         string      `json:"path"`
	CreatedDate  int64       `json:"createdDate,omitempty"`
	LastModified int64       `json:"lastModified,omitempty"`
	IsDir        bool        `json:"isDir"`
	Children     []*FileNode `json:"children,omitempty"`
}

type OrganizedTree struct {
	Dirs  []*FileNode `json:"dirs"`
	Files []*FileNode `json:"files"`
}

type FileTreeResult struct {
	Tree      interface{}
	DirCount  int64
	FileCount int64
}

// GenerateFileTree recursively generates a file tree for the given directory
func GenerateFileTree(root string, organize bool) (*FileTreeResult, error) {
	// Start counting time
	start := time.Now()
	var dirCount, fileCount int64

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
		Name:         filepath.Base(root),
		Path:         root,
		LastModified: info.ModTime().Unix(),
		IsDir:        true,
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Use a buffered channel to control the number of goroutines
	sema := make(chan struct{}, runtime.NumCPU()) // Use the number of CPUs for better concurrency control
	errCh := make(chan error, 1)                  // Error channel
	errWg := sync.WaitGroup{}                     // WaitGroup for error channel

	// Set the root node
	wg.Add(1)
	go walkDir(root, rootNode, &wg, sema, errCh, &dirCount, &fileCount)
	// Wait for all goroutines to finish
	wg.Wait()
	// Close the error channel
	close(errCh)

	// Error handling routine
	errWg.Add(1)
	go func() {
		defer errWg.Done()
		for err := range errCh {
			// Handle errors here, possibly logging them or aggregating into a single error
			if err != nil {
				utils.OutputMessage(nil, utils.LogOutput, 0, "Error: %v", err)
			}
		}
	}()
	errWg.Wait() // Wait for the error handling routine to finish

	var result interface{}
	if organize {
		result = OrganizeFileTree(rootNode)
		utils.OutputMessage(nil, utils.LogOutput, 0, "Organizing file tree for %v", rootNode.Path)
	} else {
		result = rootNode
		utils.OutputMessage(nil, utils.LogOutput, 0, "Get file tree for %v", rootNode.Path)
	}

	// Output the time taken to generate the file tree
	elapsed := time.Since(start)
	utils.OutputMessage(nil, utils.LogOutput, 0, "Total time taken: %v", elapsed)

	fileTreeResult := &FileTreeResult{
		Tree:      result,
		DirCount:  dirCount,
		FileCount: fileCount,
	}

	return fileTreeResult, nil
}

func walkDir(path string, node *FileNode, wg *sync.WaitGroup, sema chan struct{}, errCh chan error, dirCount, fileCount *int64) {
	defer wg.Done()

	// Acquire a semaphore at the start of walkDir to ensure it's released properly
	sema <- struct{}{}
	// Ensure to release semaphore whether the function exits normally or through a return
	defer func() { <-sema }()

	// List entries under the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		errCh <- err // Send the error to the error channel
		return       // Ignore directories that cannot be read
	}

	for _, entry := range entries {
		// Skip hidden files and directories
		if entry.Name()[0] == '.' {
			continue
		}
		// Get the full path of the entry
		fullPath := filepath.Join(path, entry.Name())
		// Node initialization with common properties
		fileInfo, err := entry.Info() // Get file info for common properties
		if err != nil {
			errCh <- err // Send error to the error channel
			continue
		}
		childNode := &FileNode{
			Name:         entry.Name(),
			Path:         fullPath,
			LastModified: fileInfo.ModTime().Unix(), // Set last modified for both files and directories
			IsDir:        entry.IsDir(),
		}

		if !entry.IsDir() {
			atomic.AddInt64(fileCount, 1)
			// Fill additional fields for files
			childNode.Size = fileInfo.Size()
			childNode.FileType = strings.TrimPrefix(filepath.Ext(entry.Name()), ".") // Remove dot from the extension
			childNode.CreatedDate = fileInfo.ModTime().Unix()
		}

		// If it is a directory, recursively traverse the directory
		if entry.IsDir() {
			atomic.AddInt64(dirCount, 1)
			// Use WaitGroup to add a count before recursion
			wg.Add(1)
			go walkDir(fullPath, childNode, wg, sema, errCh, dirCount, fileCount)
		}

		// If it is a file, add it to the children list
		node.Children = append(node.Children, childNode)
	}
}

// Organizes the file tree into a flat list
func OrganizeFileTree(node *FileNode) OrganizedTree {
	organizedTree := OrganizedTree{Dirs: []*FileNode{}, Files: []*FileNode{}}
	collectNodesForOrganizedTree(node, &organizedTree)

	return organizedTree
}

// Helper function to recursively collect nodes into organized dirs and files
func collectNodesForOrganizedTree(node *FileNode, organizedTree *OrganizedTree) {
	if node != nil {
		if node.IsDir {
			// For directories, append to Dirs and skip including children
			newNode := *node       // Create a copy to avoid modifying the original node
			newNode.Children = nil // We don't want children in the dirs
			organizedTree.Dirs = append(organizedTree.Dirs, &newNode)
		} else {
			// For files, simply append to Files
			organizedTree.Files = append(organizedTree.Files, node)
		}

		// Iterate through children if it's a directory
		for _, child := range node.Children {
			collectNodesForOrganizedTree(child, organizedTree)
		}
	}
}
