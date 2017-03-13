package InitializeElevator

import (
	. "elevatorProject/Driver"
	"elevatorProject/Utilities"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func InitializeElevator() ElevatorData {
	//sett heisen i en etasje
	//oppdater structen.
	//sett initialisert til true
	InitElevator()

	if GetFloorSensorSignal() == -1 {
		SetMotorDirection(DirnUp)
		for GetFloorSensorSignal() == -1 {
		}
		SetMotorDirection(DirnStop)
	}
	var initializedData ElevatorData

	initializedData.ID = Utilities.GetMacAddr()
	initializedData.Floor = GetFloorSensorSignal()
	initializedData.Direction = GetMotorDirection()
	initializedData.Status = 0
	initializedData.Initiated = true
	initializedData.ForceUpdate = false
	initializedData.ScheduledReboot = false
	SetFloorIndicator(initializedData.Floor)

	return initializedData

}

func InitializeElevatorList() [N_ELEVATORS]ElevatorData {
	var elevatorDataList [N_ELEVATORS]ElevatorData
	for i := 0; i < N_ELEVATORS; i++ {
		elevatorDataList[i].ID = ""
		for j := 0; j < N_FLOORS; j++ {
			for k := 0; k < N_BUTTONS; k++ {
				elevatorDataList[i].Orders[j][k] = 0
			}
		}
		elevatorDataList[i].Initiated = false
	}

	return elevatorDataList
}

func RunBackupProcess() {
	fmt.Println("Spawning Backup")

	// Backup listening to primary

	udpListen := DialBroadcastUDP(20100)
	buffer := make([]byte, 1024)

	for i := 0; i < 3; {
		//Need to set a read deadline for the UDP
		t := time.Now()
		udpListen.SetDeadline(t.Add(10 * time.Millisecond))
		_, udpshit, err := udpListen.ReadFrom(buffer[:])

		if udpshit != nil {
		}

		if err != nil {
			i++
			time.Sleep(200 * time.Millisecond)
		} else {
			i = 0
			time.Sleep(200 * time.Millisecond)

		}
	}
	//Lost conneciton to primary.

	//Close connection and take over as primary
	udpListen.Close()

}

func RunPrimaryProcess() {
	fmt.Println("This is now the primary!")

	// Setting up connection to backup
	udpBroadcast := DialBroadcastUDP(20100)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("localhost:%d", 20100))

	//Spawning a backup process
	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newBackup.Start()
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("I AM ALIVE")

	for {
		udpBroadcast.WriteTo(msg, addr)
		time.Sleep(200 * time.Millisecond)

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
