/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Some examples of using BitBucketFS
 *-----------------------------------------------------------------*/
package bucketfs_test

import (
	"errors"
	"fmt"
	"github.com/lordofscripts/vfs"
	"github.com/lordofscripts/vfs/bucketfs"
	"os"
)

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func ExampleBitBucketFS() {
	const (
		DUMMY_FILE_1             = "/bitbucket/test.pdf"
		DUMMY_FILE_2             = "/bitbucket/test.docx"
		DUMMY_DIR_0              = "/bitbucket"
		DUMMY_DIR_1              = "/bitbucket/Dir1"
		DUMMY_DIR_2              = "/bitbucket/Dir2"
		PERM_DIR     os.FileMode = 0775
		PERM_FILE    os.FileMode = 0666
	)

	// I. Dry Run
	// Create a vfs for dry runs where the operation is only printed/reported
	// and NO errors are produced.
	bfs := bucketfs.Create()

	_, err := bfs.Open(DUMMY_FILE_1)
	_, err = bfs.OpenFile(DUMMY_FILE_1, os.O_RDWR, PERM_FILE)
	err = bfs.Remove(DUMMY_FILE_1)
	err = bfs.Rename(DUMMY_FILE_1, DUMMY_FILE_2)
	err = bfs.Symlink(DUMMY_FILE_1, DUMMY_FILE_2)
	err = bfs.Mkdir(DUMMY_DIR_1, PERM_DIR)
	_, err = bfs.Stat(DUMMY_FILE_1)
	_, err = bfs.Lstat(DUMMY_FILE_2)
	_, err = bfs.ReadDir(DUMMY_DIR_1)
	if err != nil {
		fmt.Println(err) // check on every call. Just an example
	}
	// II. Damp Run
	// Now produce error only if FAKE objects do not exist
	const (
		NO_FILE_1 = "/tmp/dummy1.txt"
		NO_FILE_2 = "/tmp/dummy2.pdf"
		NO_DIR_1  = "/tmp/Dir3"
	)
	// The returned errors will be mostly os.PathError or os.LinkError
	// with ErrBitBucket wrapped into it.
	ErrBitBucket := errors.New("Bitbucket warning!")
	fs := bucketfs.CreateWithError(ErrBitBucket).
		// pretend these directories exist
		WithFakeDirectories([]string{DUMMY_DIR_0}).
		// pretend these files exist
		WithFakeFiles([]string{DUMMY_FILE_1, DUMMY_FILE_2})

	// Some that will NOT produce error
	_, err = bfs.Open(DUMMY_FILE_1)
	_, err = bfs.OpenFile(DUMMY_FILE_2, os.O_RDWR, PERM_FILE)
	err = bfs.Mkdir(NO_DIR_1, PERM_DIR)
	_, err = bfs.ReadDir(DUMMY_DIR_0)

	// These would throw os.PathError with wrapped ErrBitBucket
	err = bfs.Remove(NO_FILE_1)
	err = bfs.Rename(NO_FILE_1, DUMMY_FILE_2)
	err = bfs.Symlink(NO_FILE_2, DUMMY_FILE_2)
	err = bfs.Mkdir(DUMMY_DIR_1, PERM_DIR)
	_, err = bfs.Stat(NO_FILE_1)
	_, err = bfs.Lstat(NO_FILE_2)
	_, err = bfs.ReadDir(DUMMY_DIR_2)

	// These would throw os.LinkError with wrapped ErrBitBucket
	err = bfs.Rename(NO_FILE_1, DUMMY_FILE_2)
	err = bfs.Symlink(NO_FILE_2, DUMMY_FILE_2)

	// III. Dummy Bitbucket files
	var fd vfs.File
	fd, _ = vfs.Create(fs, NO_FILE_1)
	fd.Close()
}
