package memfs

import (
	"testing"

	"github.com/3JoB/vfs"
)

func TestFileInterface(t *testing.T) {
	_ = vfs.File(NewMemFile("", nil, nil))
}
