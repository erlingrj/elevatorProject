package EventController

import (
	"elevatorProject/ElevatorController"
	"elevatorProject/OrderController"
	"elevatorProject/Utilities"
	. "elevatorProject/driver"
	"fmt"
)

//Funksjoner som skal legges til ORDER MODULEN
//------------------------------------------------------------------------------------------------------------

func ArriveAtFloor(elevatorDataList [N_ELEVATORS]ElevatorData, floor int, startTimer chan TimerType) [N_ELEVATORS]ElevatorData {

	if floor == -1 {
		//Vi har forlatt en etasje, starter timeren
		startTimer <- TimerType(TimeToReachFloor)
		//elevatorData = RemoveCompletedOrders(elevatorData)
		//SetAllLights(elevatorData, AllExternalOrders())
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

			//elevatorData = RemoveCompletedOrders(elevatorData)
		}

		//ExternalOrders.UpdateElevatorData(elevatorDataList[0])
		//SetAllLights(elevatorData, AllExternalOrders())

	}
	return elevatorDataList
}

func ExternalButtonPressed(elevatorDataList [N_ELEVATORS]ElevatorData, order ElevatorOrder, newOrderTxCh chan ElevatorOrder, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) [N_ELEVATORS]ElevatorData {
	fmt.Println("enter")

	elevatorDataList = OrderController.PlaceExternalOrder(elevatorDataList, order, newOrderTxCh, updateElevatorTxCh)

	if order.Floor == GetFloorSensorSignal() && elevatorDataList[0].Direction == DirnStop {
		elevatorDataList = ArriveAtFloor(elevatorDataList, elevatorDataList[0].Floor, startTimer)

	} else if elevatorDataList[0].Direction == DirnStop {

		elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
		SetMotorDirection(elevatorDataList[0].Direction)
	}

	//	if elevatorDataList.Status == StatusIdle && elevatorDataList.Statys {}

	//	if elevatorDataList.Status == StatusIdle && elevatorDataList.Statys {}

	//OrderController.UpdateElevatorData(elevatorDataList[0])
	SetAllLights(elevatorDataList)

	return elevatorDataList

}

func LeaveFloor(elevatorDataList [N_ELEVATORS]ElevatorData, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {
	elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
	elevatorDataList = ElevatorController.RemoveCompletedOrders(elevatorDataList)
	SetDoorOpenLamp(0)
	SetAllLights(elevatorDataList)
	elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
	SetMotorDirection(elevatorDataList[0].Direction)
	updateElevatorTxCh <- elevatorDataList[0]
	return elevatorDataList
}

func InternalButtonPressed(elevatorDataList [N_ELEVATORS]ElevatorData, floor int, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) [N_ELEVATORS]ElevatorData {

	elevatorDataList = OrderController.PlaceInternalOrder(elevatorDataList, floor, updateElevatorTxCh)

	if elevatorDataList[0].Direction == DirnStop {

		if elevatorDataList[0].Floor == floor {
			elevatorDataList = ArriveAtFloor(elevatorDataList, floor, startTimer)
		} else {
			//fmt.Println("Internal")
			elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
			SetMotorDirection(elevatorDataList[0].Direction)
			//elevatorData = RemoveCompletedOrders(elevatorData)
			//SetAllLights(elevatorData, AllOrderController())
		}

	}
	//UpdateElevatorData(elevatorData)
	SetAllLights(elevatorDataList)

	return elevatorDataList

}

func ElevatorDataReceivedFromNetwork(elevatorDataRx ElevatorData, elevatorDataList [N_ELEVATORS]ElevatorData) [N_ELEVATORS]ElevatorData {
	if elevatorDataRx.ID == elevatorDataList[0].ID && elevatorDataRx.ForceUpdate == true {
		elevatorDataList[0].Orders = elevatorDataRx.Orders
		elevatorDataList[0].ForceUpdate = false
		elevatorDataList[0] = ElevatorController.GetNextDirection(elevatorDataList[0])
		SetMotorDirection(elevatorDataList[0].Direction)
		//UpdateElevatorData(elevatorData)

	} else {

		elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, elevatorDataRx.ID)] = elevatorDataRx
		//ExternalOrders.UpdateElevatorData(elevatorDataRx)
	}

	SetAllLights(elevatorDataList)

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
		SetAllLights(elevatorDataList)

	}

	return elevatorDataList

}

func ElevatorPeerUpdateFromNetwork(elevatorDataList [N_ELEVATORS]ElevatorData, onlineElevatorList PeerUpdate, updateElevatorTxCh chan ElevatorData, newOrderCh chan ElevatorOrder) [N_ELEVATORS]ElevatorData {

	//Setting all lost elevators to uninitiated, is probably unessecary.Unless two elevators fail at the same instant
	for i := 0; i < len(onlineElevatorList.Lost); i++ {
		elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.Lost[i])].Initiated = false
	}

	if onlineElevatorList.New == "" {
		OrderController.RedestributeExternalOrders(elevatorDataList, elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.Lost[len(onlineElevatorList.Lost)-1])], newOrderCh, updateElevatorTxCh)
	}

	if onlineElevatorList.New != "" {
		if FindElevatorIndex(elevatorDataList, onlineElevatorList.New) == -1 {
			elevatorDataList[FindElevatorIndex(elevatorDataList, "")].ID = onlineElevatorList.New
			elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].Initiated = true
		} else {

			elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].Initiated = true
			//If this elevator already has data stored we want to push those data back
			if OrderController.HasUnresolvedInternalOrders(elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.New)]) == true {
				elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.New)].ForceUpdate = true //To indicate that this data packet should be forced
				updateElevatorTxCh <- elevatorDataList[FindElevatorIndex(elevatorDataList, onlineElevatorList.New)]
			}
		}
	}

	if len(onlineElevatorList.Peers) == 1 && len(onlineElevatorList.Lost) > 1 {
	}
	//Lost connection to all other Elevators
	//Initialize reboot

	return elevatorDataList

}

func LeaveFloor(elevatorDataList [N_ELEVATORS]ElevatorData, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {
	elevatorDataList[0] = GetNextDirection(elevatorDataList[0])
	elevatorDataList = RemoveCompletedOrders(elevatorDataList)
	SetDoorOpenLamp(0)
	SetAllLights(elevatorDataList)
	elevatorDataList[0] = GetNextDirection(elevatorDataList[0])
	SetMotorDirection(elevatorDataList[0].Direction)
	updateElevatorTxCh <- elevatorDataList[0]
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
