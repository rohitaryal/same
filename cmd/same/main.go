package main

import (
	"flag"
	"fmt"

	"github.com/rohitaryal/same/pkg/hasher"
	"github.com/rohitaryal/same/pkg/scanner"
)

var (
	// Specifies if we should make a backup file
	// containing file hashes from given directory
	backupMode bool = true

	// Specifies if we should start integrity checkup
	// using a backup file
	checkupMode bool = false

	// Path to backup file
	backupFile = ""

	// Path that needs to be backed-up
	backupDirectory = ""

	// Hash to use
	hashMode string = "MD5"

	// Verbose mode = Extra information displayed
	verbose bool = false
)

func main() {
	flag.BoolVar(&backupMode, "b", true, "Initiates backup mode of a directory")
	flag.BoolVar(&checkupMode, "c", false, "Initiates checkup mode using a backup file")
	flag.StringVar(&backupFile, "file", "", "Path to backup file")
	flag.StringVar(&backupDirectory, "dir", "", "Directory that needs to be backed up")
	flag.StringVar(&hashMode, "hash", "MD5", "Hash to use for integrity check")
	flag.BoolVar(&verbose, "v", false, "Make the operation verbose")

	flag.Usage = func() {
		fmt.Println("same: File integrity checker")
		flag.PrintDefaults()
		fmt.Println("\nAuthor: @rohitaryal :)")
	}

	// Dont' forget this one
	flag.Parse()

	if backupMode {
		fmt.Println("Backing up: ", backupDirectory)
		backup()
	}

	if checkupMode {
		fmt.Print("Hello its checkup time.")
	}
}

func backup() {
	channel := make(chan *scanner.File, 100)

	var err error
	var head scanner.File

	go func() {
		head = scanner.Scan(backupDirectory, channel)
	}()

	scanned := 0
	errored := 0
	for file := range channel {
		// Can't hash a file
		if file.IsDirectory {
			continue
		}

		scanned += 1
		if file.Errored {
			errored += 1
		} else {
			file.Remarks, err = hasher.Hash(file.FullPath, hashMode)
			if err != nil {
				errored += 1
			}
		}

		fmt.Printf("\r [+] Total Scanned: %d   Saved: %d   Failure: %d", scanned, scanned-errored, errored)
	}

	fmt.Println("\nHead: ", head.Size)

	fmt.Println()
}
