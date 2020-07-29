package common

// Maybe MONAD optional
type Maybe interface {
	Return(value interface{}) Maybe
	Bind(func(interface{}) Maybe) Maybe
}

// Just MONAD present
type Just struct {
	Value interface{}
}

// Nothing MONAD empty
type Nothing struct{}

// Return return if present
func (j Just) Return(value interface{}) Maybe {
	return Just{value}
}

// Bind process for present
func (j Just) Bind(f func(interface{}) Maybe) Maybe {
	return f(j.Value)
}

func (n Nothing) Return(value interface{}) Maybe {
	return Nothing{}
}

func (n Nothing) Bind(f func(interface{}) Maybe) Maybe {
	return Nothing{}
}
