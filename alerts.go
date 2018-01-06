package main

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
)

//
// Review the size of the backup file and make sure
// it is within an expected range.
//
func AlertsBackupSize(filePath string, low int64, high int64) error {

	// Make sure this is really a file.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file stats
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	bytes := stat.Size()

	// Log size
	Log("Backup size is " + humanize.Bytes(uint64(bytes)) + ".")

	// Make sure the backup is bigger than the low end.
	if stat.Size() < low {
		fmt.Println("Failed Low test")
	}

	// Make sure the backup is less than the high end.
	if stat.Size() > high {
		fmt.Println("Failed High test")
	}

	// Return happy.
	return nil
}

/* End File */
