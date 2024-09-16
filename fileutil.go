/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package vfs

import (
	"os"
	"path/filepath"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// CleanPath trims leading/trailing white space, does filepath.Clean()
// and expands ~/ to the Home Directory of the current user.
func CleanPath(path string) string {
	path = strings.Trim(path, " \t")
	path = filepath.Clean(path)
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return path
}

// HasFileModeFlag checks if the 'flag' is set in the 'flags' value.
func HasModeFlag(flag int, flags int) bool {
	return flags&flag == flag
}

// HasFileModeFlag checks if the 'flag' is set in the 'flags' value.
// Example: HasFileModeFlag(os.ModeSymlink, os.ModeSymLink|0666) returns true
func HasFileModeFlag(flag os.FileMode, flags os.FileMode) bool {
	return flags&flag == flag
}
