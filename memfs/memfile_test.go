package memfs

import (
	"testing"

	"github.com/lordofscripts/vfs"
)

func TestFileInterface(t *testing.T) {
	_ = vfs.File(NewMemFile("", nil, nil))
}
