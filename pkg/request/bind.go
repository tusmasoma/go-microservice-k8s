package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func Bind(out proto.Message, r *http.Request) error {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("unsupported Content-Type: %s", r.Header.Get("Content-Type"))
	}
	return json.NewDecoder(r.Body).Decode(out)
}
