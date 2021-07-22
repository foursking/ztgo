package metadata

const (
	// C/S
	Server = "server"
	Client = "client"

	// networks
	RemoteIP   = "remote_ip"
	LocalIP    = "local_ip"
	RemoteAddr = "remote_addr"
	LocalAddr  = "local_addr"

	// service types
	GRPC    = "gRPC"
	HTTP    = "HTTP"
	TCP     = "TCP"
	UDP     = "UDP"
	InnerPB = "innerPB"

	Success = "success"
	Failure = "failure"

	Timeout = "timeout"

	ErrCode = "err_code"
	Caller  = "caller"

	Unknown = "unknown"
)
