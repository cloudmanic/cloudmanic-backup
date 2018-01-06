package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var (
	// TarCmd is the path to the `tar` executable
	TarCmd = "tar"

	// MysqlDumpCmd is the path to the `mysqldump` executable
	MysqlDumpCmd = "mysqldump"

	// Where we store backups for later storing to a better place.
	TmpPath = "/tmp/"
)

// MySQL is an `Exporter` interface that backs up a MySQL database via the `mysqldump` command
type MySQL struct {
	// DB Host (e.g. 127.0.0.1)
	Host string

	// DB Port (e.g. 3306)
	Port string

	// DB Name
	DB string

	// DB User
	User string

	// DB Password
	Password string

	// Extra mysqldump options
	// e.g []string{"--extended-insert"}
	Options []string

	// Key to encrypt the backup with.
	EncryptKey string

	// The size the backup most be bigger than (bytes)
	SizeCheckLow int64

	// The size the backup most be smaller than (bytes)
	SizeCheckHigh int64
}

//
// Export produces a `mysqldump` of the specified database, and creates a gzip compressed tarball archive.
//
func (x MySQL) Export() *ExportResult {

	result := &ExportResult{MIME: "application/x-tar"}

	dumpPath := TmpPath + fmt.Sprintf(`%v_%v.sql`, x.DB, time.Now().Unix())

	options := append(x.dumpOptions(), fmt.Sprintf(`-r%v`, dumpPath))
	out, err := exec.Command(MysqlDumpCmd, options...).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		return result
	}

	result.Path = dumpPath + ".tar.gz"
	_, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		return result
	}

	// Check to make sure the backup is within the range we expected
	AlertsBackupSize(dumpPath, x.SizeCheckLow, x.SizeCheckHigh)

	// Remove non-tar'ed version.
	os.Remove(dumpPath)

	// Encrypt the file.
	encPath, err := x.encryptFile(result.Path)

	if err != nil {
		result.Error = makeErr(err, "")
		return result
	}

	// Remove the non-encrypted path
	os.Remove(result.Path)

	// Set the encrypted path
	result.Path = encPath

	return result
}

//
// Encrypt file.
//
func (x MySQL) encryptFile(filePath string) (string, error) {

	content, err := readFromFile(filePath)

	if err != nil {
		return "", err
	}

	encrypted := encrypt(string(content), x.EncryptKey)
	writeToFile(encrypted, filePath+".enc")

	// Return new path.
	return filePath + ".enc", nil
}

//
// Mysql Dump Option
//
func (x MySQL) dumpOptions() []string {

	options := x.Options
	options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	options = append(options, fmt.Sprintf(`-P%v`, x.Port))
	options = append(options, fmt.Sprintf(`-u%v`, x.User))

	if x.Password != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Password))
	}

	options = append(options, x.DB)

	return options
}

/* End File */
