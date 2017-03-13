package EventController

import (
	. "elevatorProject/Driver"
	"elevatorProject/ElevatorController"
	. "elevatorProject/Network/network/peers"
	"elevatorProject/OrderController"
	"elevatorProject/Utilities"
	"fmt"
)

//Funksjoner som skal legges til ORDER MODULEN
//------------------------------------------------------------------------------------------------------------

func ArriveAtFloor(elevatorDataList [N_ELEVATORS]ElevatorData, floor int, startTimer chan TimerType, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {

	if floor == -1 {
		//Vi har forlatt en etasje, starter timeren
		startTimer <- TimerType(TimeToReachFloor)
		//elevatorData = RemoveCompletedOrders(elevatorData)
		//ElevatorController.SetAllLights(elevatorData, AllExternalOrders())
	} else {

		startTimer <- TimerType(TimeFloorReached)
		SetFloorIndicator(floor)
		elevatorDataList[0].Floor = floor

		if ElevatorController.CheckIfShouldStop(elevatorDataList[0]) == true {
			elevatorDataList[0].Status = StatusDoorOpen
			SetMotorDirection(DirnStop)
			SetDoorOpenLamp(1)
			//elevatorData = OpenDoors(elevatorData)
			startTimer <- TimerType(TimeToOpenDoors)
			updateElevatorTxCh <- elevatorDataList[0]

			//elevatorData = RemoveCompletedOrders(elevatorData)
		}

		//ExternalOrders.UpdateElevatorData(elevatorDataList[0])
		//ElevatorController.SetAllLights(elevatorData, AllExternalOrders())

	}
	return elevatorDataList
}

func ExternalButtonPressed(elevatorDataList [N_ELEVATORS]ElevatorData, order ElevatorOrder, newOrderTxCh chan ElevatorOrder, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) [N_ELEVATORS]ElevatorData {

	elevatorDataList = OrderController.PlaceExternalOrder(elevatorDataList, order, newOrderTxCh, updateElevatorTxCh)

	if order.Floor == GetFloorSensorSignal() && elevatorDataList[0].Direction == DirnStop {
		elevatorDataList = ArriveAtFloor(elevatorDataList, elevatorDataList[0].Floor, startTimer, updateElevatorTxCh)

	} else if elevatorDataList[0].Direction == DirnStop {

		elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
		SetMotorDirection(elevatorDataList[0].Direction)
	}

	//	if elevatorDataList.Status == StatusIdle && elevatorDataList.Statys {}

	//	if elevatorDataList.Status == StatusIdle && elevatorDataList.Statys {}

	//OrderController.UpdateElevatorData(elevatorDataList[0])
	ElevatorController.SetAllLights(elevatorDataList)

	return elevatorDataList

}

func LeaveFloor(elevatorDataList [N_ELEVATORS]ElevatorData, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {
	elevatorDataList[0].Status = StatusIdle
	elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
	elevatorDataList = ElevatorController.RemoveCompletedOrders(elevatorDataList)
	SetDoorOpenLamp(0)
	ElevatorController.SetAllLights(elevatorDataList)
	elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
	SetMotorDirection(elevatorDataList[0].Direction)
	updateElevatorTxCh <- elevatorDataList[0]
	ElevatorController.SetAllLights(elevatorDataList)
	return elevatorDataList
}

func InternalButtonPressed(elevatorDataList [N_ELEVATORS]ElevatorData, floor int, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) [N_ELEVATORS]ElevatorData {

	elevatorDataList = OrderController.PlaceInternalOrder(elevatorDataList, floor, updateElevatorTxCh)

	if elevatorDataList[0].Direction == DirnStop {

		if elevatorDataList[0].Floor == floor {
			elevatorDataList = ArriveAtFloor(elevatorDataList, floor, startTimer, updateElevatorTxCh)
		} else {
			//fmt.Println("Internal")
			elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
			SetMotorDirection(elevatorDataList[0].Direction)
			//elevatorData = RemoveCompletedOrders(elevatorData)
			//ElevatorController.SetAllLights(elevatorData, AllOrderController())
		}

	}
	//UpdateElevatorData(elevatorData)
	ElevatorController.SetAllLights(elevatorDataList)

	return elevatorDataList

}

func ElevatorDataReceivedFromNetwork(elevatorDataRx ElevatorData, elevatorDataList [N_ELEVATORS]ElevatorData, elevatorUpdateTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {

	if elevatorDataRx.ID == elevatorDataList[0].ID && elevatorDataRx.ForceUpdate == true {
		//Another elevator wants to overwrite our internal order due to network loss

		//Appending internal orders
		for i := 0; i < N_FLOORS; i++ {
			if elevatorDataRx.Orders[i][2] == 1 {
				elevatorDataList[0].Orders[i][2] = 1
			}
		}

		elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
		SetMotorDirection(elevatorDataList[0].Direction)
		fmt.Println("Internal orders from network")
		//Sending a message back to inform that we have successfully updated our orderqueue
		elevatorUpdateTxCh <- elevatorDataList[0]
		elevatorDataList[0].ForceUpdate = false

	} else if elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)].ForceUpdate == true {
		//We have sent internal orders to this elevator, check if they have been received
		fmt.Println("Internal orders to other elevator check")
		check := true
		for i := 0; i < N_FLOORS; i++ {
			if elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)].Orders[i][2] == 1 && elevatorDataRx.Orders[i][2] == 0 {
				check = false
			}
		}

		if check == false {
			//Resend the update

			elevatorUpdateTxCh <- elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)]
		} else {
			//We have sucessfully delegated the internal orders to the elevator. Update our own verison of elevatorDataRx
			elevatorDataRx.ForceUpdate = false
			elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)] = elevatorDataRx
		}
	} else if elevatorDataRx.ID != elevatorDataList[0].ID {
		fmt.Println("Other elevator")
		elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)] = elevatorDataRx
		//ExternalOrders.UpdateElevatorData(elevatorDataRx)
	}

	ElevatorController.SetAllLights(elevatorDataList)
	return elevatorDataList
}

