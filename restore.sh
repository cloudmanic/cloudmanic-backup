#!/bin/bash

# $1 = file to download

# Small tool to help restore. This configs in this script need to be updated for your personal use.
# This bash script assume you are using mc - Minio Client to access your S3 Store

# Decript key
ENCRYPT_KEY=<SOME ENCRYPT KEY>

# DB Config Stuff
DB_PORT=3306 
DB_HOST=127.0.0.1
DB_USER=<SOME USER>
DB_PASS=<SOME PASSWORD>
DB_DATABASE=<SOME DATABASE>

MC_PATH=<SOME MC PATH> # ie. wasabi/db-backups/example.com
MC_DOWNLOAD="$MC_PATH/$1"

BASE_FILE=$(echo $1 | cut -f 1 -d '.')

# Set mysql password so it does not show up in PS
export MYSQL_PWD=$DB_PASS

if [ $# -ne 1 ]
then
	echo "Usage: $0 S3FileName"
	exit 1
fi

# First we load up the .env file locally for cloudmanic-backup
echo "ENCRYPT_KEY=$ENCRYPT_KEY" > .env

# Download file from S3 store via mc - Minio Client
mc cp $MC_DOWNLOAD .

# Decript the database.
cloudmanic-backup --decrypt=$1

# Delete encrypted file
rm -f $1 

# Unzip the backup.
gunzip "$BASE_FILE.sql.tar.gz"

# Untar the backup.
tar -xf "$BASE_FILE.sql.tar"

# Move the dump into the current directory
mv "tmp/$BASE_FILE.sql" .
rm -f "$BASE_FILE.sql.tar"
rm -rf tmp

# Truncate the restore db.
echo "Removing all the tables in $DB_DATABASE ...."
TABLES=$(mysql -u $DB_USER -P $DB_PORT -h $DB_HOST $DB_DATABASE -Nse "show tables")

for t in $TABLES
do
	mysql -u $DB_USER -P $DB_PORT -h $DB_HOST $DB_DATABASE -e "DROP TABLE $t"
done

# Restore the new database
echo "Restoring from backup..."
mysql -u $DB_USER -P $DB_PORT -h $DB_HOST $DB_DATABASE < $BASE_FILE.sql

# Delete the tar file.
rm -f $BASE_FILE.sql

# Remove .env file.
rm .env
