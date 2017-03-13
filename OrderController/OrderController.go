package OrderController

import (
	//	. "elevatorProject/Network/network/peers"
	. "elevatorProject/Driver"
	"elevatorProject/Utilities"
)

func PlaceInternalOrder(elevatorDataList [N_ELEVATORS]ElevatorData, floor int, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {

	elevatorDataList[0].Orders[floor][ButtonType(2)] = 1
	updateElevatorTxCh <- elevatorDataList[0]

	return elevatorDataList
}

func PlaceExternalOrder(elevatorDataList [N_ELEVATORS]ElevatorData, order ElevatorOrder, newOrderTxCh chan ElevatorOrder, updateElevatorTxCh chan ElevatorData) [N_ELEVATORS]ElevatorData {
	order.ElevatorID = FindBestElevator(elevatorDataList, order)
	if order.ElevatorID == elevatorDataList[0].ID {
		//Oppdaterer egne ordreliste
		elevatorDataList[0].Orders[order.Floor][order.Direction] = 1

		//Sender oppdatert informasjon på nettverket
		updateElevatorTxCh <- elevatorDataList[0]

	} else {
		newOrderTxCh <- order
	}
	return elevatorDataList
}

func CalculateSingleElevatorCost(elevator ElevatorData, order ElevatorOrder) int {
	if (int(elevator.Direction) == -1 && int(order.Direction) == 1) || (int(elevator.Direction) == 1 && int(order.Direction) == 0) {
		switch elevator.Direction {
		case DirnUp:
			if order.Floor > elevator.Floor {
				return order.Floor - elevator.Floor
			} else {
				return (elevator.Floor-1)*2 + (elevator.Floor - order.Floor)
			}
		case DirnDown:
			if order.Floor < elevator.Floor {
				return elevator.Floor - order.Floor
			} else {
				return (elevator.Floor-1)*2 + (order.Floor - elevator.Floor)
			}
		}
	} else {
		switch elevator.Direction {
		case DirnUp:
			return 2*N_FLOORS - elevator.Floor - order.Floor
		case DirnDown:
			return (elevator.Floor - 1) + (order.Floor - 1)
		case DirnStop:
			distance := elevator.Floor - order.Floor
			if distance >= 0 {
				return distance
			}
			return -distance
		}
	}
	return -1
}

func FindBestElevator(elevatorDataList [N_ELEVATORS]ElevatorData, order ElevatorOrder) string {
	var minCost = 100000
	var ID string

	//var elevatorNumber = -1 //kanksje fint å bruke ID her?
	for i := 0; i < N_ELEVATORS; i++ {
		if elevatorDataList[i].Initiated {
			var thisCost = CalculateSingleElevatorCost(elevatorDataList[i], order)
			if thisCost < minCost {
				minCost = thisCost
				ID = elevatorDataList[i].ID
			}
		}
	}
	return ID //kan bare bruke ID-en til ordren
}

//tror det er best om bare en av heisene omfordeler ordre

func RedestributeExternalOrders(elevatorDataList [N_ELEVATORS]ElevatorData, lostElevator ElevatorData, newOrderCh chan ElevatorOrder, updateElevatorDataCh chan ElevatorData) {

	if Utilities.AmIMaster(elevatorDataList) == true {

		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < 2; j++ {
				if lostElevator.Orders[i][j] == 1 {
					elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, lostElevator.ID)].Orders[i][j] = 0
					newOrder := ElevatorOrder{i, j, ""}
					newOrder.ElevatorID = FindBestElevator(elevatorDataList, newOrder)

					newOrderCh <- newOrder
				}
			}
		}

		updateElevatorDataCh <- elevatorDataList[Utilities.FindElevatorIndex(elevatorDataList, lostElevator.ID)]
	}
}

func HasUnresolvedInternalOrders(elevatorData ElevatorData) bool {
	for i := 0; i < N_FLOORS; i++ {
		if elevatorData.Orders[i][2] == 1 {
			return true
		}
	}

	return false
}
