package karma

import (
	"fmt"
)

// Context is a element of key-value linked list of message contexts.
type Context struct {
	Key      string
	Value    interface{}
	Previous *Context
}

// Context adds new key-value context pair to current context list and return
// current context list.
func (context *Context) Describe(
	key string,
	value interface{},
) *Context {
	previous := &Context{
		Key:      key,
		Value:    value,
		Previous: context.Previous,
	}

	context.Previous = previous

	return context
}

// Format produces context-rich hierarchical message, which will include all
// previously declared context key-value pairs.
func (context Context) Format(
	reason Reason,
	message string,
	args ...interface{},
) Karma {
	return Karma{
		Message: fmt.Sprintf(message, args...),
		Reason:  reason,
		Context: &context,
	}
}

// Reason adds current context to the specified message. If message is not
// hierarchical, it will be converted to such.
func (context Context) Reason(reason Reason) Karma {
	if previous, ok := reason.(Karma); ok {
		context.Walk(func(key string, value interface{}) {
			previous.Context = previous.Context.Describe(key, value)
		})

		return previous
	} else {
		return Karma{
			Reason:  reason,
			Context: &context,
		}
	}
}

// Walk iterates over all key-value context pairs and calls specified
// callback for each.
func (context *Context) Walk(callback func(string, interface{})) {
	if context == nil {
		return
	}

	callback(context.Key, context.Value)

	if context.Previous != nil {
		context.Previous.Walk(callback)
	}
}

// GetKeyValuePairs returns slice of key-value context pairs, which will
// be always even, each even index is key and each odd index is value.
func (context *Context) GetKeyValuePairs() []interface{} {
	pairs := []interface{}{}

	context.Walk(func(name string, value interface{}) {
		pairs = append(pairs, name, value)
	})

	return pairs
}