# dirsync
Basic directory synchronization with golang.
Recursively copies files and subdirectories over to target directory, deleting
files not present in source path and preserving identical files without copying.

If two files have a matching name and size, checksums are compared to determine
identicality.

# Usage
Get it: `go get github.com/Varjelus/dirsync`

Import it: `import github.com/Varjelus/dirsync`


## Use it

`err := dirsync.Sync("/path/to/dir1", "/path/to/dir2")`, handle error
