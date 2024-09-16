/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"fmt"
	"os"
	"testing"

	"github.com/lordofscripts/vfs"
	. "github.com/lordofscripts/vfs/test"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FNAME string      = "dummy.txt"
	FSIZE int64       = 1024
	FMODE os.FileMode = 0765
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

// instantiates a BitBucketFileInfo (implements os.FileInfo interface) and
// verify that the methods return the correct information we set.
func Test_FileInfo(t *testing.T) {
	adm := NewUnitTestFramer("BitBucketFS FileInfo interface", t)
	defer adm.TestCaseFrame(t)(t)

	// a regular file
	var fi os.FileInfo
	fi = newBitBucketFileInfo(FNAME, FSIZE, FMODE)
	if err := verifyFileInfo(fi, false, false, FMODE); err != nil {
		t.Error(adm.CryE(nil, err))
	}

	// a directory
	fi = newBitBucketDirInfo(FNAME, FSIZE, FMODE)
	if err := verifyFileInfo(fi, true, false, os.ModeDir|FMODE); err != nil {
		t.Error(adm.CryE(nil, err))
	}

	// a symbolic link
	fi = newBitBucketFileInfo(FNAME, FSIZE, os.ModeSymlink|FMODE)
	if err := verifyFileInfo(fi, false, true, os.ModeSymlink|FMODE); err != nil {
		t.Error(adm.CryE(nil, err))
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

// checks all the public methods of os.FileInfo
func verifyFileInfo(fi os.FileInfo, isDir, isLink bool, mode os.FileMode) error {
	const TEMPLATE = "expected %v got %v"
	if fi.Name() != FNAME {
		return fmt.Errorf("Name: "+TEMPLATE, FNAME, fi.Name())
	}
	if fi.Size() != FSIZE {
		return fmt.Errorf("Size: "+TEMPLATE, FSIZE, fi.Size())
	}
	if fi.Mode() != mode {
		return fmt.Errorf("Mode: "+TEMPLATE, FMODE, fi.Mode())
	}
	if fi.IsDir() != isDir {
		return fmt.Errorf("IsDir: "+TEMPLATE, isDir, fi.IsDir())
	}
	if vfs.HasFileModeFlag(os.ModeSymlink, fi.Mode()) != isLink {
		return fmt.Errorf("Symlink: "+TEMPLATE, isLink, !isLink)
	}
	if fi.Sys() != nil {
		return fmt.Errorf("Sys: "+TEMPLATE, nil, fi.Sys())
	}
	return nil
}
