package move

import (
	"github.com/ethanlook/airbot/waypoint"
)

type Move interface {
	MoveOnMap(wp *waypoint.Waypoint) error
}

type MoveManger struct {
	// TODO

}

func NewMoveManager() Move {
	return &MoveManger{}
}

func (mm *MoveManger) MoveOnMap(wp *waypoint.Waypoint) error {
	req := mo
}
