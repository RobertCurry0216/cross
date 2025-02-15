package common

type Viewable interface {
	View(state State) string
	Init(state State)
}
