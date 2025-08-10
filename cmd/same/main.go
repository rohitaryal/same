package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/rohitaryal/same/internal/backup"
	"github.com/rohitaryal/same/internal/checkup"
	"github.com/rohitaryal/same/pkg/logger"
)

func main() {
	isBackup := flag.Bool("init", false, "Initiates hash backup for a path")
	isCheckup := flag.Bool("check", false, "Initiates integrity checkup using backup file")
	backupFileLocation := flag.String("file", "", "Path to save/saved backup file")
	hashMode := flag.String("hash", "md5", "Hash to use for integrity check [md5, size, sha256]")
	directory := flag.String("dir", ".", "Direcory to be backed up")
	flag.Parse()

	if *isCheckup && *backupFileLocation == "" {
		logger.Error("Please provide a backup file location", nil)
		return
	}

	if *backupFileLocation == "" {
		*backupFileLocation = fmt.Sprintf("backup-%d.same", time.Now().UnixMilli())
	}

	if *isBackup {
		logger.Info(fmt.Sprint("Initiating backup for: ", *directory))
		backup.Init(*directory, *backupFileLocation, *hashMode)
	} else if *isCheckup {
		logger.Info(fmt.Sprint("Initiating checkup from: ", *backupFileLocation))
		checkup.Init(*backupFileLocation, *hashMode)
	} else {
		logger.Warning("Get yourself some help :)")
		flag.Usage()
	}
}
