package main 

import (
	//"net/http"
	//"time"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	//"encoding/json"
	"os"
	//"syscall"
	"fmt"
)

type ReadWriter struct {
	io.Reader
	io.Writer
}
// TODO Make a channel for displaying and sending messaged to/from main program
// This function makes a terminal for our web proxy
func makeTerminal() {
	rw := ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	// uncomment this part if using unix/cmd
//	oldState, err := terminal.MakeRaw(int(syscall.Stdin))
//	if err != nil {
//        log.Fatal(err)
//	}
//	defer terminal.Restore(int(syscall.Stdin), oldState)
	
	term := terminal.NewTerminal(rw,">")
	// test
	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return 
		}
		if err != nil {
			log.Print(err)
		}
		if line == "" {
			continue
		}
		fmt.Fprintln(term, line)
	}
}

func main() {
	
	go makeTerminal()
	
	
	for{}
//	tr := &http.Transport{
//		MaxIdleConns:       10,
//		IdleConnTimeout:    30 * time.Second,
//		DisableCompression: true,
//	}
//	client := &http.Client{ Transport: tr }
//	
//	s := &http.Server{
//		Addr:           ":8080",
//		Handler:        ,
//		ReadTimeout:    10 * time.Second,
//		WriteTimeout:   10 * time.Second,
//		MaxHeaderBytes: 1 << 20,
//	}
//	log.Fatal(s.ListenAndServe())
	

}

