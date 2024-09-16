/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * BitBucket half-dummy filesystem implementation
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/lordofscripts/vfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	BITBUCKET_FILE_SIZE int64 = 512
	BITBUCKET_DIR_SIZE  int64 = 0

	cDIR                rune = 'Ⅾ'
	cFILE               rune = 'Ⅎ'
	cDASHED_RIGHT_ARROW rune = '⇢'
	cLINK               rune = '⛓'
)

var (
	ErrFileExists error = errors.New("File exists")
	ErrDirExists  error = errors.New("Directory exists")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ vfs.Filesystem = (*BitBucketFS)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// BitBucketFS is fake filesystem which has 2 modes of operation:
//  1. In Silent Mode no errors are produced, all operations are simply printed
//     to the console (stdout).
//  2. In Hybrid Mode a user-defined error is specified and optionally a list of
//     FAKE files and directories.  Then the outcome of any operation depends on
//     whether the requested file/directory is present in the FAKE list as it
//     would in a real file system. Then it would produce a PathError or LinkError
//     under the same conditions a a real underlying OS would do, else no error.
type BitBucketFS struct {
	mutex   *sync.RWMutex
	silent  bool  // simply print operation & param and produces no error
	err     error // if not nil, same error for all entry points
	entries map[string]iBucketNode
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Create is the silent mode constructor. All operations only print a message
// on the console. No error is produced.
// @implements vfs.Filesystem
func Create() *BitBucketFS {
	return &BitBucketFS{&sync.RWMutex{}, true, nil, make(map[string]iBucketNode, 0)}
}

// CreateWithError is the constructor for hybrid mode. Alternatively use the
// silent mode constructor followed by the WithError() fluent API call. In
// this mode you MAY use the fluent API WithFakeDirectories() and
// WithFakeFiles() calls.
// @implements vfs.Filesystem
func CreateWithError(err error) *BitBucketFS {
	if err == nil {
		panic("BitBucketFS ctor. needs err parameter")
	}
	return &BitBucketFS{&sync.RWMutex{}, false, err, make(map[string]iBucketNode, 0)}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// WithFakeDirectories registers fake directories. It only applies only to
// non-silent mode.
// In silent mode it returns without doing anything.
func (bfs *BitBucketFS) WithFakeDirectories(dirs []string) *BitBucketFS {
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	if !bfs.silent {
		for i, name := range dirs {
			dirs[i] = filepath.Clean(name)
		}

		sort.Strings(dirs)
		for _, name := range dirs {
			bfs.entries[name] = createNode(os.ModeDir | ALL_PERMS)
		}
	}

	return bfs
}

// WithFakeFiles registeres fake file objects. It only applies only to
// non-silent mode.
// In silent mode it returns without doing anything.
func (bfs *BitBucketFS) WithFakeFiles(files []string) *BitBucketFS {
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	if !bfs.silent {
		for i, name := range files {
			files[i] = filepath.Clean(name)
		}

		sort.Strings(files)

		for _, name := range files {
			bfs.entries[name] = createNode(ALL_RW_PERMS)
		}
	}

	return bfs
}

// WithError is a fluent API call that can be applied when using the silent
// mode constructor.
func (bfs *BitBucketFS) WithError(customError error) *BitBucketFS {
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	bfs.silent = false
	bfs.err = customError

	return bfs
}

// PathSeparator returns the path separator
func (bfs *BitBucketFS) PathSeparator() uint8 {
	return '/'
}

// Remove deletes the file/directory called name.
// Errors: fs.PathError
func (bfs *BitBucketFS) Remove(name string) error {
	return bfs.executePath("Remove", name)
}

// Rename renames oldPath to newPath
// Errors: os.LinkError
func (bfs *BitBucketFS) Rename(oldPath, newPath string) error {
	const Oper = "Rename"
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	oldPath = filepath.Clean(oldPath)
	newPath = filepath.Clean(newPath)

	if bfs.silent {
		fmt.Printf("\t⚡ Rename %s %c %s\n", oldPath, cDASHED_RIGHT_ARROW, newPath)
	} else {
		node, oldExists := bfs.entries[oldPath]
		if !oldExists {
			return &os.LinkError{Oper, oldPath, newPath, fmt.Errorf("OldPath does not exist. %w", bfs.err)}
		}

		_, newExists := bfs.entries[newPath]
		if newExists {
			return &os.LinkError{Oper, oldPath, newPath, fmt.Errorf("NewPath exists. %w", bfs.err)}
		}

		// on Rename the old name disappears
		delete(bfs.entries, oldPath)
		// and the new name is added
		bfs.entries[newPath] = node
		fmt.Printf("\t⚡ Rename %c %s %c %s\n", xlateIsDir(node.IsDir()), oldPath, cDASHED_RIGHT_ARROW, newPath)
	}

	return nil
}

// Symlink symbolic links/points newName to oldName
// Errors: os.LinkError
func (bfs *BitBucketFS) Symlink(oldName, newName string) error {
	const Oper = "Symlink"
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	oldName = filepath.Clean(oldName)
	newName = filepath.Clean(newName)

	if bfs.silent {
		fmt.Printf("\t⚡ Symlink %s ⇉ %s\n", newName, oldName)
	} else {
		node, oldExists := bfs.entries[oldName]
		if !oldExists { // nothing to point to
			return &os.LinkError{Oper, oldName, newName, fmt.Errorf("OldName (target) does not exist. %w", bfs.err)}
		}

		_, newExists := bfs.entries[newName]
		if newExists { // source can't be created
			return &os.LinkError{Oper, oldName, newName, fmt.Errorf("NewName exists. %w", bfs.err)}
		}

		// For Symlink the old remains and the new is created
		bfs.entries[newName] = node.LinkTo(oldName)
		fmt.Printf("\t⚡ Symlink %c %s ⇉ %s\n", xlateIsDir(node.IsDir()), newName, oldName)
	}

	return nil
}

// Mkdir creates a FAKE directory.
// Errors: fs.PathError
func (bfs *BitBucketFS) Mkdir(name string, perm os.FileMode) error {
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	var errx error = nil
	name = filepath.Clean(name)
	if bfs.silent {
		fmt.Printf("\t⚡ Mkdir %s\n", name)
	} else {
		if node, exists := bfs.entries[name]; exists {
			fmt.Printf("\t⚡ Mkdir %c %s\n", xlateIsDir(node.IsDir()), name)

			if node.IsDir() {
				errx = &os.PathError{"Mkdir", name, errors.Join(bfs.err, ErrDirExists)}
			} else {
				errx = &os.PathError{"Mkdir", name, errors.Join(bfs.err, ErrFileExists)}
			}
		} else {
			fmt.Printf("\t⚡ Mkdir %c %s\n", xlateIsDir(true), name)
			bfs.entries[name] = createNode(os.ModeDir | perm)
		}
	}
	return errx
}

// Open opens the file.
// Errors: fs.PathError
func (bfs *BitBucketFS) Open(name string) (vfs.File, error) {
	return bfs.OpenFile(name, os.O_RDONLY, 0)
}

// OpenFile returns dummy error
// Errors: fs.PathError
func (bfs *BitBucketFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	const Oper = "Open/OpenFile"
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	// Create the File Descriptor when file opened or created
	onFileOpened := func(name string, mode int, perm os.FileMode) vfs.File {
		fmt.Printf("\t⚡ %s %s\n", Oper, name)
		return newBitBucketFile(name, mode, perm, bfs.err)
	}

	name = filepath.Clean(name)
	if bfs.silent {
		return onFileOpened(name, flag, perm), nil
	} else {
		if node, exists := bfs.entries[name]; exists {
			if node.IsDir() {
				return nil, &os.PathError{Oper, name, errors.Join(bfs.err, vfs.ErrIsDirectory)}
			} else {
				return onFileOpened(name, flag, perm), nil
			}
		} else {
			// doesn't exist but O_CREATE given
			if vfs.HasModeFlag(os.O_CREATE, flag) {
				return onFileOpened(name, flag, perm), nil
			}
			return nil, &os.PathError{Oper, name, bfs.err}
		}
	}
}

// Stat returns information about the file. If the file is a Symbolic Link,
// then the link is followed and it instead returns information about the
// file it points to.
// Errors: nil or fs.PathError
func (bfs *BitBucketFS) Stat(name string) (os.FileInfo, error) {
	const Oper string = "Stat"

	if bfs.silent {
		fmt.Printf("\t⚡ %s %s\n", Oper, name)
		return defaultFileInfo, nil
	}

	node, cleaned := bfs.findNode(name)
	if node == nil {
		return nil, &os.PathError{Oper, cleaned, os.ErrNotExist}
	}

	if node.IsLink() {
		var err error
		node, err = bfs.dereferenceSymLink(Oper, node.Target())
		if err != nil {
			return nil, err // broken symlink
		}
		// target node found!
		cleaned = node.ClientData().(string)
		return bfs.internalStat(Oper, cleaned) // Stat target node
	} else {
		return bfs.internalStat(Oper, cleaned) // Stat this node
	}
}

// Same as Stat() except that if the file is a Symbolic Link, it returns
// information about the symlink, not the file it points to.
// Errors: nil or fs.PathError
func (bfs *BitBucketFS) Lstat(name string) (os.FileInfo, error) {
	return bfs.internalStat("Lstat", name)
}

// ReadDir returns PathError in non-silent mode  if path is not in FakeDirectories
func (bfs *BitBucketFS) ReadDir(path string) ([]os.FileInfo, error) {
	bfs.mutex.RLock()
	defer bfs.mutex.RUnlock()

	Empty := make([]os.FileInfo, 0)
	path = filepath.Clean(path)
	if bfs.silent {
		fmt.Printf("\t⚡ ReadDir %s\n", path)
		return Empty, nil
	} else {
		if node, exists := bfs.entries[path]; exists {
			if node.IsDir() {
				return Empty, nil // TODO? is it worth to fake []FileInfo?
			} else {
				return Empty, &os.PathError{"ReadDir", path, errors.Join(bfs.err, vfs.ErrNotDirectory)}
			}
		} else {
			return Empty, &os.PathError{"ReadDir", path, bfs.err}
		}
	}
}

/* ----------------------------------------------------------------
 *				P r i v a t e		M e t h o d s
 *-----------------------------------------------------------------*/

// Stat() a file regardless of whether it is SymLink or not. Does NOT
// follow a link, that is left up to Stat()
// Errors: nil or fs.PathError
func (bfs *BitBucketFS) internalStat(operation, name string) (os.FileInfo, error) {
	bfs.mutex.RLock()
	defer bfs.mutex.RUnlock()

	var finfo BitBucketFileInfo = defaultFileInfo
	var node iBucketNode
	var exists bool = false
	name = filepath.Clean(name)
	if node, exists = bfs.entries[name]; exists {
		if node.IsDir() {
			finfo = newBitBucketDirInfo(name, node.Size(), node.Mode())
		} else {
			finfo = newBitBucketFileInfo(name, node.Size(), node.Mode())
		}
	} else {
		finfo.IName = name
	}

	if bfs.silent {
		fmt.Printf("\t⚡ %s %c %s\n", operation, xlateIsDir(finfo.IsDir()), name) // default FileInfo
		return finfo, nil
	} else {
		if !exists {
			return nil, &os.PathError{"Stat", name, bfs.err}
		} else {
			fmt.Printf("\t⚡ %s %c %s\n", operation, xlateIsDir(finfo.IsDir()), name)
			return finfo, nil
		}
	}
}

// findNode() locates the ynode and returns it, else returns nil. The 2nd
// return value is the cleaned 'name' with expanded ~/ if present.
func (bfs *BitBucketFS) findNode(name string) (iBucketNode, string) {
	bfs.mutex.RLock()
	defer bfs.mutex.RUnlock()

	name = vfs.CleanPath(name)
	node, exists := bfs.entries[name]
	if !exists {
		node = nil
	}
	return node, name
}

// dereferenceSymLink proxies on Filesystem calls where we need to follow
// symbolic links until finding the real target. When successful it returns
// the target ynode and nil error and iBucketNode.ClientData() retrieves the
// actual path/name of the end/target node after dereferencing.
func (bfs *BitBucketFS) dereferenceSymLink(operation, name string) (iBucketNode, error) {
	fmt.Printf("\t%c %s\n", cLINK, name)

	node, cname := bfs.findNode(name)
	if node == nil { // broken symbolic link
		return nil, &os.PathError{Op: operation, Path: cname, Err: os.ErrNotExist}
	}

	if node.IsLink() {
		var err error
		node, err = bfs.dereferenceSymLink(operation, node.Target())
		if err != nil {
			return nil, &os.PathError{Op: operation, Path: cname, Err: os.ErrNotExist}
		}
	} else {
		node.WithClientData(cname)
	}

	return node, nil
}

// execute a single-parameter filesystem operation. In silent mode only the
// action is printed. In non-silent IF the filesystem object is named
// in the corresponding WithFakeDirectories() or WithFakeFiles() then the
// action is printed with the type of object, and if not named therein
// an os.PathError{} is returned.
// Operations: Remove
// Errors: os.PathError
func (bfs *BitBucketFS) executePath(op, name string) error {
	bfs.mutex.Lock()
	defer bfs.mutex.Unlock()

	var errx error = nil
	name = filepath.Clean(name)
	if bfs.silent {
		fmt.Printf("\t⚡ %s %s\n", op, name)
	} else {
		if node, exists := bfs.entries[name]; exists {
			fmt.Printf("\t⚡ %s %c %s\n", op, xlateIsDir(node.IsDir()), name)
			delete(bfs.entries, name)
		} else {
			errx = &os.PathError{op, name, bfs.err}
		}
	}
	return errx
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func xlateIsDir(b bool) rune {
	if b {
		return cDIR
	}
	return cFILE
}
