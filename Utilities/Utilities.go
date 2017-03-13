package Utilities

import (
	. "elevatorProject/Network/network/peers"
	. "elevatorProject/driver"
	"fmt"
	"net"
)

func AmIMaster(elevatorDataList [N_ELEVATORS]ElevatorData) bool {
	for i := 1; i < len(elevatorDataList); i++ {
		if elevatorDataList[i].Initiated == true && elevatorDataList[i].ID > elevatorDataList[0].ID {
			return false
		}

	}

	return true
}

func getMacAddr() string {

	var currentNetworkHardwareName string

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		currentNetworkHardwareName = interf.Name

	}

	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	macAddress := netInterface.HardwareAddr
	id := macAddress.String()

	return id
}

func PrintOrderList(elevatorDataList [N_ELEVATORS]ElevatorData) {

	fmt.Printf("		UP		DOWN		INTERNAL \n")
	fmt.Printf("-------------------------------------------------------------\n")
	for i := 0; i < N_FLOORS; i++ {
		fmt.Printf("Floor %d", i+1)
		for j := 0; j < N_BUTTONS; j++ {
			fmt.Printf("		%d", elevatorDataList[0].Orders[i][j])

		}
		fmt.Printf("\n")

	}
	fmt.Printf("-------------------------------------------------------------")
	fmt.Printf("\n")
	fmt.Printf("Direction: %d", MotorDirection(elevatorDataList[0].Direction))
	fmt.Printf("\n")
	fmt.Printf("-------------------------------------------------------------")
	fmt.Printf("\n")

}

func FindElevatorIndex(elevatorDataList [N_ELEVATORS]ElevatorData, elevatorID string) int {
	index := -1
	for i := 0; i < N_ELEVATORS; i++ {
		if elevatorDataList[i].ID == elevatorID {
			return i
		}
	}
	return index
}