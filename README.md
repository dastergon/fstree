# fstree

fstree is a recursive directory listing tool that generates a depth indented listing of files and sub-directories in a tree-like format.

fstree is an implementation of the [tree(1)](https://linux.die.net/man/1/tree) command line utility in Go.

## Installation

    go get -u github.com/dastergon/fstree

## Example
Execute the following command to the current working directory.

```bash
$ fstree .
```

Expected output:
```
.
├── LICENSE
├── README.md
└── fstree.go

0 directories, 3 files
```

## Usage

```
Usage: fstree [options...] <target>
  -L int Max display depth of the directory tree. (default -1)
  -a	All files are printed (included hidden files).
  -d	List directories only.
  -f	Prints the full path prefix for each file.
  -filelimit int Do not descend directories that contain more than # entries. (default -1)
  -o string Send output to filename.
  -p	Print the file type and permissions for each file.
  -version Outputs the version of fstree.
```
