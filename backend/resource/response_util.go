package resource

import (
	"encoding/json"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"net/http"
)

func SendPlainText(sender backend.CallResourceResponseSender, text string) error {
	return sendResourceResponse(
		sender,
		http.StatusOK,
		map[string][]string{
			"content-type": {"text/plain"},
		},
		[]byte(text),
	)
}

func SendJSON(sender backend.CallResourceResponseSender, obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return sendResourceResponse(
		sender,
		http.StatusOK,
		map[string][]string{
			"content-type": {"application/json"},
		},
		body,
	)
}

func sendResourceResponse(
	sender backend.CallResourceResponseSender,
	status int,
	headers map[string][]string,
	body []byte,
) error {
	return sender.Send(&backend.CallResourceResponse{
		Status:	status,
		Headers:	headers,
		Body:	body,
	})
}
