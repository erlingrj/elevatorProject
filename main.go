package main


import (
	"fmt"
	"elevatorProject/EventController"
	"elevatorProject/InitializeElevator"
	"elevatorProject/Network"
	. "elevatorProject/Network/network/peers"
	. "elevatorProject/Driver"
	"elevatorProject/Timer"
	"elevatorProject/Utilities"
	"time"
)

func main() {

	InitializeElevator.RunBackupProcess()

	time.Sleep(1 * time.Second)
	go InitializeElevator.RunPrimaryProcess()

	elevatorData := InitializeElevator.InitializeElevator()
	ElevatorMasterList := InitializeElevator.InitializeElevatorList()
	ElevatorMasterList[0] = elevatorData

	updateElevatorRxCh := make(chan ElevatorData, 500)
	updateElevatorTxCh := make(chan ElevatorData, 500)

	startTimer := make(chan TimerType, 50)
	timeOut := make(chan TimerType, 50)

	newOrderTxCh := make(chan ElevatorOrder, 50)
	newOrderRxCh := make(chan ElevatorOrder, 50)

	peerUpdateCh := make(chan PeerUpdate, 50)
	peerTxEnableCh := make(chan bool)

	arriveAtFloorCh := make(chan int)
	externalButtonCh := make(chan ElevatorOrder, 50)
	internalButtonCh := make(chan int, 50)

	Utilities.PrintOrderList(ElevatorMasterList)

	go Network.RunNetwork(elevatorData, updateElevatorTxCh, updateElevatorRxCh, newOrderTxCh, newOrderRxCh, peerUpdateCh, peerTxEnableCh)

	go ReadFloorSensors(arriveAtFloorCh)
	go ReadButtonSensors(externalButtonCh, internalButtonCh)

	go Timer.RunTimer(timeOut, startTimer)

	for {
		select {

		case msg1 := <-arriveAtFloorCh:
			ElevatorMasterList = EventController.ArriveAtFloor(ElevatorMasterList, msg1, startTimer, updateElevatorTxCh)
			Utilities.PrintOrderList(ElevatorMasterList)

		case msg2 := <-externalButtonCh:
			ElevatorMasterList = EventController.ExternalButtonPressed(ElevatorMasterList, msg2, newOrderTxCh, updateElevatorTxCh, startTimer)
			fmt.Printf("External button pressed\n")
			fmt.Printf("-------------------------------------------------------------\n")

		case msg3 := <-internalButtonCh:
			ElevatorMasterList = EventController.InternalButtonPressed(ElevatorMasterList, msg3, updateElevatorTxCh, startTimer)
			fmt.Printf("Internal button pressed\n")
			fmt.Printf("-------------------------------------------------------------\n")

		case msg4 := <-updateElevatorRxCh:
			ElevatorMasterList = EventController.ElevatorDataReceivedFromNetwork(msg4, ElevatorMasterList, updateElevatorTxCh)

		case msg5 := <-newOrderRxCh:
			ElevatorMasterList = EventController.OrderReceivedFromNetwork(msg5, ElevatorMasterList, updateElevatorTxCh)

		case msg6 := <-peerUpdateCh:
			fmt.Println(msg6)
			ElevatorMasterList = EventController.ElevatorPeerUpdateFromNetwork(ElevatorMasterList, msg6, updateElevatorTxCh, newOrderTxCh)

		case timeout := <-timeOut:
			ElevatorMasterList = EventController.TimeOut(ElevatorMasterList, timeout, updateElevatorTxCh)

		}

	}

}
