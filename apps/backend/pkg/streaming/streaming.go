package streaming

import (
	"fmt"
	"io"
	"mime/multipart"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
	ChunkSize   = 32 * 1024         // 32KB
)

func ReadFile(file multipart.File, header *multipart.FileHeader) ([]byte, error) {
	if header.Size > MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d MB)", 
			header.Size, MaxFileSize/(1024*1024))
	}
	
	if header.Size == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	
	// For small files (< 5MB), read directly
	if header.Size <= 5*1024*1024 {
		return io.ReadAll(file)
	}
	
	// For large files, use chunked reading
	data := make([]byte, 0, header.Size)
	chunk := make([]byte, ChunkSize)
	
	for {
		n, err := file.Read(chunk)
		if n > 0 {
			data = append(data, chunk[:n]...)
		}
		
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
	}
	
	return data, nil
}