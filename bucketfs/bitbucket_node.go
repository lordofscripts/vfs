/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A BITBUCKET FILESYSTEM NODE (LITE & NORMAL)
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"fmt"
	"os"
	"strings"

	"github.com/lordofscripts/vfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	DEF_USER  string = "boot"
	DEF_GROUP string = "foot"

	ALL_PERMS    os.FileMode = 0777
	ALL_RW_PERMS os.FileMode = 0666
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var (
	_ iBucketNode = (*bucketNodeLite)(nil)
	_ iBucketNode = (*bucketNode)(nil)
)

// iBucketNode is an internal interface that allows us to compile
// BitBucketFS with minimal node information (bucketNodeLite) or slightly
// more functional (bucketNode).
type iBucketNode interface {
	// setup supplementary client data
	WithClientData(data any) iBucketNode
	// setup as symbolic link
	WithLink(name string) iBucketNode
	// Is it a directory?
	IsDir() bool
	// Is it a symbolic link?
	IsLink() bool
	// file size
	Size() int64
	// the full file mode (permissions/fs.ModePerm & type/fs.ModeType)
	Mode() os.FileMode
	// file permissions (masked by fs.ModePerm)
	Perms() os.FileMode
	// target file/dir (only if symbolic link)
	Target() string
	// Create a symbolic link representation node. It sets os.ModeSymLink in Mode()
	LinkTo(string) iBucketNode
	// Get client data (usually nil)
	ClientData() any
	// Stringer shows Perms() as string. For symbolic links also what it points to.
	String() string
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// BucketNodeLite is the minimal (low overhead) version of a BitBucket
// filesystem node which only has notion of whether the filesystem object
// is a directory or not.
type bucketNodeLite struct {
	isDir  bool
	isLink bool
	target string
	extra  any
}

// BucketNode is a slightly more functional (yet minimal) node information
// that remembers a filesystem object's type (directory, device, link, etc.)
// and (Unix) permissions and its size.
type bucketNode struct {
	fmode  os.FileMode
	fsize  int64
	target string
	extra  any
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// newBucketNodeLite is a light version of a BitBucketFS node
func newBucketNodeLite(mode os.FileMode) *bucketNodeLite {
	return &bucketNodeLite{vfs.HasFileModeFlag(os.ModeDir, mode),
		vfs.HasFileModeFlag(os.ModeSymlink, mode),
		"", nil}
}

// newBucketNodeLite is a more functional version of a BitBucketFS node
func newBucketNode(mode os.FileMode) *bucketNode {
	return &bucketNode{mode, 0, "", nil}
}

// createNode creates a BitBucketFS node. The developer chooses to select
// therein whether it is a lite version or a functional version.
func createNode(mode os.FileMode) iBucketNode {
	var node iBucketNode
	// Uncomment the next line to use lite nodes
	//node = newBucketNodeLite(mode)
	// Uncomment the next line to use functional nodes instead.
	node = newBucketNode(mode)
	return node
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (min *bucketNodeLite) WithClientData(data any) iBucketNode {
	min.extra = data
	return min
}

// setup as symbolic link
func (min *bucketNodeLite) WithLink(name string) iBucketNode {
	if len(name) > 0 {
		min.target = name
		min.isLink = true
	}
	return min
}

func (min *bucketNodeLite) IsDir() bool {
	return min.isDir
}

func (min *bucketNodeLite) IsLink() bool {
	return min.isLink && len(min.target) > 0
}

func (min *bucketNodeLite) Size() int64 {
	return 0
}

func (min *bucketNodeLite) Mode() os.FileMode {
	mode := ALL_PERMS
	if min.isDir {
		mode = os.ModeDir | mode
	}
	if min.isLink {
		mode = os.ModeSymlink | mode
	}
	return mode
}

func (min *bucketNodeLite) Perms() os.FileMode {
	return ALL_PERMS
}

func (min *bucketNodeLite) Target() string {
	return min.target
}

func (min *bucketNodeLite) LinkTo(key string) iBucketNode {
	node := newBucketNodeLite(os.ModeSymlink | ALL_PERMS)
	node.target = key
	return node
}

func (min *bucketNodeLite) ClientData() any {
	return min.extra
}

func (min *bucketNodeLite) String() string {
	var mode os.FileMode
	if min.isDir {
		mode = os.ModeDir | ALL_PERMS
	}
	if min.isLink {
		mode = os.ModeSymlink | ALL_PERMS
	}
	return fmt.Sprintf("%s", mode)
}

func (b *bucketNode) WithClientData(data any) iBucketNode {
	b.extra = data
	return b
}

// setup as symbolic link
func (b *bucketNode) WithLink(name string) iBucketNode {
	if len(name) > 0 {
		b.target = name
		b.fmode = os.ModeSymlink | b.fmode
	}
	return b
}

func (b *bucketNode) IsDir() bool {
	return vfs.HasFileModeFlag(os.ModeDir, b.fmode)
}

func (b *bucketNode) IsLink() bool {
	return vfs.HasFileModeFlag(os.ModeSymlink, b.fmode)
}

func (b *bucketNode) Size() int64 {
	return b.fsize
}

func (b *bucketNode) Mode() os.FileMode {
	return b.fmode
}

func (b *bucketNode) Perms() os.FileMode {
	return b.fmode & os.ModePerm
}

func (b *bucketNode) Target() string {
	return b.target
}

func (b *bucketNode) LinkTo(key string) iBucketNode {
	var n *bucketNode = nil
	if len(strings.Trim(key, " \t")) > 0 {
		n = newBucketNode(os.ModeSymlink | ALL_PERMS)
		n.fsize = 0
		n.target = key
	}
	return n
}

func (b *bucketNode) ClientData() any {
	return b.extra
}

func (b *bucketNode) String() string {
	targetStr := ""
	if b.IsLink() && len(b.target) > 0 {
		targetStr = " ->" + b.target
	}
	return fmt.Sprintf("%s%s", b.fmode, targetStr)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
