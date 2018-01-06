package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ExportResult is the result of an export operation... duh
type ExportResult struct {
	// Path to exported file
	Path string

	// MIME type of the exported file (e.g. application/x-tar)
	MIME string

	// Any error that occurred during `Export()`
	Error *Error
}

// Storer takes an `ExportResult` and move it somewhere! To a cloud storage service, for instance...
type Storer interface {
	Store(result *ExportResult, directory string) *Error
}

//
// To hands off an ExportResult to a `Storer` interface and invokes its Store() method.
// The directory argument is passed along too. If `store` is `nil`, the the method will simply move the export
// result to the specified directory (via the `mv` command)
//
func (x *ExportResult) To(directory string, store Storer) *Error {

	if store == nil {
		out, err := exec.Command("mv", x.Path, directory+x.Filename()).Output()
		return makeErr(err, string(out))
	}

	storeErr := store.Store(x, directory)
	if storeErr != nil {
		return storeErr
	}

	err := os.Remove(x.Path)
	return makeErr(err, "")
}

//
// Filename returns the just filename component of the `Path` attribute
//
func (x ExportResult) Filename() string {
	_, filename := filepath.Split(x.Path)
	return filename
}

/* End File */
