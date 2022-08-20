package types

type ErrInvalidData struct {
	message string
}

func NewErrInvalidData(msg string) ErrInvalidData {
	return ErrInvalidData{message: msg}
}

func (e ErrInvalidData) Error() string {
	return e.message
}

type ErrDuplicatedOrder struct{}

func (ErrDuplicatedOrder) Error() string {
	return "duplicated order"
}

type ErrFlightImpossible struct{}

func (ErrFlightImpossible) Error() string {
	return "flight impossible for provided date and launchpad"
}
