/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Implements fs.File interface for BitBucketFS
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/lordofscripts/vfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	ErrSeek error = errors.New("Seek error")
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ vfs.File = (*BitBucketFile)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// BitBucketFile mocks a File returning an error on every operation
// To create a BitBucketFileFS returning a dummy File instead of an error
// you can your own DummyFS:
//
//	type writeDummyFS struct {
//		Filesystem
//	}
//
//	func (fs writeDummyFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
//		return DummyFile(dummyError), nil
//	}

// BitBucketFile represents a dummy File
type BitBucketFile struct {
	name  string // name of the fake file
	size  int64  // fake size of the file
	at    int64  // read pointer of the fake file
	mode  int
	perm  os.FileMode
	err   error // error to return
	mutex *sync.RWMutex
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newBitBucketFile(name string, flag int, perm os.FileMode, err error) BitBucketFile {
	return BitBucketFile{name, 0, 0, flag, perm & os.ModePerm, err, &sync.RWMutex{}}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (f BitBucketFile) String() string {
	return fmt.Sprintf("::%s %d@%d", f.name, f.size, f.at)
}

// Name returns the name of the fake file.
func (f BitBucketFile) Name() string {
	return f.name
}

// Sync returns an error
func (f BitBucketFile) Sync() error {
	return f.err
}

// Truncate reduces the file size to 'size'.
// Errors: none.
func (f BitBucketFile) Truncate(size int64) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if size < f.size {
		f.size = size
	}
	return nil
}

// Close closes the file descriptor.
// Errors: none
func (f BitBucketFile) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.at = 0
	return nil
}

// Write augments the size of the file but does not store any content.
// Errors: none
func (f BitBucketFile) Write(p []byte) (n int, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if p != nil {
		f.size += int64(len(p))
	}
	return len(p), nil
}

// Read pretends it read content.
// Errors: io.EOF
func (f BitBucketFile) Read(p []byte) (n int, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	remains := f.size - f.at
	if int64(len(p)) < remains {
		f.at += int64(len(p))
		p = make([]byte, len(p))
		n = len(p)
		err = nil
	} else {
		f.at = f.size - 1
		n = int(remains)
		p = make([]byte, n)
		err = errors.Join(io.EOF, f.err)
	}
	return n, err
}

// ReadAt pretends to read starting at offset.
// Errors: io.EOF
func (f BitBucketFile) ReadAt(p []byte, offset int64) (n int, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if offset < f.size {
		f.at = offset
		return f.Read(p)
	} else {
		n = 0
		err = errors.Join(io.EOF, f.err)
	}
	return n, f.err
}

// Seek advances the fake file pointer.
// Errors. ErrSeek or the error given to the hybrid BitBucketFS constructor.
func (f BitBucketFile) Seek(offset int64, whence int) (int64, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var newOffset int64
	if whence == 0 {
		newOffset = offset
	} else if whence == 1 {
		newOffset = f.at + offset
	} else if whence == 2 {
		newOffset = f.size - offset
	} else {
		return 0, f.err
	}

	if newOffset < 0 || newOffset >= f.size {
		return 0, errors.Join(f.err, ErrSeek)
	}

	f.at = newOffset
	return f.at, nil
}
