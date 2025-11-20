package fastapi

import (
	"fmt"
	"io"
)

func (app *App) newClientCodec(rw io.ReadWriter) (link.COdec, error) {
	return app.newCodec(rw, app.newResponse), nil
}

func (app *App) newResponse(serviceID, messageID byte) (Message, error) {
	if service := app.services[serviceID]; service != nil {
		if msg := service.(Service).NewResponse(messageID); msg != nil {
			return msg, nil
		}
		return nil, DecodeError{fmt.Sprintf("Unsupported Message Type: [%d, %d]", serviceID, messageID)}
	}
	return nil, DecodeError{fmt.Sprintf("Unsupported Service: [%d, %d]", serviceID, messageID)}
}
