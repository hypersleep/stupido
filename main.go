package main

import(
	"log"
	"net"
	"flag"
	"os"
	"time"
	"bufio"
)

func main() {
	port := flag.String("port", "12202", "port to listen")
	path := flag.String("path", "/var/log/stupido", "path to write logs")
	name := flag.String("name", "default", "name of your log stream")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:" + *port)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	dir := *path + "/" + *name

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		log.Println("Path", dir, "not exist. Trying to create...")

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	filepath := dir + "/" + *name + ".log"

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	go func(f *os.File) {
		f.Sync()
		time.Sleep(2 * time.Second)
	}(f)

	log.Println("Stupido started at:", *port, "Writing at:", filepath)

	for {
		conn, err := l.Accept()
        if err != nil {
            log.Println("Error accepting: ", err.Error())
        }

		go func(conn net.Conn) {
			log.Println("Handling new connection...")

			defer func() {
				log.Println("Closing connection...")
				conn.Close()
			}()

			bufReader := bufio.NewReader(conn)

			for {
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))

				// Read tokens delimited by newline
				bytes, err := bufReader.ReadBytes('\n')
				if err != nil {
					log.Println(err)
					return
				}

				_, err = f.Write(bytes)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn)

	}
}
