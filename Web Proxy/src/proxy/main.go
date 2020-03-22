package main 

import (
	"net"
	"net/http"
	"time"
	"io"
	"io/ioutil"
	"log"
	"crypto/tls"
	"strings"
	"strconv"
)

const dateFormat = "Mon, 02 Jan 2006 15:04:05 MST"

type cachedResponse struct {
	response *http.Response // most of the response to use the headers
	responseBody []byte // processed body of the response, since on receival it's streamed (can only store it in this form)
}

// cache is shared
var httpCache = map[string]cachedResponse{}

// list of blocked URL's
var blockedURLs =  map[string]bool{}

func main() {
	
	s := &http.Server{
		Addr:           ":42070",
		Handler:        http.HandlerFunc(httpsRequestHandler),
		TLSNextProto: 	make(map[string]func(*http.Server, *tls.Conn, http.Handler)), //Disable HTTP/2
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":420",http.HandlerFunc(blockRequestHandler)))
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":421",http.HandlerFunc(unblockRequestHandler)))
	}()
	log.Fatal(http.ListenAndServe(":42069",http.HandlerFunc(httpRequestHandler)))	
}

// This function handles incoming requests and responds to them if needed
func httpRequestHandler(w http.ResponseWriter, req *http.Request) {	
	log.Printf("received the http request from : %s", req.URL.String())
	
	//check if we blocked this address
	for blockedURL, isBlocked := range blockedURLs {
		if strings.Contains(req.URL.String(), blockedURL) && isBlocked {
			w.WriteHeader(http.StatusLocked)
			log.Printf("Request blocked according to administrator policy.")
			return
		}	
	}
	client := &http.Client{}	
	url := req.URL
    url.Host = req.Host
    url.Scheme = "http"
    
    // take response from cache the response and send to client, if it's not expired and us there
    cachedResp, exists := httpCache[req.URL.String()]
    if exists {
    	var unformattedDate string
    	var maxAge int
    	var err error
    	for head, values := range cachedResp.response.Header {
			if head == "Cache-Control" {
				for _, value := range values {
					if strings.Contains(value, "max-age") {
						newVal := strings.Split(value, ",") //safety guard for wrongly parsed/constructed headers
						value = newVal[0]
						maxAge, err = strconv.Atoi(value[8:])
						if err != nil {
							log.Fatalf("Failed converting max-age value to integer: %v", err)
						}
					}
				}
			} else if head == "Date" {
				unformattedDate = values[0]
			}
		}
    	formattedDate, err := time.Parse(dateFormat, unformattedDate)
    	if err != nil {
    		log.Fatalf("Failed formatting Date from header into variable: %v", err)
    	}
    	expiryTime := formattedDate.Add(time.Duration(maxAge)*time.Second)
    	if expiryTime.After(time.Now()){
			cachedBody := cachedResp.responseBody
		    if err != nil {
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		        return
		    }
		    log.Printf("Read the cached body")
			for head, values := range cachedResp.response.Header {
				for _, value := range values {
				w.Header().Add(head, value)
				}
			}
			w.Write(cachedBody)
			log.Printf("Supplied cached response")
			return
    	} else {
    		log.Printf("Response timed out, fetching new one")
    	}
    }
	   
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
	
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	
	//conditionally cache the response
	canCache := true
	for head, values := range resp.Header {
		if head == "Cache-Control" {
			for _, value := range values {
				if value == "no-cache" || value == "no-store" || value == "private" {
					canCache = false
					log.Printf("Response not cache-able")
				}
			}
		}
	}
	if canCache {
		toCache := cachedResponse{
			response : resp,
			responseBody : body,
		}
		httpCache[req.URL.String()] = toCache
		log.Printf("Response cached")
	}
	
	// clone the response and send to client
    defer resp.Body.Close()
    log.Printf("Read the body")
	for head, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(head, value)
		}
	}
	w.Write(body)
	log.Printf("Finished cloning")		
}

func httpsRequestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("received https request from : %s", req.URL.String())
	//check if we blocked this address
	for blockedURL, isBlocked := range blockedURLs {
		if strings.Contains(req.URL.String(), blockedURL) && isBlocked {
			w.WriteHeader(http.StatusLocked)
			log.Printf("Request blocked according to administrator policy.")
			return
		}	
	}
	if req.Method != http.MethodConnect {
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        return
    }
	log.Printf("Dialing...")
	serverConn, err := net.DialTimeout("tcp", req.Host, time.Second*10)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    // Hijacking connection leaves it to us to manage the connection manually
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Could not Hijack the connection", http.StatusInternalServerError)
        return
    }
    log.Printf("Hijacking finished")
    clientConn, _, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    // create a https transfer tunnel between server and the client through our proxy, 
    // works both ways independently
    go tunnel(serverConn, clientConn)
    go tunnel(clientConn, serverConn)
    log.Printf("Finished tunelling")
}

func tunnel(dest io.WriteCloser, src io.ReadCloser) {
    defer func() { _ = dest.Close() }()
    defer func() { _ = src.Close() }()
    _, _ = io.Copy(dest, src)
}

// This method handles the block command requests from console
func blockRequestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received block request")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Problem reading block request body: %v", err)
	}
	url := string(body)
	blockedURLs[url] = true
}

// This method handles the unblock command requests from console
func unblockRequestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received unblock request")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Problem reading unblock request body: %v", err)
	}
	url := string(body)
	blockedURLs[url] = false
}
