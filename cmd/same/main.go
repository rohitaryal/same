package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/fatih/color"
	"github.com/rohitaryal/same/pkg/hasher"
	"github.com/rohitaryal/same/pkg/scanner"
)

var Color = map[string]string{
	"ERROR":   color.HiRedString("[✕]"),
	"INFO":    color.HiCyanString("[i]"),
	"LOADING": color.HiBlueString("[*]"),
	"SUCCESS": color.HiGreenString("[✓]"),
	"WARNING": color.HiYellowString("[✕]"),
}

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

	flag.StringVar(&backupFile, "file", fmt.Sprint("backup", time.Now().Unix()), "Path to save backup file")

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

	// Append `.same` extension
	if path.Ext(backupFile) != ".same" {
		backupFile += ".same"
	}

	if backupMode {
		fmt.Println("[i] Backing up: ", backupDirectory)
		backup()
	}

	if checkupMode {
		fmt.Println("[]")
	}
}

func backup() {
	channel := make(chan *scanner.File, 100)

	backup, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("%s Failed to create backup file.", Color["ERROR"])
		return
	}

	defer backup.Close()

	var head scanner.File

	go func() {
		head = scanner.Scan(backupDirectory, channel)
	}()

	scanned := 0
	errored := 0
	for file := range channel {
		// Can't hash a directory
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

		fmt.Printf("\r%s Total Scanned: %d   Saved: %d   Failure: %d", Color["LOADING"], scanned, scanned-errored, errored)
	}

	fmt.Printf("\r%s Total Scanned: %d   Saved: %d   Failure: %d", Color["SUCCESS"], scanned, scanned-errored, errored)

	encoder := gob.NewEncoder(backup)
	if err = encoder.Encode(head); err != nil {
		backup.Close()
		fmt.Printf("\n%s Failed to write to backup file.\n", Color["ERROR"])
	}

	fmt.Printf("\n%s Hash backup saved to: %s\n", Color["SUCCESS"], backupFile)
}
