package driver

const N_ELEVATORS = 3

const N_FLOORS = 4

// Number of buttons (and corresponding lamps) on a per-floor basis
const N_BUTTONS = 3

type MotorDirection int8

const (
	DirnDown = -1 + iota
	DirnStop
	DirnUp
)

type TimerType int

const (
	TimeFloorReached = -1 + iota
	TimeToReachFloor
	TimeToOpenDoors
)

type ButtonType int

const (
	ButtonCallUp = iota
	ButtonCallDown
	ButtonInternal
)

type ElevatorStatus int

const (
	StatusIdle = iota
	StatusDoorOpen
	StatusMoving
)

type ElevatorOrder struct {
	Floor      int
	Direction  int //-1 for ned, 1 for opp
	ElevatorID string
}

type ElevatorData struct {
	Floor           int
	Direction       MotorDirection
	Orders          [N_FLOORS][N_BUTTONS]int
	Status          ElevatorStatus
	ID              string
	Initiated       bool
	ForceUpdate     bool
	ScheduledReboot bool
}
