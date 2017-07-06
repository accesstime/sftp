package main

import (
	"fmt"
	"os"

	"./config"
	log "github.com/Sirupsen/logrus"
	"github.com/koding/multiconfig"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {

	sftpConf := &config.Config{}
	m := multiconfig.NewWithPath("./config/config.toml")
	log.Infof("Loading configuration...")
	err := m.Load(sftpConf)
	if err != nil {
		log.Fatalf("Failed to load configuration. %v", err)
	}
	log.Infof("Configuration loaded as: %+v", *sftpConf)
	config := &ssh.ClientConfig{
		User:            sftpConf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpConf.Password),
		},
	}
	client, err := ssh.Dial("tcp", sftpConf.Server, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	fmt.Println("Successfully connected to ssh server.")

	// open an SFTP session over an existing ssh connection.

	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// Open the source file

	srcFile, err := sftp.Open(sftpConf.ServerPath + sftpConf.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Create the destination file

	dstFile, err := os.Create(sftpConf.LocalPath + sftpConf.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// Copy the file

	srcFile.WriteTo(dstFile)
}
