package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"port-forwording/internal/common"
	"port-forwording/internal/model"
	"sync"
)

func init() {
	Init()
}

func Init() {
	var pfInfos []model.PortForwarding
	result := common.DB.Find(&pfInfos)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	fmt.Println(pfInfos)

	for _, prInfo := range pfInfos {
		go Create(prInfo.ListenPort, prInfo.TargetIp, prInfo.TargetPort)
	}
}

func Create(listenPort uint, targetIP string, targetPort uint) {
	localAddr := fmt.Sprintf("0.0.0.0:%d", listenPort)
	serverAddr := fmt.Sprintf("%s:%d", targetIP, targetPort)

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error Close listen:", err)
			return
		}
	}(listener)
	log.Printf("Listening on %s, forwarding to %s\n", localAddr, serverAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(clientConn, serverAddr)
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
