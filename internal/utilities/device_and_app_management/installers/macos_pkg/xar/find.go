package xar

import (
    "bytes"
    "fmt"
    "io"
    "path/filepath"
    "strings"
)

// FileMatch represents a found file and its metadata
type FileMatch struct {
    ID       uint64
    FullPath string
    Size     int64
    Type     FileType
}

// FindFilesByName searches for all files with the given name in the XAR archive,
// regardless of their location in the directory structure
func FindFilesByName(reader *Reader, filename string) []FileMatch {
    var matches []FileMatch
    
    for id, file := range reader.File {
        // Match against base filename, ignoring path
        if filepath.Base(file.Name) == filename {
            matches = append(matches, FileMatch{
                ID:       id,
                FullPath: file.Name,
                Size:     file.Size,
                Type:     file.Type,
            })
        }
    }
    
    return matches
}

// FindFilesByPattern searches for all files matching a pattern (e.g., "*.txt", "config.*")
func FindFilesByPattern(reader *Reader, pattern string) []FileMatch {
    var matches []FileMatch
    
    for id, file := range reader.File {
        match, err := filepath.Match(pattern, filepath.Base(file.Name))
        if err == nil && match {
            matches = append(matches, FileMatch{
                ID:       id,
                FullPath: file.Name,
                Size:     file.Size,
                Type:     file.Type,
            })
        }
    }
    
    return matches
}

// GetFileContents returns the contents of a file as a byte slice
func GetFileContents(reader *Reader, fileId uint64) ([]byte, error) {
    file, exists := reader.File[fileId]
    if !exists {
        return nil, fmt.Errorf("file with ID %d not found in archive", fileId)
    }
    
    if file.Type != FileTypeFile {
        return nil, fmt.Errorf("ID %d is not a file", fileId)
    }
    
    rc, err := file.Open()
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer rc.Close()
    
    return io.ReadAll(rc)
}

// GetFileContentsByName finds all files with the given name and returns their contents
func GetFileContentsByName(reader *Reader, filename string) map[string][]byte {
    matches := FindFilesByName(reader, filename)
    contents := make(map[string][]byte)
    
    for _, match := range matches {
        if data, err := GetFileContents(reader, match.ID); err == nil {
            contents[match.FullPath] = data
        }
    }
    
    return contents
}

// SearchFileContents searches for files containing the given text
func SearchFileContents(reader *Reader, searchText string) ([]FileMatch, error) {
    var matches []FileMatch
    
    for id, file := range reader.File {
        if file.Type != FileTypeFile {
            continue
        }
        
        data, err := GetFileContents(reader, id)
        if err != nil {
            continue
        }
        
        if bytes.Contains(data, []byte(searchText)) {
            matches = append(matches, FileMatch{
                ID:       id,
                FullPath: file.Name,
                Size:     file.Size,
                Type:     file.Type,
            })
        }
    }
    
    return matches, nil
}

// BuildDirTree builds a tree representation of the XAR archive's directory structure
func BuildDirTree(reader *Reader) *DirNode {
    root := &DirNode{
        Name:     "/",
        Children: make(map[string]*DirNode),
        Files:    make(map[string]*FileMatch),
    }
    
    for id, file := range reader.File {
        parts := strings.Split(strings.Trim(file.Name, "/"), "/")
        current := root
        
        // Navigate/create the directory structure
        for i := 0; i < len(parts)-1; i++ {
            if _, exists := current.Children[parts[i]]; !exists {
                current.Children[parts[i]] = &DirNode{
                    Name:     parts[i],
                    Children: make(map[string]*DirNode),
                    Files:    make(map[string]*FileMatch),
                }
            }
            current = current.Children[parts[i]]
        }
        
        // Add the file to its directory
        fileName := parts[len(parts)-1]
        if file.Type == FileTypeDirectory {
            if _, exists := current.Children[fileName]; !exists {
                current.Children[fileName] = &DirNode{
                    Name:     fileName,
                    Children: make(map[string]*DirNode),
                    Files:    make(map[string]*FileMatch),
                }
            }
        } else {
            current.Files[fileName] = &FileMatch{
                ID:       id,
                FullPath: file.Name,
                Size:     file.Size,
                Type:     file.Type,
            }
        }
    }
    
    return root
}

// DirNode represents a node in the directory tree
type DirNode struct {
    Name     string
    Children map[string]*DirNode
    Files    map[string]*FileMatch
}

// PrintDirTree prints the directory structure
func PrintDirTree(node *DirNode, indent string) {
    fmt.Printf("%s%s/\n", indent, node.Name)
    
    // Print files in current directory
    for fileName, file := range node.Files {
        fmt.Printf("%s  %s (%d bytes)\n", indent, fileName, file.Size)
    }
    
    // Recursively print subdirectories
    for _, child := range node.Children {
        PrintDirTree(child, indent+"  ")
    }
}
