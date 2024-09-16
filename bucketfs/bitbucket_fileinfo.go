/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * BitBucketFS's BitBucketFileInfo FileInfo interface implementation
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"fmt"
	"os"
	"time"

	"github.com/lordofscripts/vfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var (
	_               os.FileInfo       = (*BitBucketFileInfo)(nil)
	defaultFileInfo BitBucketFileInfo = BitBucketFileInfo{"", 0, ALL_PERMS, time.Now(), false, nil}
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// BitBucketFileInfo mocks a os.FileInfo returning default values on every operation
// Struct fields can be set.
// @implements fs.FileInfo interface
type BitBucketFileInfo struct {
	IName    string
	ISize    int64
	IMode    os.FileMode
	IModTime time.Time
	IDir     bool
	ISys     any
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// a FileInfo directory object which implements fs.FileInfo
func newBitBucketFileInfo(name string, size int64, fmode os.FileMode) BitBucketFileInfo {
	if vfs.HasFileModeFlag(os.ModeDir, fmode) {
		fmode = (^os.ModeDir & fmode)
	}
	return BitBucketFileInfo{name, size, fmode, time.Now(), false, nil}
}

// a FileInfo file object which implements fs.FileInfo
func newBitBucketDirInfo(name string, size int64, fmode os.FileMode) BitBucketFileInfo {
	return BitBucketFileInfo{name, size, os.ModeDir | fmode, time.Now(), true, nil}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Stringer renders the FileInfo in a way VERY similar to the ls -l command.
func (fi BitBucketFileInfo) String() string {
	return fmt.Sprintf("%s %7s %7s %8d %s %s", fi.IMode,
		DEF_USER, DEF_GROUP,
		fi.ISize,
		fi.IModTime.Format("01 Jan 2006 15:04"),
		fi.IName)
}

// Name returns the field IName
func (fi BitBucketFileInfo) Name() string {
	return fi.IName
}

// Size returns the field ISize
func (fi BitBucketFileInfo) Size() int64 {
	return fi.ISize
}

// Mode returns the field IMode
func (fi BitBucketFileInfo) Mode() os.FileMode {
	return fi.IMode
}

// ModTime returns the field IModTime
func (fi BitBucketFileInfo) ModTime() time.Time {
	return fi.IModTime
}

// IsDir returns the field IDir
func (fi BitBucketFileInfo) IsDir() bool {
	return fi.IDir
}

// Sys returns the field ISys
func (fi BitBucketFileInfo) Sys() any {
	return fi.ISys
}
