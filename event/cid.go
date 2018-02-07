package event

// A CID is a caller ID that entities use to trigger and bind functionality
type CID int

// E is shorthand for GetEntity(int(cid))
// But we apparently forgot we added this shorthand,
// because this isn't used anywhere.
func (cid CID) E() interface{} {
	return GetEntity(int(cid))
}

// Parse returns the given cid, or the entity's cid
// if the given cid is 0. This way, multiple entities can be
// composed together by passing 0 down to lower tiered constructors, so that
// the topmost entity is stored once and bind functions will
// bind to the topmost entity.
func (cid CID) Parse(e Entity) CID {
	if cid == 0 {
		return e.Init()
	}
	return cid
}
