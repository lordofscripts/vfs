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

	. "github.com/lordofscripts/vfs/test"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

func Test_FileNodeLite(t *testing.T) {
	const (
		FSIZE int64       = 1024
		FMODE os.FileMode = 0765
	)
	adm := NewUnitTestFramer("BitBucketFS Lite File Node interface", t)
	defer adm.TestCaseFrame(t)(t)

	tests := []struct {
		object iBucketNode
		mode   os.FileMode // Lite node always returns ALL_PERMS, keeps no mode
		dir    bool
		size   int64
		target string
		client any
		want   error
	}{
		// with Regular file
		{newBucketNodeLite(FMODE), ALL_PERMS, false, 0, "", nil, nil},
		{newBucketNodeLite(FMODE), ALL_PERMS, false, 0, "", "data", nil},
		{newBucketNodeLite(FMODE), ALL_PERMS, false, 0, "", 45, nil},
		{newBucketNodeLite(os.ModeSymlink | FMODE), os.ModeSymlink | ALL_PERMS, false, 0, "/usr/sbin/shutdown", nil, nil},
		{newBucketNodeLite(os.ModeSymlink | FMODE), os.ModeSymlink | ALL_PERMS, false, 0, "/usr/sbin/shutdown", true, nil},
		// with Directory
		{newBucketNodeLite(os.ModeDir | FMODE), os.ModeDir | ALL_PERMS, true, 0, "", nil, nil},
		{newBucketNodeLite(os.ModeDir | FMODE), os.ModeDir | ALL_PERMS, true, 0, "", "data", nil},
		{newBucketNodeLite(os.ModeDir | FMODE), os.ModeDir | ALL_PERMS, true, 0, "", 45, nil},
		{newBucketNodeLite(os.ModeDir | FMODE), os.ModeSymlink | os.ModeDir | ALL_PERMS, true, 0, "/usr/sbin", nil, nil},
		{newBucketNodeLite(os.ModeDir | FMODE), os.ModeSymlink | os.ModeDir | ALL_PERMS, true, 0, "/usr/sbin", true, nil},
	}

	for nr, tc := range tests {
		// setup
		node := tc.object
		if tc.client != nil { // setup client data
			node.WithClientData(tc.client)
		}
		if tc.target != "" { // setup symlink
			node.WithLink(tc.target)
		}

		// validate
		if err := verifyFileNode(nr, node, tc.dir, tc.size, tc.mode, tc.target, tc.client); err != nil {
			t.Error(adm.CryE(nil, err))
		}
	}
}

func Test_FileNode(t *testing.T) {
	const (
		FSIZE int64       = 1024
		FMODE os.FileMode = 0765
	)
	adm := NewUnitTestFramer("BitBucketFS Extended File Node interface", t)
	defer adm.TestCaseFrame(t)(t)

	tests := []struct {
		object iBucketNode
		mode   os.FileMode
		dir    bool
		size   int64
		target string
		client any
		want   error
	}{
		// with Regular file
		{newBucketNode(FMODE), FMODE, false, 0, "", nil, nil},
		{newBucketNode(FMODE), FMODE, false, 0, "", "data", nil},
		{newBucketNode(FMODE), FMODE, false, 0, "", 45, nil},
		{newBucketNode(FMODE), os.ModeSymlink | FMODE, false, 0, "/usr/sbin/reboot", nil, nil},
		{newBucketNode(FMODE), os.ModeSymlink | FMODE, false, 0, "/usr/sbin/reboot", true, nil},
		// with Directory
		{newBucketNode(os.ModeDir | FMODE), os.ModeDir | FMODE, true, 0, "", nil, nil},
		{newBucketNode(os.ModeDir | FMODE), os.ModeDir | FMODE, true, 0, "", "data", nil},
		{newBucketNode(os.ModeDir | FMODE), os.ModeDir | FMODE, true, 0, "", 45, nil},
		{newBucketNode(os.ModeDir | FMODE), os.ModeSymlink | os.ModeDir | FMODE, true, 0, "/usr/sbin", nil, nil},
		{newBucketNode(os.ModeDir | FMODE), os.ModeSymlink | os.ModeDir | FMODE, true, 0, "/usr/sbin", true, nil},
	}
	for nr, tc := range tests {
		// setup
		node := tc.object
		if tc.client != nil { // setup client data
			node.WithClientData(tc.client)
		}
		if tc.target != "" { // setup symlink
			node.WithLink(tc.target)
		}

		fmt.Println(node)
		// validate
		if err := verifyFileNode(nr, node, tc.dir, tc.size, tc.mode, tc.target, tc.client); err != nil {
			t.Error(adm.CryE(nil, err))
		}
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

// checks all the public methods of os.FileInfo
func verifyFileNode(nr int, node iBucketNode, isDir bool, size int64, mode os.FileMode, target string, client any) error {
	const TEMPLATE = "case #%d expected %v got %v"
	if node.IsDir() != isDir {
		return fmt.Errorf("IsDir: "+TEMPLATE, nr, isDir, node.IsDir())
	}
	if node.Mode() != mode {
		return fmt.Errorf("Mode: "+TEMPLATE, nr, mode, node.Mode())
	}
	perm := mode & os.ModePerm
	if node.Perms() != perm {
		return fmt.Errorf("Perm: "+TEMPLATE, nr, perm, node.Perms())
	}
	if node.Size() != size {
		return fmt.Errorf("Size: "+TEMPLATE, nr, size, node.Size())
	}
	if node.Target() != target {
		return fmt.Errorf("Target: "+TEMPLATE, nr, target, node.Target())
	}
	if target != "" {
		if node.IsLink() != true {
			return fmt.Errorf("IsLink: "+TEMPLATE, nr, true, node.IsLink())
		}
	}
	if client != nil && node.ClientData() != client {
		return fmt.Errorf("ClientData: "+TEMPLATE, nr, client, node.ClientData())
	}
	return nil
}
