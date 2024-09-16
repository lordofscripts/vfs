/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Library of Test Utilities
 *-----------------------------------------------------------------*/
package test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ColorBlack       Color = "\u001b[30m"
	ColorRed               = "\u001b[31m"
	ColorLightRed          = "\u001b[91m"
	ColorGreen             = "\u001b[32m"
	ColorBrown             = "\u001b[33m"
	ColorYellow            = "\u001b[93m"
	ColorPurple            = "\u001b[35m"
	ColorLightPurple       = "\u001b[95m"
	ColorBlue              = "\u001b[34m"
	ColorMagenta           = "\u001b[35m"
	ColorCyan              = "\u001b[36m"
	SlowBlink              = "\u001b[5m"
	BlinkOff               = "\u001b[25m"
	ColorReset             = "\u001b[0m"
)

var (
	cliVerbose      bool = false
	cliVerboseLevel int  = 0 // 0=NONE, 1=SOME, 2=ALL
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

// Test module initialization. This allows us to pass parameters
// to the Unit Testing, for example:
//
//	clear; VERBOSE=2 go test github.com/user/project
//	go test ./... -silent
func init() {
	levelStr, isSet := os.LookupEnv("VERBOSE")
	if isSet {
		if level, err := strconv.Atoi(levelStr); err == nil {
			cliVerboseLevel = level
			cliVerbose = (cliVerboseLevel > 0)
		}
	} else {
		cliVerbose = false
		cliVerboseLevel = 0
	}
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Color string

type UnitTestFramer struct {
	t     *testing.T
	name  string
	tally *testCase
	multi bool
}

type testCase struct {
	failed bool
	count  int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewUnitTestFramer(testName string, t *testing.T) *UnitTestFramer {
	return &UnitTestFramer{t, testName, newTestCase(), false}
}

func newTestCase() *testCase {
	return &testCase{false, 0}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (t testCase) Clear() {
	t.failed = false
	t.count = 0
}

// AsMultiple causes the setup title to become indented as a sub-testcase
func (u *UnitTestFramer) AsMultiple() *UnitTestFramer {
	u.multi = true
	return u
}

func (u *UnitTestFramer) Title(title string) *UnitTestFramer {
	if cliVerbose {
		fmt.Println(ColorBrown, "\t§ ", title, ColorReset)
	}
	return u
}

func (u *UnitTestFramer) TestCaseFrame(tb testing.TB) func(tb testing.TB) {
	if cliVerbose {
		var template = "%s➤ %s%s\n"
		if u.multi {
			template = "\t%s§ %s%s\n"
		}
		fmt.Printf(template, ColorBrown, u.name, ColorReset)
	}

	return func(tb testing.TB) {
		fmt.Print(ColorYellow)
		u.Outcome()
		fmt.Println(ColorReset)
	}
}

func (u *UnitTestFramer) Outcome() {
	if !u.tally.failed {
		fmt.Printf("\t* ✔ OK")
	} else {
		fmt.Printf("\t* ✘ FAILED %d subcases\n", u.tally.count)
	}
}

func (u *UnitTestFramer) Cry(format string, v ...any) error {
	u.tally.failed = true
	u.tally.count += 1
	return fmt.Errorf("MSG::"+format, v...)
}

func (u *UnitTestFramer) CryV(expected, got any, format string, v ...any) error {
	u.tally.failed = true
	u.tally.count += 1
	strCmp := fmt.Sprintf("\n\t:: Expected %s but Got %s", expected, got)
	return fmt.Errorf("MSG::"+format+strCmp, v...)
}

func (u *UnitTestFramer) CryE(expected, got error) error {
	u.tally.failed = true
	u.tally.count += 1
	return fmt.Errorf("Expected error %T but got %T > %s", expected, got, got)
}

// Clear resets all the counters so that the instance can be reused in another
// test (sub)case.
func (u *UnitTestFramer) Clear() {
	u.tally.Clear()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/*func setup() {
	fmt.Println("-------------- Unit Tests ---------------")
}

func teardown() {
	fmt.Println("---------------- The End ----------------")
}

func TestMain(m *testing.M) {
    setup()
    defer teardown()

    code := m.Run()
    os.Exit(code)
}*/

func Group(title string) {
	if cliVerbose {
		fmt.Println(ColorCyan, ">> ", title, " <<", ColorReset)
	}
}

func Title(title string) {
	if cliVerbose {
		fmt.Println(ColorBrown, "➤ ", title, ColorReset)
	}
}

func SubTitle(title string) {
	if cliVerbose {
		fmt.Println(ColorBrown, "\t§ ", title, ColorReset)
	}
}

func Verbose(format string, args ...any) {
	if cliVerbose {
		fmt.Printf("\t"+format+"\n", args...)
	}
}

// compares error types, the 1st must NOT be nil
func IsSameErrorType(err1, err2 error) bool {
	//_, ok := err1.(reflect.TypeOf(err2))
	return reflect.TypeOf(err1) == reflect.TypeOf(err2)
}
