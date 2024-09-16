/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/lordofscripts/vfs"
	. "github.com/lordofscripts/vfs/test"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

func Test_IFile(t *testing.T) {
	_ = vfs.File(newBitBucketFile("dummy.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666, ErrAny))
}

// File{}.Open() existing file should NOT fail
func Test_Open(t *testing.T) {
	const (
		NAME        = "File.Open()"
		FAKE_FILE   = "/mnt/dummy.txt" // exists
		ABSENT_FILE = "/mnt/savvy.doc" // doesn't exist
	)

	ErrBitBucket := errors.New("Open Error!")
	tests := []struct {
		name     string
		filename string
		want     error
	}{
		{"Open existing", FAKE_FILE, nil},
		{"Open non-existing (PathError)", ABSENT_FILE, &fs.PathError{}},
	}

	Title(NAME)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adm := NewUnitTestFramer(tt.name, t).AsMultiple()
			teardownCase := adm.TestCaseFrame(t)
			defer teardownCase(t)

			fs := CreateWithError(ErrBitBucket).WithFakeFiles([]string{FAKE_FILE})

			if fd, err := fs.Open(tt.filename); !IsSameErrorType(tt.want, err) {
				t.Error(adm.Cry("Open failed %T vs %T %s", err, tt.want, err))
			} else {
				if fd != nil {
					fd.Close()
				}
			}
		})
	}
}

// File{}.Truncate() should truncate the file to specified size
func Test_Truncate(t *testing.T) {
	const NAME = "File.Truncate()"
	t.Run(NAME, func(t *testing.T) {
		adm := NewUnitTestFramer(NAME, t)
		teardownCase := adm.TestCaseFrame(t)
		defer teardownCase(t)

		const FILEOBJ = "/mnt/dummy.txt"
		ErrBitBucket := errors.New("Truncate Error!")
		fs := CreateWithError(ErrBitBucket).WithFakeFiles([]string{FILEOBJ})

		if fd, err := fs.Open(FILEOBJ); err != nil {
			t.Error(adm.Cry("could not open %T %s", err, err))
		} else {
			// TODO: BitBucketFileInfo should keep track of file size
			if err := fd.Truncate(0); err != nil {
				t.Error(adm.Cry("Truncate failed (%T): %w", err, err))
			}
			// TODO: need a way to get file size via Stat()
			fd.Close()
		}
	})
}

// File{}.Name() must return the fully-qualified filename
func Test_Name(t *testing.T) {
	const NAME = "File.Name()"
	t.Run(NAME, func(t *testing.T) {
		adm := NewUnitTestFramer(NAME, t)
		teardownCase := adm.TestCaseFrame(t)
		defer teardownCase(t)

		const FILEOBJ = "/mnt/dummy.txt"
		ErrBitBucket := errors.New("Name Error!")
		fs := CreateWithError(ErrBitBucket).WithFakeFiles([]string{FILEOBJ})

		if fd, err := fs.Open(FILEOBJ); err != nil {
			t.Error(adm.Cry("could not open %T %s", err, err))
		} else {
			got := fd.Name()
			if got != FILEOBJ {
				t.Error(adm.CryV(FILEOBJ, got, "Name mismatch"))
			}
			fd.Close()
		}
	})
}

// File{}.Sync() should always throws error
func Test_Sync(t *testing.T) {
	const (
		NAME      = "File.Sync()"
		FAKE_FILE = "/mnt/dummy.txt"
	)
	t.Run(NAME, func(t *testing.T) {
		adm := NewUnitTestFramer(NAME, t)
		teardownCase := adm.TestCaseFrame(t)
		defer teardownCase(t)

		ErrBitBucket := errors.New("Sync Error!")
		fs := CreateWithError(ErrBitBucket).WithFakeFiles([]string{FAKE_FILE})
		if fd, err := fs.Open(FAKE_FILE); err != nil {
			t.Error(adm.Cry("Open result %T %s", err, err))
		} else {
			if err := fd.Sync(); err != ErrBitBucket {
				t.Error(adm.CryE(ErrBitBucket, err))
			}
			fd.Close()
		}
	})
}

// File{}.Write()
func Test_Write(t *testing.T) {
	const (
		NAME      = "File.Write()"
		FAKE_FILE = "/mnt/dummy.txt"
	)
	t.Run(NAME, func(t *testing.T) {
		adm := NewUnitTestFramer(NAME, t)
		teardownCase := adm.TestCaseFrame(t)
		defer teardownCase(t)

		ErrBitBucket := errors.New("Test Error!")
		fs := CreateWithError(ErrBitBucket).WithFakeFiles([]string{FAKE_FILE})
		if fd, err := fs.Open(FAKE_FILE); err != nil {
			t.Error(adm.Cry("Open result %T %s", err, err))
		} else {
			data := []byte("This is a test")
			n, err := fd.Write(data)
			if n != len(data) {
				t.Error(adm.CryV(len(data), n, "Write length"))
			}
			if err != nil {
				t.Error(adm.CryE(nil, err))
			}
			if fd != nil {
				fd.Close()
			}
		}
	})
}

func Test_Any(t *testing.T) {
	Title("Any Title VERBOSE")
	t.Run("NAME", func(t *testing.T) {
		frame := NewUnitTestFramer("NAMEX", t)
		teardownCase := frame.TestCaseFrame(t)
		defer teardownCase(t)

		fmt.Println("nothing")
	})
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

func setup() {
	Group("-------------- BucketFS.File Unit Tests ---------------")
}

func teardown() {
	fmt.Println("---------------- The End ----------------")
}

func TestMain(m *testing.M) {
	setup()
	defer teardown()

	code := m.Run()
	os.Exit(code)
}
