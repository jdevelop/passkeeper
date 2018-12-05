package status

type StatusControl interface {
	SelfCheckInprogress() error
	SelfCheckComplete() error
	SelfCheckFailure(reason error) error
	ReadyToTransmit() error
	TransmissionComplete() error
	TransmissionFailure(reason error) error
}
