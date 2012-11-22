package aisandbox

// Each constructor starts with name and description.
// All of them accept any amount of []float64 coordinates as the last parameter.

// NOTE: Defend is different from the other commands.
// NewDefend accepts any amount of float64 slices that can have duration component as third value.
// If duration is left out, it defaults to 0, which server translates to minimum allowed value.
// 
// Any mix of slices that have length two or three is allowed, slices of any other length are skipped.
//
// tl;dr:
// NewDefend("bacon1", "Is delicious", []float64{1.2, 2.3, 3.4}, []float64{7.2, -2.6})
// Would command bacon1 to defend direction [1.2, 2.3] for 3.4 seconds
// and then change direction to [7.2, -2.6] for as short time as the server allows and then repeat.
//
// NewDefend("bacon1", "Stare at wall", []float64{21, -29.2})
// Works like the old Defend did, bot stares at single direction until told otherwise. Duration is ignored.
func NewDefend(name, description string, directions ...[]float64) *Defend {
	direction := make([]*FacingDirection, 0, len(directions))
	for _, v := range directions {
		switch len(v) {
		case 2:
			direction = append(direction, &FacingDirection{v, 0})
		case 3:
			direction = append(direction, &FacingDirection{v[:2], v[2]})
		}
	}

	return &Defend{
		Bot:              name,
		FacingDirections: direction,
		Description:      description,
	}
}

// If direction == nil, direction is ignored and bot will look forward while moving.
func NewAttack(name, description string, direction []float64, target ...[]float64) *Attack {
	command := &Attack{
		Bot:         name,
		Target:      target,
		Description: description,
	}
	if direction != nil || len(direction) != 2 {
		command.LookAt = direction
	}
	return command
}

func NewCharge(name, description string, target ...[]float64) *Charge {
	return &Charge{
		Bot:         name,
		Target:      target,
		Description: description,
	}
}

func NewMove(name, description string, target ...[]float64) *Move {
	return &Move{
		Bot:         name,
		Target:      target,
		Description: description,
	}
}
