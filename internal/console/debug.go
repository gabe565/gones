package console

type Debug uint8

const (
	DebugDisabled = iota
	DebugWait
	DebugStepFrame
	DebugRunRender
)
