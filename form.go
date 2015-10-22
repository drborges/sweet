package sweet

import (
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	client *http.Client
	Action string
	Method string
	Fields url.Values
}

func NewForm() *Form {
	return &Form{
		Method: "POST",
		Fields: make(url.Values),
		client: http.DefaultClient,
	}
}

func (form *Form) Submit() (*http.Response, error) {
	req, err := http.NewRequest(form.Method, form.Action, strings.NewReader(form.Fields.Encode()))
	if err != nil {
		return nil, err
	}
	return form.client.Do(req)
}
