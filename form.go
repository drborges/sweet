package sweet

import (
	"net/http"
	"net/url"
	"strings"
	"strconv"
)

type Form struct {
	client   *http.Client
	endpoint string
	Action   string
	Method   string
	Fields   url.Values
}

func NewForm() *Form {
	return &Form{
		Method: "POST",
		Fields: make(url.Values),
		client: http.DefaultClient,
	}
}

func (form *Form) SetEndpoint(endpoint string) *Form {
	form.endpoint = endpoint
	return form
}

func (form *Form) Submit() (*http.Response, error) {
	payload := form.Fields.Encode()
	target := form.endpoint + form.Action
	req, err := http.NewRequest(form.Method, target, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	return form.client.Do(req)
}
