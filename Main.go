package main

import (
	"github.com/dlukes/Ground8/Tunnel"
	"github.com/tcnksm/go-input"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"runtime"
)

func main() {

	// Show the current version:
	log.Println(`Ground8 v0.0.1, based on SSHTunnel v1.3.0`)

	// Allow Go to use all CPUs:
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Read the configuration from the command-line args:
	readFlags()

	// Prompting UI
	ui := &input.UI{}

	// If username and/or password were not provided, prompt for them
	for true {
		if username == `` {
			query := `Username for ` + serverAddrString + `?`
			name, err := ui.Ask(query, &input.Options{
				Default:  os.Getenv("GROUND8_NAME"),
				Required: true,
				Loop:     true,
			})
			if err != nil {
				log.Fatalln(`Error reading username: ` + err.Error())
			} else {
				username = name
			}
		} else if password == `` {
			query := `Password for ` + username + ` at ` + serverAddrString + `?`
			passwd, err := ui.Ask(query, &input.Options{
				Default:  os.Getenv("GROUND8_PASSWD"),
				Required: true,
				Loop:     true,
				Mask:     true,
			})
			if err != nil {
				log.Fatalln(`Error reading password: ` + err.Error())
			} else {
				password = passwd
			}
		} else {
			break
		}
	}

	// Create the SSH configuration:
	Tunnel.SetPassword4Callback(password)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.PasswordCallback(Tunnel.PasswordCallback),
			ssh.KeyboardInteractive(Tunnel.KeyboardInteractiveChallenge),
		},
	}

	// Create the local end-point:
	localListener := Tunnel.CreateLocalEndPoint(localAddrString)

	// Accept client connections (will block forever):
	Tunnel.AcceptClients(localListener, config, serverAddrString, remoteAddrString)
}
