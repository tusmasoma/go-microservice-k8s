package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/tusmasoma/go-microservice-k8s/proto/web"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func Send(status int, msg proto.Message, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "failed to encode JSON: %v"}`, err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Error("failed to write response: %v\n", err)
	}
}

func OK(msg proto.Message, w http.ResponseWriter, r *http.Request) {
	Send(http.StatusOK, msg, w, r)
}

func Error(err error, w http.ResponseWriter, r *http.Request) {
	msg := &web.Error{}
	msg.Message = err.Error()
	// TODO: output logs such as sentry
	Send(http.StatusInternalServerError, msg, w, r)
}

func BadRequest(err error, w http.ResponseWriter, r *http.Request) {
	msg := &web.Error{}
	msg.Message = err.Error()
	Send(http.StatusBadRequest, msg, w, r)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

type CustomResponseWriter struct {
	responseWriter http.ResponseWriter
	statusCode     int
}

func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		responseWriter: w,
	}
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.responseWriter.Write(b)
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *CustomResponseWriter) StatusCode() int {
	return w.statusCode
}
