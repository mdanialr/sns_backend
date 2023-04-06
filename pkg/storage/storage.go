package storage

import "io"

// IStorage all implementation of Storage pkg should use this interface as
// their signature/guideline.
type IStorage interface {
	// Save do save given reader to given string as the file path and name
	// including the file extension. This Save will be responsible for closing
	// the reader, so you may safely call this in separate goroutine.
	Save(io.ReadCloser, string)
	// Remove do remove given file path.
	Remove(string)
}