func OrderReceivedFromNetwork(order ElevatorOrder, elevatorDataList [N_ELEVATORS]ElevatorData, elevatorUpdateTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {

	if order.ElevatorID == elevatorDataList[0].ID {
		//The order belongs to this elevator
		elevatorDataList[0].Orders[order.Floor][order.Direction] = 1

		//Broadcast the updated elevator struct
		elevatorUpdateTxCh <- elevatorDataList[0]

		if elevatorDataList[0].Direction == DirnStop {
			elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
			SetMotorDirection(elevatorDataList[0].Direction)
		}

		//ExternalOrders.UpdateElevatorData(elevatorDataList[0])
		ElevatorController.SetAllLights(elevatorDataList)

	}

	return elevatorDataList

}

func ElevatorPeerUpdateFromNetwork(elevatorDataList [N_ELEVATORS]ElevatorData, onlineElevatorList PeerUpdate, updateElevatorTxCh chan ElevatorData, newOrderCh chan ElevatorOrder) [N_ELEVATORS]ElevatorData {

	//Setting all lost elevators to uninitiated, is probably unessecary.Unless two elevators fail at the same instant
	for i := 0; i < len(onlineElevatorList.Lost); i++ {
		elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.Lost[i])].Initiated = false
		fmt.Println("Initialized = false")
	}

	if onlineElevatorList.New == "" {
		fmt.Println("Lost elevator")
		OrderController.RedestributeExternalOrders(elevatorDataList, elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.Lost[len(onlineElevatorList.Lost)-1])], newOrderCh, updateElevatorTxCh)
	}

	if onlineElevatorList.New != "" {
		if Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New) == -1 {
			fmt.Println("New elevator")
			elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, "")].ID = onlineElevatorList.New
			elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].Initiated = true
		} else {
			fmt.Println("New old elevator")
			elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].Initiated = true
			//If this elevator already has data stored we want to push those data back
			if OrderController.HasUnresolvedInternalOrders(elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New)]) == true {
				elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].ForceUpdate = true //To indicate that this data packet should be forced
				updateElevatorTxCh <- elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, onlineElevatorList.New)]
			}
		}
	}

	if elevatorDataList[0].Initiated == false {
		fmt.Println("Delete external orders")
		for i := 0; i < N_FLOORS; i++ {
			elevatorDataList[0].Orders[i][0] = 0
			elevatorDataList[0].Orders[i][1] = 0

		}
	}
	//Lost connection to all other Elevators
	//Initialize reboot

	return elevatorDataList

}

func TimeOut(elevatorDataList [N_ELEVATORS]ElevatorData, timeout TimerType, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {

	if timeout == TimeToOpenDoors {
		elevatorDataList = LeaveFloor(elevatorDataList, updateElevatorTxCh)
	} else if timeout == TimeToReachFloor {
		SetMotorDirection(DirnStop)
		panic("Cant reach floor, sensor/engine error!")
	}

	return elevatorDataList
}
