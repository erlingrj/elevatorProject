package main

import (
	"fmt"
	"os/exec"
	"time"
	"net"
	"log"
	"syscall"
	"os"
)



func main() {

	RunBackupProcess()
	time.Sleep(1*time.Second)
	RunPrimaryProcess()

	select {}

}

func RunBackupProcess() {
	fmt.Println("Spawning Backup")

	// Backup listening to primary

	udpListen := DialBroadcastUDP(20100)
	buffer := make([]byte, 1024)

	for i:=0; i<3; {
		//Need to set a read deadline for the UDP
		t := time.Now()
		udpListen.SetDeadline(t.Add(10*time.Millisecond))
		_, udpshit, err := udpListen.ReadFrom(buffer[:])

		if udpshit != nil {}
		if err != nil {
			fmt.Println("Backup failed read nr: ", i)
			i++
			time.Sleep(200*time.Millisecond)
		} else {
			i = 0

		}
		}
		//Lost conneciton to primary.



		//Close connection and take over as primary
		udpListen.Close()



	}



	func RunPrimaryProcess() {
		fmt.Println("This is now the primary!")

		// Setting up connection to backup
		udpBroadcast:= DialBroadcastUDP(20100)
		addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", 20100))

		//Spawning a backup process
		newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run primary.go")
		err := newBackup.Start()
		if err != nil {log.Fatal(err)}

		msg := []byte("I AM ALIVE")

		for i := 0; i<20; i++{
			udpBroadcast.WriteTo(msg, addr)
			time.Sleep(50*time.Millisecond)


		}




	}

	func DialBroadcastUDP(port int) net.PacketConn {
		s, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
		syscall.Bind(s, &syscall.SockaddrInet4{Port: port})

		f := os.NewFile(uintptr(s), "")
		conn, _ := net.FilePacketConn(f)
		f.Close()

		return conn
	}
