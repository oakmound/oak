package physics

import "github.com/oakmound/oak/dlog"

const frozen = -64

// ForceVector is a vector that has some force and can operate on entites with mass
type ForceVector struct {
	Vector
	Force *float64
}

// NewForceVector returns a force vector
func NewForceVector(direction Vector, force float64) ForceVector {
	return ForceVector{Vector: direction, Force: &force}
}

// DefaultForceVector returns a force vector that converts the mass given
// into a force float
func DefaultForceVector(delta Vector, mass float64) ForceVector {
	return NewForceVector(delta, delta.Magnitude()*mass)
}

// GetForce is a self-returning call
func (f ForceVector) GetForce() ForceVector {
	return f
}

// GetForce on a non-force vector returns a zero-value for force
func (v Vector) GetForce() ForceVector {
	return ForceVector{v, new(float64)}
}

// A Mass can have forces applied against it
type Mass struct {
	mass float64
}

// SetMass of an object
func (m *Mass) SetMass(weight float64) {
	if weight > 0 {
		m.mass = weight
	} else {
		dlog.Warn("Tried to set somethings mass 0 or below", weight)
	}
}

//GetMass returns the mass of an object
func (m *Mass) GetMass() float64 {
	return m.mass
}

// Freeze changes a pushables mass such that it can no longer be pushed.
func (m *Mass) Freeze() {
	m.mass = frozen
}

// Pushable is implemented by anything that has mass and therefore can be pushed.
type Pushable interface {
	GetDelta() Vector
	GetMass() float64
}

// A Pushes can push Pushable things by having an associated ForceVector
type Pushes interface {
	GetForce() ForceVector
}

// Push applies the force from the pushing object its target
func Push(a Pushes, b Pushable) {
	dlog.Warn("Pushing", b.GetMass())
	if b.GetMass() <= 0 {
		if b.GetMass() != frozen {
			dlog.Warn("Pushed an object with", b.GetMass(), " mass")
		}
		return
	}
	//Copy a's force so that we dont change the original when we scale it later
	fdirection := a.GetForce().Copy()
	totalF := *a.GetForce().Force / b.GetMass()
	b.GetDelta().Add(fdirection.Normalize().Scale(totalF))
	dlog.Warn("Total Force was ", totalF, " fdirection ", fdirection.GetX(), fdirection.GetY())
}

// NOTE
// IMPORTANT
// NEVER CALL A RANDOM FUNCTION TO RESOLVE PHYSICS
// PHYSICS MUST BE DETERMINISTIC
// OTHERWISE BAD RATS HAPPENS
