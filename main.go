package main

//legg filene i GOPATH/src (finn ved å skrive go env i terminal.)

import (
	"fmt"
	//"elevatorProject/ElevatorController"
	"elevatorProject/EventController"
	"elevatorProject/InitializeElevator"
	"elevatorProject/Network"
	. "elevatorProject/Network/network/peers"
	//"elevatorProject/OrderController"
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
			//fsmArriveAtFloor(msg)
			ElevatorMasterList = EventController.ArriveAtFloor(ElevatorMasterList, msg1, startTimer, updateElevatorTxCh)

		case msg2 := <-externalButtonCh:
			//elevatorData = fsmExternalButtonPressed(elevatorData, msg)
			//Utilities.PrintOrderList(ElevatorMasterList)
			ElevatorMasterList = EventController.ExternalButtonPressed(ElevatorMasterList, msg2, newOrderTxCh, updateElevatorTxCh, startTimer)

		case msg3 := <-internalButtonCh:
			//Utilities.PrintOrderList(ElevatorMasterList)
			ElevatorMasterList = EventController.InternalButtonPressed(ElevatorMasterList, msg3, updateElevatorTxCh, startTimer)
		case msg4 := <-updateElevatorRxCh:
			ElevatorMasterList = EventController.ElevatorDataReceivedFromNetwork(msg4, ElevatorMasterList, updateElevatorTxCh)

		case msg5 := <-newOrderRxCh:
			ElevatorMasterList = EventController.OrderReceivedFromNetwork(msg5, ElevatorMasterList, updateElevatorTxCh)

			//elevatorData = OrderReceivedOrder(elevatorData, msg)
		case msg6 := <-peerUpdateCh:
			fmt.Println(msg6)
			ElevatorMasterList = EventController.ElevatorPeerUpdateFromNetwork(ElevatorMasterList, msg6, updateElevatorTxCh, newOrderTxCh)

		case timeout := <-timeOut:
			ElevatorMasterList = EventController.TimeOut(ElevatorMasterList, timeout, updateElevatorTxCh)

		}

	}

}
