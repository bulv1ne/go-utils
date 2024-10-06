package monaderrors

// Option is a monad-like structure that wraps a value or represents nothing.
type Option[T any] struct {
	value T
	err   error
}

// Some creates an Option with a value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value}
}

// None creates an Option without a value and an error.
func None[T any](err error) Option[T] {
	return Option[T]{err: err}
}

// IsSome checks if the Option has no error.
func (o Option[T]) IsSome() bool {
	return o.err == nil
}

// IsError checks if the Option has an error.
func (o Option[T]) IsError() bool {
	return o.err != nil
}

// Map applies a function to the value inside the Option, if there is no error.
func (o Option[T]) Map(f func(T) (T, error)) Option[T] {
	if o.IsError() {
		return None[T](o.err)
	}
	newValue, err := f(o.value)
	return Option[T]{value: newValue, err: err}
}

// Map applies a function to the value inside the Option, if there is no error.
func Map[T1, T2 any](o Option[T1], f func(T1) (T2, error)) Option[T2] {
	if o.IsError() {
		return None[T2](o.err)
	}
	value, err := f(o.value)
	return Option[T2]{value: value, err: err}
}

// FlatMap (Bind) chains computations on the Option, if there is no error.
func (o Option[T]) FlatMap(f func(T) Option[T]) Option[T] {
	if o.IsError() {
		return None[T](o.err)
	}
	return f(o.value)
}

// FlatMap (Bind) chains computations on the Option, if there is no error.
func FlatMap[T1, T2 any](o Option[T1], f func(T1) Option[T2]) Option[T2] {
	if o.IsError() {
		return None[T2](o.err)
	}
	return f(o.value)
}

func (o Option[T]) Unwrap() (T, error) {
	return o.value, o.err
}
