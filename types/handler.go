package types

// Handler defines the core of the application's
// state transition function
type Handler func(ctx Context, msg Msg) Result

// AnteHandler authenticates transactions before
// their internal messages are handled.
type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)

// Result is the result of a transaction
type Result struct {
	Code ABCICodeType
	Data []byte
	Log  string

	GasWanted int64
	GasUsed   int64
	Fee       sdk.Coins

	Tags Tags
}
