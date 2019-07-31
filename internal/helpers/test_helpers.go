package helpers

import (
	"log"

	"github.com/pquerna/ffjson/ffjson"
)

func GetErrorResponse(msg string, _type string) *string {
	errResponse := make(map[string]interface{})
	errField := make(map[string]interface{})

	errResponse["error"] = errField
	if _type != "" {
		errField["type"] = _type
	}
	if msg != "" {
		errField["message"] = msg
	} else {
		errField["message"] = nil
	}

	resp, err := ffjson.Marshal(errResponse)
	if err != nil {
		log.Fatalf("error while marshalling err response: %s\n", err.Error())
		return nil
	}

	var respPointer = string(resp)
	return &respPointer
}
