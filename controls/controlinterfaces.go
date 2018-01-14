package controls

// StatusControl defines the methods to be used to interact with the UI components on the board.
type StatusControl interface {
	SelfCheckInprogress() error
	SelfCheckComplete() error
	SelfCheckFailure(reason error) error
	ReadyToTransmit() error
	TransmissionComplete() error
	TransmissionFailure(reason error) error
}

// InputControl defines the button actions.
type InputControl interface {
	OnClickUp()
	OnClickDown()
	OnClickOk()
}

// Defines the display methods.
type DisplayControl interface {
	Refresh()
	ScrollUp(lines int)
	ScrollDown(lines int)
}
