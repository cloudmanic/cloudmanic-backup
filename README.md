### Generate public-private keypair

https://www.everythingcli.org/secure-mysqldump-script-with-encryption-and-compression/

```openssl req -x509 -nodes -newkey rsa:2048 -keyout mysqldump-secure.priv.pem -out mysqldump-secure.pub.pem```



go run *.go --decrypt=/Users/spicer/Development/cloudmanic-backup/tmp/cloudmanic_com_1515207144.sql.tar.gz.enc


go run *.go --backup