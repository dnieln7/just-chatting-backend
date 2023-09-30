package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseOK(writer http.ResponseWriter) {
	writer.WriteHeader(200)
	writer.Write([]byte{})
}

func ResponseNoContent(writer http.ResponseWriter) {
	writer.WriteHeader(204)
	writer.Write([]byte{})
}

func ResponseJson(writer http.ResponseWriter, code int, payload interface{}) {
	bytes, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal payload")

		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bytes)
}

func ResponseJsonError(writer http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Code is greather than 499: ", message)
	}

	ResponseJson(writer, code, errorResponse{
		Code:    code,
		Message: message,
	})
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
