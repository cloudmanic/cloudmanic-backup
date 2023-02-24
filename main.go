package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jasonlvhit/gocron"
	_ "github.com/jpfuentes2/go-env/autoload"
)

// CMD Action
type CmdAction struct {
	Action string
	Value  string
}

//
// Main()
//
func main() {
	// Check for CLI flags.
	mode := parseCmdFlags()

	// Figure out what to do based on command line flags (or lack there of)
	switch mode.Action {

	// Start in daemon mode (runs like cron on a timer.)
	case "daemon":
		runDaemon()

	// Run a complete backup once.
	case "run-once":
		runFullBackup()

	// Decrypt the file we pass in.
	case "decrypt":
		decryptFile(mode.Value, os.Getenv("ENCRYPT_KEY"))
	}
}

//
// Parse command line flags.
//
func parseCmdFlags() *CmdAction {
	// Is this a decrypt action?
	decryptCmd := flag.String("decrypt", "", "Path to file to decrypt.")

	// Just run a backup once?
	runBackupOnce := flag.Bool("backup", false, "Run a backup now.")

	// Parse flags
	flag.Parse()

	if len(*decryptCmd) > 0 {
		return &CmdAction{Action: "decrypt", Value: *decryptCmd}
	}

	if *runBackupOnce {
		return &CmdAction{Action: "run-once", Value: ""}
	}

	// If we get here we assume daemon mode.
	return &CmdAction{Action: "daemon", Value: ""}
}

//
// Run in Daemon mode
//
func runDaemon() {
	Log("Starting Cloudmanic Backup In Daemon Mode: Backing up every " + os.Getenv("HOURS_BETWEEN_BACKUPS") + " Hours.")

	// Get the duration
	dur, _ := strconv.ParseInt(os.Getenv("HOURS_BETWEEN_BACKUPS"), 10, 64)

	// Setup jobs we need to run
	gocron.Every(uint64(dur)).Hours().Do(runFullBackup)

	// function Start start all the pending jobs
	<-gocron.Start()
}

//
// Run one full backup.
//
func runFullBackup() {
	// Setup the MYSQL connection (the only db we support at the moment)
	m := MySQL{
		Host:          cfg.MysqlHost,
		Port:          cfg.MysqlPort,
		DB:            cfg.MysqlDb,
		User:          cfg.MysqlUser,
		Password:      cfg.MysqlPassword,
		EncryptKey:    cfg.EncryptKey,
		SizeCheckLow:  cfg.DbSizeCheckLow,
		SizeCheckHigh: cfg.DbSizeCheckHigh,
	}

	// Store this backup with S3 (the only storage we support at the moment)
	store := &S3{
		Region:       cfg.ObjectRegion,
		Bucket:       cfg.ObjectBucket,
		AccessKey:    cfg.ObjectAccessKeyId,
		ClientSecret: cfg.ObjectSecreteAccessKey,
		EndPoint:     cfg.ObjectEndpoint,
	}

	// Backup the database and then send it to our database store.
	result := m.Export()

	if result.Error != nil {
		Log("Backup failed (001): " + result.Error.err.Error())
		return
	}

	err := result.To(cfg.BackupDbStoreDir, store)

	if err != nil {
		Log("Backup failed (002): " + err.err.Error())
		return
	}

	// Send health check notice.
	if len(cfg.PintSuccessUrl) > 0 {

		resp, err := http.Get(cfg.PintSuccessUrl)

		if err != nil {
			Log("Could send health check - " + cfg.PintSuccessUrl)
		}

		defer resp.Body.Close()

	}

	// Successful backup.
	Log("Backup Success: " + result.Path)
}

//
// Decrypt a file we pass in from the command line.
//
func decryptFile(file string, key string) {
	// Read file.
	content, err := readFromFile(file)

	if err != nil {
		Log(err.Error())
		os.Exit(1)
	}

	// Decrypt file.
	decrypted := decrypt(string(content), key)

	// Write decrypted file.
	writeToFile(decrypted, file[:len(file)-4])
}

//
// Log.
//
func Log(logStr string) {
	fmt.Println(logStr)
}

/* End File */
