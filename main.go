package main

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func initializeDatabase(db *sql.DB) error {
	// 创建表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS port_forwarding (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			listen_port INTEGER NOT NULL,
			target_ip TEXT NOT NULL,
			target_port INTEGER NOT NULL
		);

		INSERT INTO port_forwarding (listen_port, target_ip, target_port) VALUES
			(6789, '10.200.59.71', 3306);
	`)
	return err
}

var db *gorm.DB

func main() {
	db, err := sql.Open("sqlite", "./.port_forwarding.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查数据库表是否存在，不存在则创建并初始化
	if err := initializeDatabase(db); err != nil {
		log.Fatal(err)
	}

	// 查询表中的数据
	rows, err := db.Query("SELECT listen_port, target_ip, target_port FROM port_forwarding")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var listenPort, targetPort int
		var targetIP string
		if err := rows.Scan(&listenPort, &targetIP, &targetPort); err != nil {
			log.Fatal(err)
		}

		go func(listenPort int, targetIP string, targetPort int) {
			localAddr := fmt.Sprintf("0.0.0.0:%d", listenPort)
			serverAddr := fmt.Sprintf("%s:%d", targetIP, targetPort)

			listener, err := net.Listen("tcp", localAddr)
			if err != nil {
				fmt.Println("Error listening:", err)
				return
			}
			defer listener.Close()
			fmt.Printf("Listening on %s, forwarding to %s\n", localAddr, serverAddr)

			for {
				clientConn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}

				go handleConnection(clientConn, serverAddr)
			}
		}(listenPort, targetIP, targetPort)
	}

	// 主程序持续运行
	for {
		time.Sleep(time.Minute)
		fmt.Println("Main function still running...")
	}
}

func handleConnection(clientConn net.Conn, serverAddr string) {
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer serverConn.Close()

	// 启动两个goroutine，分别将客户端的数据转发到服务器，将服务器的数据转发到客户端
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		copyData(clientConn, serverConn)
	}()

	go func() {
		defer wg.Done()
		copyData(serverConn, clientConn)
	}()

	wg.Wait()
}

func copyData(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Println("Error copying data:", err)
	}
}
