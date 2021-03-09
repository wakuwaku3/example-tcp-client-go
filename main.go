package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := os.Getenv("TCP_SERVER_HOST")
	port := os.Getenv("TCP_SERVER_PORT")
	address := fmt.Sprintf("%s:%s", host, port)

	for {
		func() {
			connection := createConnection(address)
			defer connection.Close()
			ch1 := listenConnection(connection)
			ch2 := sendMessage(connection)
			select {
			case err := <-ch1:
				fmt.Fprintf(os.Stdin, "サーバーから切断されました\n%v\n", err)
			case err := <-ch2:
				fmt.Fprintf(os.Stdin, "サーバーから切断されました\n%v\n", err)
			}
		}()
	}
}

func createConnection(address string) net.Conn {
	fmt.Fprintf(os.Stdin, "%v に接続を開始します\n", address)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stdin, "接続に失敗しました\n%v\n", err)
		time.Sleep(time.Second * 5)
		return createConnection(address)
	}
	fmt.Fprintf(os.Stdin, "%v への接続に成功しました\n", address)
	return connection
}

func sendMessage(connection net.Conn) <-chan (error) {
	ch := make(chan (error))
	go func() {
		for {
			fmt.Print("> ")

			stdin := bufio.NewScanner(os.Stdin)
			if !stdin.Scan() {
				ch <- errors.New("scan failed")
				close(ch)
				return
			}
			_, err := connection.Write([]byte(stdin.Text() + "\n"))
			if err != nil {
				ch <- err
				close(ch)
				return
			}
		}
	}()

	return ch
}

func listenConnection(connection net.Conn) <-chan (error) {
	ch := make(chan (error))
	go func() {
		for {
			var response = make([]byte, 4*1024)

			if _, err := connection.Read(response); err != nil {
				ch <- err
				close(ch)
				return
			}

			fmt.Printf("Server> %s \n", response)
		}
	}()

	return ch
}
