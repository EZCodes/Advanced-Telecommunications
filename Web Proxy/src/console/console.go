package main 

import (
	"bytes"
	"bufio"
	"os"
	"log"
	"strings"
	"net/http"
	"fmt"
)

// this file implements the console, that is used to constrol the proxy
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to proxy management console")
	fmt.Println("To block a website, write 'block: url'")
	fmt.Println("To unblock, write 'unblock: url'")
	for scanner.Scan() {
		command := scanner.Text()
		if strings.HasPrefix(command, "block") {
			fields := strings.Fields(command)
			b := []byte(fields[1])
			resp, err := http.Post("http://localhost:420", "text/plain" , bytes.NewBuffer(b))
			if err != nil {
				log.Printf("Failed blocking the address: %v", err)
			} else {
				defer resp.Body.Close()
				log.Printf("block message delivered")
			}
		} else if strings.HasPrefix(command, "unblock") {
			fields := strings.Fields(command)
			b := []byte(fields[1])
			resp, err := http.Post("http://localhost:421", "text/plain" ,bytes.NewBuffer(b))
			if err != nil {
				log.Printf("Failed unblocking the address: %v", err)
			} else {
				defer resp.Body.Close()
				log.Printf("Unblock message delivered")
			}
			
		} else {
			fmt.Println("Invalid input, please try again")
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

