package normalize

import "golang.org/x/sys/unix"

// renameNew renames from to to and fails when to exists.
func renameNew(from, to string) (string, error) {
	return to, unix.RenamexNp(from, to, unix.RENAME_EXCL)
}
