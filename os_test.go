package vfs

import (
	"os"
	"testing"
)

func TestOSInterface(t *testing.T) {
	_ = Filesystem(OS())
}

func TestOSCreate(t *testing.T) {
	fs := OS()

	f, err := fs.OpenFile("/tmp/test123", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Errorf("Create: %s", err)
	}
	if err = f.Close(); err != nil {
		t.Errorf("Close: %s", err)
	}
	f2, err := fs.Open("/tmp/test123")
	if err != nil {
		t.Errorf("Open: %s", err)
	}
	if err := f2.Close(); err != nil {
		t.Errorf("Close: %s", err)
	}
	if err := fs.Remove(f.Name()); err != nil {
		t.Errorf("Remove: %s", err)
	}
}
