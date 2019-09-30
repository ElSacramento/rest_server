package parse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, resp json.RawMessage) error {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		return err
	}
	return nil
}

// Parse request query params.
func RequestQueryJSON(r *http.Request) (json.RawMessage, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	rParams := make(map[string]interface{}, len(r.Form))
	for k, v := range r.Form {
		if len(v) == 1 {
			rParams[k] = v[0]
			continue
		}
		rParams[k] = v
	}

	rawJSON, err := json.Marshal(rParams)
	if err != nil {
		return nil, err
	}
	return rawJSON, nil
}

// Parse request body.
func RequestBodyJSON(r *http.Request) (json.RawMessage, error) {
	if r.Body == nil {
		return []byte(`{}`), nil
	}
	body, err := ioutil.ReadAll(r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return nil, err
	}
	return body, nil
}
