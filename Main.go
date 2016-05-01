package main

import (
	"bytes"
	"github.com/dlukes/Ground8/Tunnel"
	"github.com/tcnksm/go-input"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) (string, error) {
	conn, err := ssh.Dial(`tcp`, hostname, config)
	if err != nil {
		return ``, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return ``, err
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)

	return stdoutBuf.String(), err
}

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
				Default:  os.Getenv(`GROUND8_NAME`),
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
				Default:  os.Getenv(`GROUND8_PASSWD`),
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

	// Launch Cloud9 on remote (or get pid/port of already running instance)
	var pid, port, killCmd, remoteAddrString string
	pidPort, err := executeCmd(cloud9Cmd, serverAddrString, config)
	if err != nil {
		log.Println(`Unable to launch or get Cloud9 server on ` + serverAddrString + `.`)
		os.Exit(1)
	} else {
		pidPort := strings.Split(pidPort, "\n")
		pid, port = pidPort[0], pidPort[1]
		killCmd = `kill ` + pid
		remoteAddrString = `localhost:` + port
		log.Println(`Cloud9 running as PID ` + pid + ` on ` + remoteAddrString + ` via ` + serverAddrString + `.`)
	}

	// Create the local end-point:
	localListener := Tunnel.CreateLocalEndPoint(localAddrString)

	log.Println(`Navigate to http://`+localAddrString+` in your browser to open Cloud9.`,
		`Press Ctrl-C to shut down Cloud9 on the remote and terminate the session.`)

	// Set up a handler for interrupt signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		log.Println(`Killing remote Cloud9 server...`)
		_, err := executeCmd(killCmd, serverAddrString, config)
		if err != nil {
			log.Println(`Error killing Cloud9 server: `, err.Error())
			os.Exit(1)
		} else {
			log.Println(`Success.`)
			os.Exit(0)
		}
	}()

	// Accept client connections (will block forever):
	Tunnel.AcceptClients(localListener, config, serverAddrString, remoteAddrString)
}
