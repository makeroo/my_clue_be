package web

import "github.com/makeroo/my_clue_be/internal/platform/data"

// Request is an incoming request to be served.
type Request struct {
	// UserIO is the user who issued the request.
	UserIO  *UserIO
	ReqID   int
	Body    interface{}
	handler RequestHandler
}

// SendError returns an error message to the user.
func (req *Request) SendError(err error) {
	req.UserIO.send <- data.MessageFrame{
		Header: data.MessageHeader{
			Type:  data.MessageError,
			ReqID: req.ReqID,
		},
		Body: data.NotifyError{
			// FIXME: actually, Error() is safe only if error is a game.Error...
			Error: err.Error(),
		},
	}
}

// SendMessage returns a response to the user.
func (req *Request) SendMessage(messageType data.MessageType, body interface{}) {
	req.UserIO.send <- data.MessageFrame{
		Header: data.MessageHeader{
			Type:  messageType,
			ReqID: req.ReqID,
		},
		Body: body,
	}
}
