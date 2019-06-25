package types

//ShutdownSupernodeOperation represents shutdown_supernode operation data.
type ShutdownSupernodeOperation struct {
	Owner string `json:"owner"`
}

//Type function that defines the type of operation ShutdownSupernodeOperation.
func (op *ShutdownSupernodeOperation) Type() OpType {
	return TypeShutdownSupernode
}

//Data returns the operation data ShutdownSupernodeOperation.
func (op *ShutdownSupernodeOperation) Data() interface{} {
	return op
}
