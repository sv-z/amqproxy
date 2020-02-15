package spec091

const (
	frameMethod        = 1
	frameHeader        = 2
	frameBody          = 3
	frameHeartbeat     = 8
	frameMinSize       = 4096
	frameEnd           = 206 // "\xCE"
	replySuccess       = 200
	ContentTooLarge    = 311
	NoRoute            = 312
	NoConsumers        = 313
	ConnectionForced   = 320
	InvalidPath        = 402
	AccessRefused      = 403
	NotFound           = 404
	ResourceLocked     = 405
	PreconditionFailed = 406
	FrameError         = 501
	SyntaxError        = `502`
	CommandInvalid     = 503
	ChannelError       = 504
	UnexpectedFrame    = 505
	ResourceError      = 506
	NotAllowed         = 530
	NotImplemented     = 540
	InternalError      = 541
)
