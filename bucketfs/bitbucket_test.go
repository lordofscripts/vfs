/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package bucketfs

import (
	"errors"
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
	CHR_OK          rune = '✔'
	CHR_NOK         rune = '✘'
	CHR_BEAK_RARROW      = '➤'
	CHR_SECTION          = '§'

	// Not in Fake*
	cFILE1 = "/tmp/not-exists/test1.txt"
	cFILE2 = "/tmp/not-exists/test2.pdf"
	cFILE3 = "/tmp/not-exists/deletions/oldname.doc"
	cFILE4 = "/tmp/not-exists/deletions/newname.doc"

	cDIR1 = "/tmp/not-exists"

	// Used in Fake*
	cFAKE_DIR1  = "/tmp/faked/Dir1"
	cFAKE_DIR11 = "/tmp/faked/Dir1/Dir11"
	cFAKE_DIR2  = "/tmp/faked/Dir2"

	cFAKE_FILE1  = "/tmp/faked/file1.txt"
	cFAKE_FILE11 = "/tmp/faked/Dir1/file11.doc"
	cFAKE_FILE2  = "/tmp/faked/Dir2/file2.jpg"

	cCASE_TITLE_TEMPLATE    = "➤ BitBucketFS Filesystem{}.%s()\n"
	cSUBCASE_SILENT         = "\t§ Silent mode"
	cSUBCASE_HYBRID         = "\t§ WithError But Faked"
	cSUBCASE_HYBRID_WITHERR = "\t§ WithError IF NOT Faked"
)

var (
	ErrAny error = errors.New("Any BitBucket error")

	fakeDirs = []string{
		cFAKE_DIR1,
		cFAKE_DIR2,
		cFAKE_DIR11,
	}
	fakeFiles = []string{
		cFAKE_FILE1,
		cFAKE_FILE2,
		cFAKE_FILE11,
	}
)

/* ----------------------------------------------------------------
 *				BitBucketFS  U n i t  T e s t
 *-----------------------------------------------------------------*/
func Test_IFilesystem(t *testing.T) {
	_ = vfs.Filesystem(Create())
}

// vfs.CreateWithError(error) must be called with non-nil error
func Test_BadConstructor(t *testing.T) {
	frame := NewUnitTestFramer("Bad parameter in ctor. » panic", t)
	defer frame.TestCaseFrame(t)(t)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CreateWithError ctor did not panic")
		}
	}()

	// should panic because a default error must be specified
	_ = CreateWithError(nil)
}

// goes through all vfs.Filesystem interface methods with BitBucketFS
// in silent mode, i.e. operation name gets printed but NO error is ever
// produced regardless of whether the file would normally be expected to
// exist or not. All operations should return nil.
func Test_BitBucketFS_Silent(t *testing.T) {
	adm := NewUnitTestFramer("BitBucketFS Silent mode (Dry run)", t)
	defer adm.TestCaseFrame(t)(t)

	fs := Create()
	adm.Title("Open")
	if _, err := fs.Open(cFAKE_FILE1); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("OpenFile")
	if _, err := fs.OpenFile(cFAKE_FILE2, 0, 0); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("Remove")
	if err := fs.Remove(cFAKE_FILE11); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("Rename")
	if err := fs.Rename(cFILE3, cFILE4); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("Mkdir")
	if err := fs.Mkdir(cFAKE_DIR1, 0); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("Stat")
	if _, err := fs.Stat(cFILE1); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("Lstat")
	if _, err := fs.Lstat(cFILE2); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	adm.Title("ReadDir")
	if _, err := fs.ReadDir(cFAKE_DIR1); err != nil {
		t.Error(adm.CryE(nil, err))
	}
}

