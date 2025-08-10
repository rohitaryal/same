// Package backup starts backup operation on specified directory
package backup

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/rohitaryal/same/pkg/hasher"
	"github.com/rohitaryal/same/pkg/logger"
	"github.com/rohitaryal/same/pkg/scanner"
)

func Init(directory, backupFileLocation, hashMode string) {
	channel := make(chan *scanner.File, 100)

	backup, err := os.Create(backupFileLocation)
	if err != nil {
		logger.Error(fmt.Sprint("Failed to create backup file: ", backupFileLocation), err)
		return
	}

	defer backup.Close()

	var head scanner.File

	go func() {
		head = scanner.Scan(directory, channel)
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
		fmt.Print("\033[2K\r")
		fmt.Printf("\r%s Total Scanned: %d   Saved: %d   Failure: %d", logger.Color["LOADING"], scanned, scanned-errored, errored)
	}

	fmt.Printf("\r%s Total Scanned: %d   Saved: %d   Failure: %d", logger.Color["SUCCESS"], scanned, scanned-errored, errored)

	encoder := gob.NewEncoder(backup)
	if err = encoder.Encode(head); err != nil {
		backup.Close()
		fmt.Printf("\n%s Failed to write to backup file.\n", logger.Color["ERROR"])
		return
	}

	fmt.Printf("\n%s Hash backup saved to: %s\n", logger.Color["SUCCESS"], backupFileLocation)
}
