package sweet

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Session map[string]string

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

func (form *Form) Submit() (*http.Response, Session, error) {
	payload := form.Fields.Encode()
	target := form.endpoint + form.Action

	req, err := http.NewRequest(form.Method, target, strings.NewReader(payload))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	resp, err := form.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	url, err := url.Parse(target)
	if err != nil {
		return nil, nil, err
	}

	session := make(Session)
	for _, cookie := range form.client.Jar.Cookies(url) {
		session[cookie.Name] = cookie.Value
	}

	return resp, session, err
}
