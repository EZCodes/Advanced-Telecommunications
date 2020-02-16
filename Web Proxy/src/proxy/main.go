package main 

import (
	"net"
	"net/http"
	"time"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	//"encoding/json"
	"os"
	//"syscall"
	"fmt"
	//"crypto/tls"
	//"crypto/x509"
)

type ReadWriter struct {
	io.Reader
	io.Writer
}

//var cache = map[http.Request]http.Response

// TODO Make a channel for displaying and sending messaged to/from main program, maybe make this function just reader
// This function makes a terminal for our web proxy
func makeTerminal() {
	rw := ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	// uncomment this part if using unix like console and uncomment syscall package
//	oldState, err := terminal.MakeRaw(int(syscall.Stdin))
//	if err != nil {
//        log.Fatal(err)
//	}
//	defer terminal.Restore(int(syscall.Stdin), oldState)
	
	term := terminal.NewTerminal(rw,"WebProxy >")
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
// TODO make a handler for requests and a map to store them.
func main() {
	
	//go makeTerminal()
	certFile := "src/proxy/cert.pem"
	keyFile := "src/proxy/key.pem"
//	go log.Fatal(http.ListenAndServe(":42069",http.HandlerFunc(httpRequestHandler)))
	log.Fatal(http.ListenAndServeTLS(":42070",certFile, keyFile, http.HandlerFunc(httpsRequestHandler)))
	//log.Printf("Listeners set up successfully")
	
}

//TODO finish this fucntion
// This function handles incoming requests and responds to them if needed
func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	
	log.Printf("received the request")
	url := req.URL
    url.Host = req.Host
    url.Scheme = "http"
    
	// clone and forward
	proxyReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		log.Print(err)
	}
	proxyReq.Header = req.Header
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error sending forwarding request %v",err)
	}
	log.Printf("forwarded the request and received the response")	
	
	// clone the response and send to client
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("Read the body")
	for head, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(head, value)
		}
	}
	w.Write(body)
	if(resp.StatusCode != 200) {
		w.WriteHeader(resp.StatusCode)
	}
	log.Printf("Finished redirection")		
}

func httpsRequestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("received https request")
	if req.Method != http.MethodConnect {
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        return
    }
	log.Printf("Dialing...")
	destConn, err := net.DialTimeout("tcp", req.Host, time.Second*30)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
        return
    }
    log.Printf("Conncetion Hijacked")
    clientConn, _, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    go transfer(destConn, clientConn)
    go transfer(clientConn, destConn)
    log.Printf("Finished tunelling")
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
    defer func() { _ = dest.Close() }()
    defer func() { _ = src.Close() }()
    _, _ = io.Copy(dest, src)
}
