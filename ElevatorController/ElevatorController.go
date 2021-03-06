package ElevatorController

import (
	. "elevatorProject/Driver"
)

func SetAllLights(elevatorDataList [N_ELEVATORS]ElevatorData) {
	for k := 0; k < N_ELEVATORS; k++ {
		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < N_BUTTONS-1; j++ {
				if k == 0 {
					SetButtonLamp(ButtonType(j), i, elevatorDataList[k].Orders[i][j])
				} else if elevatorDataList[k].Orders[i][j] == 1 && k > 0 {
					SetButtonLamp(ButtonType(j), i, elevatorDataList[k].Orders[i][j])
				}
				SetButtonLamp(ButtonType(2), i, elevatorDataList[0].Orders[i][2])
			}
		}
	}
}

func RemoveCompletedOrders(elevatorDataList [N_ELEVATORS]ElevatorData) [N_ELEVATORS]ElevatorData {

	switch elevatorDataList[0].Direction {

	case DirnUp:

		if NoOrdersAboveCurrentFloor(elevatorDataList[0]) {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallDown] = 0 //hvis de som skal opp ikke trykker videre, slettes denne, og det er litt uheldig
		} else if NoOrdersBelowCurrentFloor(elevatorDataList[0]) && elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallUp] == 0 {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallDown] = 0 //hvis de som skal opp ikke trykker videre, slettes denne, og det er litt uheldig
		}
		elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallUp] = 0
		elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonInternal] = 0

	case DirnDown:

		if NoOrdersBelowCurrentFloor(elevatorDataList[0]) {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallUp] = 0
		} else if NoOrdersAboveCurrentFloor(elevatorDataList[0]) && elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallDown] == 0 {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallUp] = 0
		}
		elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallDown] = 0
		elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonInternal] = 0

	case DirnStop:
		if NoOrdersBelowCurrentFloor(elevatorDataList[0]) {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallDown] = 0
		}
		if NoOrdersAboveCurrentFloor(elevatorDataList[0]) {
			elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonCallUp] = 0
		}
		elevatorDataList[0].Orders[elevatorDataList[0].Floor][ButtonInternal] = 0
	}
	return elevatorDataList
}

func GetNextDirection(elevatorData ElevatorData) ElevatorData {

	if NoOrdersAboveCurrentFloor(elevatorData) && NoOrdersAtCurrentFloor(elevatorData) && NoOrdersBelowCurrentFloor(elevatorData) {
		elevatorData.Direction = DirnStop
		goto end
	}

	switch elevatorData.Direction {

	case DirnUp:
		if NoOrdersAboveCurrentFloor(elevatorData) {
			elevatorData.Direction = DirnDown
		} else {
			elevatorData.Direction = DirnUp
		}

	case DirnDown:
		if NoOrdersBelowCurrentFloor(elevatorData) {
			elevatorData.Direction = DirnUp
		} else {
			elevatorData.Direction = DirnDown
		}

	case DirnStop:
		if NoOrdersAtCurrentFloor(elevatorData) {
			if NoOrdersBelowCurrentFloor(elevatorData) {
				elevatorData.Direction = DirnUp
			} else if NoOrdersAboveCurrentFloor(elevatorData) {
				elevatorData.Direction = DirnDown
			}
		} else if elevatorData.Status != StatusDoorOpen {

			if NoOrdersAboveCurrentFloor(elevatorData) && !NoOrdersBelowCurrentFloor(elevatorData) {
				elevatorData.Direction = DirnDown
			} else if !NoOrdersAboveCurrentFloor(elevatorData) && NoOrdersBelowCurrentFloor(elevatorData) {
				elevatorData.Direction = DirnUp
			}
		}
	}

end:
	return elevatorData
}

func CheckIfShouldStop(elevatorData ElevatorData) bool {
	switch elevatorData.Direction {

	case DirnUp:
		if elevatorData.Orders[elevatorData.Floor][ButtonCallUp] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
			return true
		} else if NoOrdersAboveCurrentFloor(elevatorData) {
			return true
		}

	case DirnDown:
		if elevatorData.Orders[elevatorData.Floor][ButtonCallDown] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
			return true
		} else if NoOrdersBelowCurrentFloor(elevatorData) {
			return true
		}

	case DirnStop:
		return true
	}
	return false
}

func NoOrdersAboveCurrentFloor(elevatorData ElevatorData) bool {
	if elevatorData.Floor == N_FLOORS-1 {
		return true
	}

	for i := elevatorData.Floor + 1; i < N_FLOORS; i++ {
		if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
			return false
		}
	}
	return true
}

func NoOrdersBelowCurrentFloor(elevatorData ElevatorData) bool {
	if elevatorData.Floor == 0 {
		return true
	}
	for i := 0; i < elevatorData.Floor; i++ {
		if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
			return false
		}
	}
	return true
}

func NoOrdersAtCurrentFloor(elevatorData ElevatorData) bool {
	for i := 0; i < N_BUTTONS; i++ {
		if elevatorData.Orders[elevatorData.Floor][i] != 0 {
			return false
		}
	}
	return true
}
