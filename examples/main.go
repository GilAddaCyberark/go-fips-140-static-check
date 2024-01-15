package main

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"log"
)

func main() {
	sum := 3 + 2
	fmt.Printf("Sum: %d\n", sum)

	//	calcluate md5 for string
	str := "your string here"
	h := md5.New()
	io.WriteString(h, str)
	fmt.Printf("MD5 hash: %x\n", h.Sum(nil))

	// build a tls listener based on the default config
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

}
