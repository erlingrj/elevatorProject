package Timer

import (
  . "elevatorProject/Driver"
  "time"
)

const openDoorTime = 3
const reachingFloorTime = 3

func RunTimer(timeout chan TimerType, startTimer chan TimerType) {

  doorTimer := time.NewTimer(0)
  doorTimer.Stop()

  floorTimer := time.NewTimer(0)
  floorTimer.Stop()

  for {

    select {

    case timerType := <-startTimer:

      if timerType == TimeToOpenDoors {
        doorTimer.Reset(openDoorTime * time.Second)

      } else if timerType == TimeToReachFloor {
        floorTimer.Reset(reachingFloorTime * time.Second)

      } else if timerType == TimeFloorReached {
        floorTimer.Stop()
      }

    case <-floorTimer.C:
      timeout <- TimerType(0)

    case <-doorTimer.C:
      timeout <- TimerType(1)

    }

  }

}
