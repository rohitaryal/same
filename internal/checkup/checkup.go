// Package checkup provides method to initiate integrity
// checkup using a backup file
package checkup

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rohitaryal/same/pkg/hasher"
	"github.com/rohitaryal/same/pkg/logger"
	"github.com/rohitaryal/same/pkg/scanner"
)

func Init(backupFileLocation, hashMode string) {
	file, err := os.Open(backupFileLocation)
	if err != nil {
		logger.Error(fmt.Sprint("Failed to open backup file: ", backupFileLocation), err)
		return
	}

	decoder := gob.NewDecoder(file)
	var loadedRoot scanner.File
	if err = decoder.Decode(&loadedRoot); err != nil {
		logger.Error(fmt.Sprint("Failed to read from backup file: ", backupFileLocation), err)
		return
	}

	success := 0
	errored := 0

	channel := make(chan *scanner.File)

	go func() {
		nestedCheck(&loadedRoot, channel, hashMode)

		close(channel)
	}()

	var invalidFiles []*scanner.File

	for file := range channel {
		if file.Errored {
			errored += 1
			invalidFiles = append(invalidFiles, file)
		} else {
			success += 1
		}

		fmt.Print("\033[2K\r")
		fmt.Printf("\r%s Valid Integrity: %d \t Invalid Integrity: %d", logger.Color["LOADING"], success, errored)
	}

	fmt.Printf("\r%s Valid Integrity: %d \t Invalid Integrity: %d", logger.Color["SUCCESS"], success, errored)

	if len(invalidFiles) > 0 {
		fmt.Printf("\n%s %d files with compromised integrity:", logger.Color["WARNING"], len(invalidFiles))

		for _, file := range invalidFiles {
			fmt.Printf("\n\t%s %s [%s]", color.BlueString("-"), file.FullPath, file.Remarks)
		}
	}

	fmt.Println()
}

func nestedCheck(file *scanner.File, channel chan *scanner.File, hashMode string) {
	if file.IsDirectory {
		for _, child := range file.Contents {
			nestedCheck(child, channel, hashMode)
		}
		return
	}

	res, err := hasher.CompareHash(file.FullPath, file.Remarks, hashMode)
	// Yes always forgot to invert
	file.Errored = !res

	if err != nil {
		if os.IsNotExist(err) {
			file.Remarks = "DELETED"
		} else if os.IsPermission(err) {
			file.Remarks = "DENIED"
		} else {
			file.Remarks = "UNKNOWN"
		}
	} else if !res {
		file.Remarks = "INVALID_HASH"
	}

	channel <- file
}
