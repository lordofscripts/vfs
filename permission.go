/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package vfs

import (
	"os"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Permission int

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (pe Permission) String() string {
	var str string
	p := int(pe)
	primary := pe.PrimaryMode()

	if primary == os.O_RDWR { // O_RDWR   int = syscall.O_RDWR open the file read-write
		str = "rw-"
	} else if primary == os.O_RDONLY { // O_RDONLY int = syscall.O_RDONLY open the file read-only
		str = "r--"
	} else if primary == os.O_WRONLY { // O_WRONLY int = syscall.O_WRONLY open the file write-only
		str = "-w-"
	} else {
		panic("Couldn't figure out Permission")
	}

	// These values can be ORed together
	if p&os.O_APPEND != 0 { // O_APPEND int = syscall.O_APPEND append data to the file when writing
		str += "a"
	}
	if p&os.O_CREATE != 0 { // O_CREATE int = syscall.O_CREAT create a new file if none exists
		str += "c"
	}
	if p&os.O_EXCL != 0 { // O_EXCL   int = syscall.O_EXCL used with O_CREATE, file must not exist
		str += "e"
	}
	if p&os.O_SYNC != 0 { // O_SYNC   int = syscall.O_SYNC open for synchronous I/O
		str += "s"
	}
	if p&os.O_TRUNC != 0 { // O_TRUNC  int = syscall.O_TRUNC truncate regular writable file when opened
		str += "t"
	}

	return str
}

// PrimaryMode returns only the permission part corresponding to either of:
// os.O_RDWR, os.O_RDONLY or os.O_WRONLY
func (pe Permission) PrimaryMode() int {
	const PRIM int = 0x3
	return int(pe) & PRIM
}

// Flags returns the permission value with its PrimaryMode bits set to zero.
func (pe Permission) Flags() int {
	const PRIM uint = 0x3
	return int(uint(pe) & ^PRIM)
}
