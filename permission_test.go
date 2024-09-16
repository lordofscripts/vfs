/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package vfs

import (
	"fmt"
	"os"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/
func Test_PrimaryPermission(t *testing.T) {
	fmt.Println("➤ Primary Permission value assertion")
	if os.O_RDONLY > 3 || os.O_WRONLY > 3 || os.O_RDWR > 3 {
		t.Errorf("\t✘ Primary permission will not work!")
	} else {
		fmt.Println("\t✔ OK")
	}
}

func Test_PrimaryMode(t *testing.T) {
	fmt.Println("➤ Permission.PrimaryMode()")

	p := Permission(os.O_RDONLY | os.O_APPEND | os.O_TRUNC)
	primary := p.PrimaryMode()
	if primary != os.O_RDONLY {
		t.Errorf("\t✘ Expected O_RDONLY (%x) but got %x", os.O_RDONLY, primary)
	} else {
		fmt.Println("\t✔ OK")
	}
}

func Test_Permission(t *testing.T) {
	fmt.Println("➤ Permission enum")
	type subCaseT struct {
		P Permission
		S string
	}
	subCases := []subCaseT{
		{Permission(os.O_RDONLY | os.O_APPEND), "r--a"},
		{Permission(os.O_WRONLY | os.O_APPEND), "-w-a"},
		{Permission(os.O_RDWR | os.O_APPEND), "rw-a"},
		{Permission(os.O_RDONLY | os.O_CREATE), "r--c"},
		{Permission(os.O_WRONLY | os.O_CREATE), "-w-c"},
		{Permission(os.O_RDWR | os.O_CREATE), "rw-c"},
		{Permission(os.O_RDONLY | os.O_CREATE | os.O_EXCL), "r--ce"},
		{Permission(os.O_WRONLY | os.O_CREATE | os.O_EXCL), "-w-ce"},
		{Permission(os.O_RDWR | os.O_CREATE | os.O_EXCL), "rw-ce"},
		{Permission(os.O_RDONLY | os.O_CREATE | os.O_TRUNC), "r--ct"},
		{Permission(os.O_WRONLY | os.O_CREATE | os.O_TRUNC), "-w-ct"},
		{Permission(os.O_RDWR | os.O_CREATE | os.O_TRUNC), "rw-ct"},
		{Permission(os.O_RDONLY | os.O_EXCL), "r--e"},
		{Permission(os.O_WRONLY | os.O_EXCL), "-w-e"},
		{Permission(os.O_RDWR | os.O_EXCL), "rw-e"},
		{Permission(os.O_RDONLY | os.O_SYNC), "r--s"},
		{Permission(os.O_WRONLY | os.O_SYNC), "-w-s"},
		{Permission(os.O_RDWR | os.O_SYNC), "rw-s"},
		{Permission(os.O_RDONLY | os.O_TRUNC), "r--t"},
		{Permission(os.O_WRONLY | os.O_TRUNC), "-w-t"},
		{Permission(os.O_RDWR | os.O_TRUNC), "rw-t"},
	}
	for _, sub := range subCases {
		if sub.P.String() != sub.S {
			t.Errorf("\t✘ Expected %q but got %q", sub.S, sub.P)
		} else {
			fmt.Printf("\t✔ Permission %d is %s\n", sub.P, sub.S)
		}
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/
