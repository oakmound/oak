package event

// An Entity is an element which can be bound to,
// in that it has a CID. All Entities need to implement
// is an Init function which should call NextID(e) and
// return that id:
// 		func (f *Foo) Init() event.CID {
//			f.CID = event.NextID(f)
//          return f.CID
// 		}
// In a multi-window setup each window may have its own
// callerMap, in which case event.NextID should be replaced
// with a NextID call on the appropriate callerMap.
type Entity interface {
	Init() CID
}

// Q: Why does every entity need its own implementation
// of Init()? Why can't it get that method definition
// from struct embedding?
//
// A: Because the CallerMap will store whatever struct is
// passed in to NextID. In a naive implementation:
// type A struct {
//     DefaultEntity
// }
//
// type DefaultEntity struct {
//     event.CID
// }
//
// func (de *DefaultEntity) Init() event.CID {
//     de.CID = event.NextID(de)
//     return de.CID
// }
//
// func main() {
//     ...
//     a := &A{}
//     cid := a.Init()
//     ent := event.GetEntity(cid)
//     _, ok := ent.(*A)
//     // ok is false, ent is type *DefaultEntity
// }
//
// So to effectively do this you would need something like:
// func DefaultEntity(parent interface{}) *DefaultEntity {}
// ... where the structure would store and pass down the parent.
// This introduces empty interfaces, would make initalization
// more difficult, and would use slightly more memory.
//
// Feel free to use this idea in your own implementations, but
// this package will not provide this structure at this time.
