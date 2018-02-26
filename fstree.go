package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xlab/treeprint"
)

const (
	version = "1.0"
)

var (
	// -a : All files are printed (by default hidden files not printed)
	flagPrintAllFiles = flag.Bool("a", false, "All files are printed (included hidden files).")
	// -d : List directories only
	flagDirsOnly = flag.Bool("d", false, "List directories only.")
	// -f: Prints the full path prefix for each file
	flagFullPathPrefix = flag.Bool("f", false, "Prints the full path prefix for each file.")
	//--filelimit :  Do not descend directories that contain more than # entries.
	flagFileLimit = flag.Int("filelimit", -1, "Do not descend directories that contain more than # entries.")
	// -L : Max display depth of the directory tree.
	flagDirLimit = flag.Int("L", -1, "Max display depth of the directory tree.")
	// -o : Send output to filename
	flagWriteToFile = flag.String("o", "", "Send output to filename.")
	// -p : Print the file type and permissions for each file (as per ls -l).
	flagPermissions = flag.Bool("p", false, "Print the file type and permissions for each file.")
	// -version: Outputs the version of fstree
	flagVersion = flag.Bool("version", false, "Outputs the version of fstree.")

	directoriesCount = 0
	filesCount       = 0
)

// traverseFileSystem traverses the filesystem in a depth-first fashion and
// builds a filesystem hierarchy in a tree-like format.
func traverseFilesystem(currentPath string, root bool, tree treeprint.Tree, depth int) {
	var outputName string

	// Add a new branch to the tree.
	if !root {
		outputName = filepath.Base(currentPath)
		if *flagFullPathPrefix {
			outputName = currentPath
		}

		if *flagPermissions {
			fileInfo, _ := os.Stat(currentPath)
			tree = tree.AddMetaBranch(fileInfo.Mode(), outputName)
		} else {
			tree = tree.AddBranch(outputName)
		}

	}
	// Iterate through the files of the directory
	files, _ := ioutil.ReadDir(currentPath)
	for _, f := range files {
		// Check the depth level and stop when -L flag is used.
		if depth == *flagDirLimit {
			continue
		}

		// Do not descend directories that contain more than # entries when
		// --filelimit flag is used.
		if *flagFileLimit >= 0 && len(files) >= *flagFileLimit {
			if currentPath == "." {
				fmt.Printf("fstree: [ %d entries exceeds filelimit, not opening dir]\n", len(files))
			}
			break
		}

		// Check if a regular file or  directory is hidden,
		// and ignore it if flag flagPrintAllFiles is not used.
		if !*flagPrintAllFiles && strings.HasPrefix(f.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(currentPath, f.Name())
		if f.IsDir() {
			directoriesCount++
			traverseFilesystem(fullPath, false, tree, depth+1)
		} else {
			filesCount++
			// Ignore directories if flag -d is used.
			if *flagDirsOnly {
				continue
			}
			// Add a new node with the filename to the current branch.
			outputName = f.Name()
			if *flagFullPathPrefix {
				outputName = fullPath
			}
			if *flagPermissions {
				tree.AddMetaNode(f.Mode(), outputName)
			} else {
				tree.AddNode(outputName)
			}
		}
	}
}

// printFileSystemStats builds a consistent output
// for the directory and file stats.
func printFileSystemStats() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(directoriesCount))
	if directoriesCount == 1 {
		buf.WriteString(" directory, ")
	} else {
		buf.WriteString(" directories, ")
	}
	buf.WriteString(strconv.Itoa(filesCount))
	if filesCount == 1 {
		buf.WriteString(" file")
	} else {
		buf.WriteString(" files")
	}
	return buf.String()
}

func main() {

	// Modfify the default usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <target>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *flagVersion {
		fmt.Println("fstree:", version)
		os.Exit(0)
	}

	if *flagDirLimit == 0 {
		fmt.Println("fstree: Invalid level, must be greater than 0.")
		os.Exit(1)
	}

	flags := flag.Args()
	// If the the user does not provide a target directory,
	// default to "." (current working directory)
	currentPath := "."
	if len(flags) != 0 {
		currentPath = flags[0]
	}
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist.\n", currentPath)
		os.Exit(1)
	}

	// Create a tree object.
	tree := treeprint.New()

	// Traverse the filesystem and build the filesystem tree.
	traverseFilesystem(currentPath, true, tree, 0)

	// Generate the list of directories and files in a tree-like format.
	if *flagWriteToFile != "" {
		// Write the filesystem tree to a file when -o flag is used.
		file, err := os.Create(*flagWriteToFile)
		if err != nil {
			fmt.Printf("Cannot create file: %s\n", *flagWriteToFile)
			os.Exit(1)
		}
		defer file.Close()

		if _, err := file.Write(tree.Bytes()); err != nil {
			fmt.Printf("Cannot write to file: %s\n", *flagWriteToFile)
			os.Exit(1)
		}
		if _, err := file.WriteString(printFileSystemStats()); err != nil {
			fmt.Printf("Cannot write to file: %s\n", *flagWriteToFile)
			os.Exit(1)
		}
	} else {
		// Print the filesystem tree.
		fmt.Println(tree.String())
		// Print the stats from the traversed filesystem.
		fmt.Println(printFileSystemStats())
		os.Exit(0)
	}
}
