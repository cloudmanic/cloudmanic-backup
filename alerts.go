package main

import (
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

		// Log alert
		Log("Failed Alert = Database backup was less than " + humanize.Bytes(uint64(low)))

		// Setup Email
		subject := "[Cloudmanic Backup] The " + cfg.BackupName + " Database Backup Was Too Small In Size."
		text := "The " + cfg.BackupName + " database backup was less than " + humanize.Bytes(uint64(low)) + " in size. It was " + humanize.Bytes(uint64(bytes)) + "."
		html := "<p>" + text + "</p>"

		// Send Email
		EmailSend(cfg.AlertEmail, subject, html, text)
	}

	// Make sure the backup is less than the high end.
	if stat.Size() > high {

		// Log alert
		Log("Failed Alert - Database backup was bigger than " + humanize.Bytes(uint64(high)))

		// Setup Email
		subject := "[Cloudmanic Backup] The " + cfg.BackupName + " Database Backup Was Too Big In Size."
		text := "The " + cfg.BackupName + " database backup was bigger than " + humanize.Bytes(uint64(high)) + " in size. It was " + humanize.Bytes(uint64(bytes)) + "."
		html := "<p>" + text + "</p>"

		// Send Email
		EmailSend(cfg.AlertEmail, subject, html, text)
	}

	// Return happy.
	return nil
}

/* End File */
