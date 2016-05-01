package main

import (
	"flag"
)

func readFlags() {
	flag.StringVar(&serverAddrString, `server`, `jupyter.korpus.cz:22`, `The (remote) SSH server, e.g. 'my.host.com', 'my.host.com:22', '127.0.0.1:22', 'localhost:22'.`)
	flag.StringVar(&localAddrString, `local`, `localhost:9999`, `The local end-point of the tunnel, e.g. '127.0.0.1:50000', 'localhost:50000'.`)
	flag.StringVar(&username, `user`, ``, `The user's name for the SSH server.`)
	flag.StringVar(&password, `pwd`, ``, `The user's password for the SSH server.`)
	flag.StringVar(&cloud9Cmd, `cloud9`, `/opt/cloud9/cloud9.py`, `Full path to cloud9.py on remote server.`)
	flag.Parse()
}
