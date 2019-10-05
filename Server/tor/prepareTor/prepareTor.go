// +build windows
//
// Package tor provides a wrapper around the tor proxy command
package prepareTor

import (
	"archive/zip"
	//"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	//"strings"
	"syscall"

	"github.com/x3419/TorgoBot/Server/tor/utils"
)

const (
	// IsReadyMessage indicates that tor is ready for connections
	IsReadyMessage = "Bootstrapped 100%"
	TOR_ZIP_URL    = "https://dist.torproject.org/torbrowser/8.5.5/tor-win32-0.4.1.5.zip"
)

// Tor wraps the tor command line
type Tor struct {
	RootPath string    // RootPath is the path to extract the tor bundle zip
	Cmd      *exec.Cmd // Cmd is the tor proxy command
}

// New returns a new Tor instance
func New() *Tor {
	rootPath := "C:\\testaroo"
	//rootPath := "C:\\Users\\Analyst\\Desktop\\Tor Browser\\Browser"
	//fmt.Println("In New")
	return &Tor{RootPath: rootPath}
}

// Download downloads the tor zip distribution
func (t *Tor) Download(dst io.Writer) error {
	//fmt.Println("In Download")
	resp, err := http.Get(TOR_ZIP_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, &utils.DownloadProgressReader{
		Reader: resp.Body,
		Lenght: resp.ContentLength,
	})

	return err
}

// DownloadAndExtract download the tor bundle and extract to a given path
// It also returns the last known error
func (t *Tor) DownloadAndExtract() error {
	//fmt.Println("Not in DownloadAndExtract")
	if ok := utils.FileExists(t.GetExecutable()); ok {
		return nil
	}

	var buf bytes.Buffer
	err := t.Download(&buf)
	if err != nil {
		return err
	}

	zipWriter, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return err
	}

	for _, file := range zipWriter.File {
		err = t.extractFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tor) extractFile(file *zip.File) error {
	//fmt.Println("Not in extractFile")
	path := filepath.Join(t.RootPath, file.Name)
	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}
	return nil
}

// Start starts the tor proxy. Start blocks until an error occur or tor bootstraping
// is done
func (t *Tor) Start() error {
	cmd := t.GetExecutable()
	t.Cmd = exec.Command(cmd)
	//fmt.Println("Command executed")
	t.Cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	//stdout, err := t.Cmd.StdoutPipe()
	_, err := t.Cmd.StdoutPipe()

	if err != nil {
		return err
	}

	err = t.Cmd.Start()
	if err != nil {
		return err
	}

	/*scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), IsReadyMessage) {
			break
		}
	}*/

	return nil
}

// GetExecutable returns the path to tor.exe
func (t *Tor) GetExecutable() string {
	//return fmt.Sprintf("%s\\Tor\\tor.exe", t.RootPath)
	return fmt.Sprintf("%s\\Tor\\tor.exe", t.RootPath)
}

// Kill kill the tor process
func (t *Tor) Kill() error {
	err := t.Cmd.Process.Kill()
	if err != nil {
		return err
	}

	return nil
}

// Clean delete the tor folder
func (t *Tor) Clean() error {
	dir, _ := filepath.Split(t.GetExecutable())
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}
