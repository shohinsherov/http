package main

import (
	"github.com/shohinsherov/http/pkg/server"
	//"time"
//	"fmt"
//	"io/ioutil"
	//"strconv"
//	"strings"
//	"bytes"
	//"io"
	"log"
	"net"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}

}

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host,port))

	srv.Register("/payments/{id}/{ds}", func(req *server.Request) {
		id := req.PathParams["id"]
		log.Print(id)
	})

	/*srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web-site"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" + 
			"Content-Type: text/html\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			body,
		 ))
	if err != nil {
		log.Print(err)
	}
	})
	
	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" + 
			"Content-Type: text/html\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			body,
		 ))
	if err != nil {
		log.Print(err)
	}
	})
*/
	log.Print("server run in ",host +":"+ port)
	return srv.Start()
}


/*func execute(host string, port string) (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			// Идём обслуживать следующего
			continue
		}

		err = handle(conn)
		if err != nil {
			log.Print(err)
			// Идём обслуживать следующего
			continue
		}
	}

}

func handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	// .....

	//conn.Write([]byte("Hello!\r\n"))

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		return nil
	}

	if err != nil {
		log.Print(err)
		return err
	}
	//log.Printf("%s", buf[:n])

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {}

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {

	}

	if version != "HTTP/1.1" {

	}

	if path == "/" {
		body, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			log.Print(err)
			return fmt.Errorf("can't read index.html: %w", err)
		}

		marker := "{{year}}"
		year := time.Now().Year()
		body = bytes.ReplaceAll(body, []byte(marker), []byte(strconv.Itoa(year)))

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" + 
				string(body),
		))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}
*/