package gimple

import (
	"errors"
)

// Gimple represents simple container for services
// Go version of PHP container called Pimple
type Gimple struct {
	values map[string]Callable
}

// Callable represents function type for callbacks of container
type Callable func(context *Gimple) interface{}

// ExtendCallable represents function type for extend callbacks of container
type ExtendCallable func(origin interface{}, context *Gimple) interface{}

// New creates new instance of Gimple container
func New() *Gimple {
	return &Gimple{make(map[string]Callable)}
}

// SetService sets new service in container
func (gimple *Gimple) SetService(id string, value Callable) {
	gimple.values[id] = value
}

// GetService gets new service from container
func (gimple *Gimple) GetService(id string) (interface{}, error) {
	value, ok := gimple.values[id]
	if !ok {
		return nil, errors.New("There is no service with id " + id)
	}

	return value(gimple), nil
}

// ServiceExists checks if container exists
func (gimple *Gimple) ServiceExists(id string) bool {
	_, ok := gimple.values[id]
	return ok
}

// UnsetService unsets service from container
func (gimple *Gimple) UnsetService(id string) {
	delete(gimple.values, id)
}

// Share returns a closure that stores the result of the given service definition
// for uniqueness in the scope of this instance of Gimple.
func (gimple *Gimple) Share(callback Callable) Callable {
	var obj interface{}

	return func(c *Gimple) interface{} {
		if obj == nil {
			obj = callback(c)
		}
		return obj
	}
}

// Protect protects a callable from being interpreted as a service.
// This is useful when you want to store a callable as a parameter.
func (gimple *Gimple) Protect(callback Callable) Callable {
	return func(c *Gimple) interface{} {
		return callback
	}
}

// Raw gets a parameter or the closure defining an object.
func (gimple *Gimple) Raw(id string) (Callable, error) {
	value, ok := gimple.values[id]
	if !ok {
		return nil, errors.New("There is no service with id " + id)
	}

	return value, nil
}

// Extend extends an object definition.
// Useful when you want to extend an existing object definition,
// without necessarily loading that object.
func (gimple *Gimple) Extend(id string, callback ExtendCallable) (Callable, error) {
	factory, ok := gimple.values[id]
	if !ok {
		return nil, errors.New("There is no service with id " + id)
	}

	extender := func(c *Gimple) interface{} {
		return callback(factory(c), c)
	}

	return extender, nil
}

// Keys returns all defined value names.
func (gimple *Gimple) Keys() []string {
	keys := make([]string, 0, len(gimple.values))
	for k := range gimple.values {
		keys = append(keys, k)
	}
	return keys
}
