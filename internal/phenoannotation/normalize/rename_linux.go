package normalize

import "golang.org/x/sys/unix"

// renameNew renames from to to and fails when to exists.
func renameNew(from, to string) (string, error) {
	return to, unix.Renameat2(unix.AT_FDCWD, from, unix.AT_FDCWD, to, unix.RENAME_NOREPLACE)
}
