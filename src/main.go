package main

import "github.com/bengtrj/simple-sync/sshclient"

import "os"

import "log"

import "fmt"

func main() {
	passd := os.Getenv("PASSWORD")

	testSSH("34.228.39.123:22", passd)
	testSSH("34.235.139.164:22", passd)
}

func testSSH(address, passd string) {
	client, err := sshclient.DialWithPasswd(address, "root", passd)
	if err != nil {
		log.Fatal(err)
	}

	script := `
	ls -la
	service --status-all
	`

	out, err := client.Cmd(script).SmartOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
