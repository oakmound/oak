// Package event propagates events through entities with given caller IDs.
// It sets up a subscribe-publish model with the Bind and Trigger functions.
// In a slight change to the sub-pub model, event allows bindings to occur
// in an explicit order through assigning priority to individual bind calls.
package event