/* ----------------------------------------------------------------
 *	Filesystem.Open()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Open(t *testing.T) {
	// TODO: open symlink
	const Oper = "Open"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if _, err := fs.Open(cFILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	fs = CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if _, err := fs.Open(cFAKE_FILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.PathError
	if _, err := fs.Open(cFILE1); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s\n\tGot: %T", Oper, err, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Open()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_OpenFile(t *testing.T) {
	const Oper = "OpenFile"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if _, err := fs.OpenFile(cFILE1, os.O_RDWR, 0666); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	fs = CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if _, err := fs.OpenFile(cFAKE_FILE1, os.O_RDWR, 0666); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.PathError
	if _, err := fs.OpenFile(cFILE1, os.O_RDWR, 0666); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Remove()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Remove(t *testing.T) {
	// TODO: remove symlink
	const Oper = "Remove"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if err := fs.Remove(cFILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	fs = CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if err := fs.Remove(cFAKE_FILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.PathError
	if err := fs.Remove(cFILE1); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Rename()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Rename(t *testing.T) {
	// TODO: rename symlink
	const Oper = "Rename"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if err := fs.Rename(cFILE1, cFILE2); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	fs = CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if err := fs.Rename(cFAKE_FILE1, cFILE2); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.LinkError
	if err := fs.Rename(cFILE1, cFILE2); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Symlink()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Symlink(t *testing.T) {
	const Oper = "Symlink"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs1 := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if err := fs1.Symlink(cFILE1, cFILE2); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	var e *os.LinkError
	fs2 := CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if err := fs2.Symlink(cFAKE_FILE1, cFILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	if err := fs2.Symlink(cFILE3, cFILE2); !errors.As(err, &e) {
		t.Errorf("%s expected os.LinkError got: %T %s", Oper, err, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Remove()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Mkdir(t *testing.T) {
	const Oper = "Mkdir"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if err := fs.Mkdir(cFAKE_DIR1, 0777); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	// II. Non-silent
	fs = CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (error because DIR exists as Fake)
	var e *os.PathError
	fmt.Println(cSUBCASE_HYBRID)
	if err := fs.Mkdir(cFAKE_DIR1, 0777); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (no error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	if err := fs.Mkdir(cDIR1, 0777); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Stat()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Stat(t *testing.T) {
	// TODO: Stat symlink should follow
	const Oper = "Stat"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs1 := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if _, err := fs1.Stat(cFILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	fs1 = nil

	// II. Non-silent
	fs2 := CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if _, err := fs2.Stat(cFAKE_FILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.PathError
	if _, err := fs2.Stat(cFILE1); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Stat()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_Lstat(t *testing.T) {
	// TODO: Lstat symlink & non-symlink
	const Oper = "Lstat"
	fmt.Printf(cCASE_TITLE_TEMPLATE, Oper)

	ok := true
	// I. Silent
	fs1 := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if _, err := fs1.Lstat(cFILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	fs1 = nil

	// II. Non-silent
	fs2 := CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1 on faked (no error)
	fmt.Println(cSUBCASE_HYBRID)
	if _, err := fs2.Lstat(cFAKE_FILE1); err != nil {
		t.Errorf("%s nil expected: %s", Oper, err)
		ok = false
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	var e *os.PathError
	if _, err := fs2.Lstat(cFILE1); !errors.As(err, &e) {
		t.Errorf("%s os.PathError expected: %s", Oper, err)
		ok = false
	}

	outcome(ok)
}

/* ----------------------------------------------------------------
 *	Filesystem.Stat()
 *-----------------------------------------------------------------*/
func Test_BitBucketFS_ReadDir(t *testing.T) {
	const Oper = "ReadDir"
	adm := NewUnitTestFramer(Oper, t)
	defer adm.TestCaseFrame(t)(t)

	// I. Silent
	fs1 := Create()
	// 1.1 any FS object
	fmt.Println(cSUBCASE_SILENT)
	if _, err := fs1.ReadDir(cFILE1); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	fs1 = nil

	// II. Non-silent
	fs2 := CreateWithError(ErrAny).
		WithFakeDirectories(fakeDirs).
		WithFakeFiles(fakeFiles)
	// 2.1.1 on faked (error because exists as FILE)
	var e *os.PathError
	fmt.Println(cSUBCASE_HYBRID + " as FILE")
	if _, err := fs2.ReadDir(cFAKE_FILE1); !errors.As(err, &e) {
		t.Error(adm.CryE(nil, err))
	}
	// 2.1.2 on faked (no error because it exists as DIR)
	fmt.Println(cSUBCASE_HYBRID + " as DIR")
	if _, err := fs2.ReadDir(cFAKE_DIR1); err != nil {
		t.Error(adm.CryE(nil, err))
	}
	// 2.2 on non-faked (error)
	fmt.Println(cSUBCASE_HYBRID_WITHERR)
	if _, err := fs2.ReadDir(cFILE1); !errors.As(err, &e) {
		//t.Errorf("%s os.PathError expected: %s", Oper, err)
		t.Error(adm.CryE(e, err))
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

func outcome(ok bool) {
	if ok {
		fmt.Println("\t* ✔ OK")
	} else {
		fmt.Println("\t* ✘ FAILED")
	}
}
