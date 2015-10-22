package sweet_test

import (
	"github.com/drborges/sweet"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	page = `
		<html>
			<body>
				<form id="first-form" action="/form1" method="POST">
					<input type="text" name="user" value="diego" />
					<input type="text" name="auth_token" value="super_secret" />
					<input type="button" name="submit" value="OK" />
				</form>

				<form id="second-form" action="/form2" method="GET"></form>
			</body>
		</html>`
)

func TestExtractFormFromReader(t *testing.T) {
	form, err := sweet.New().FromReader(strings.NewReader(page)).Select("#first-form").ExtractForm()

	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	expectedFormAction := "/form1"
	expectedFormMethod := "POST"
	expectedFormBody := url.Values{}
	expectedFormBody.Add("auth_token", "super_secret")
	expectedFormBody.Add("user", "diego")

	if form.Action != expectedFormAction {
		t.Errorf("Expected %v, got %v", expectedFormAction, form.Action)
	}

	if form.Method != expectedFormMethod {
		t.Errorf("Expected %v, got %v", expectedFormMethod, form.Method)
	}

	if !reflect.DeepEqual(form.Fields, expectedFormBody) {
		t.Errorf("Expected %v, got %v", expectedFormBody, form.Fields)
	}
}

func TestExtractFormFromURL(t *testing.T) {
	form, err := sweet.New().FromURL("https://github.com/login").Select(".auth-form form").ExtractForm()

	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	expectedFormAction := "/session"
	expectedFormMethod := "POST"

	if form.Action != expectedFormAction {
		t.Errorf("Expected %v, got %v", expectedFormAction, form.Action)
	}

	if form.Method != expectedFormMethod {
		t.Errorf("Expected %v, got %v", expectedFormMethod, form.Method)
	}

	if form.Fields.Get("login") != "" {
		t.Errorf("Expected empty, got %v", form.Fields.Get("login"))
	}

	if form.Fields.Get("password") != "" {
		t.Errorf("Expected empty, got %v", form.Fields.Get("password"))
	}

	if form.Fields.Get("authenticity_token") == "" {
		t.Errorf("Expected not empty")
	}
}

func TestExtractFormErrNotFound(t *testing.T) {
	form, err := sweet.New().FromReader(strings.NewReader(page)).Select("#not-existent").ExtractForm()

	expectedErr := sweet.ErrNotFound{"#not-existent"}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected %v, got %v", expectedErr, err.Error())
	}

	if form != nil {
		t.Errorf("Expected %v, got %v", nil, form)
	}
}

func TestExtractFormErrEmptyForm(t *testing.T) {
	form, err := sweet.New().FromReader(strings.NewReader(page)).Select("#second-form").ExtractForm()

	expectedErr := sweet.ErrEmptyForm{"#second-form"}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected %v, got %v", expectedErr, err.Error())
	}

	if form != nil {
		t.Errorf("Expected %v, got %v", nil, form)
	}
}

func TestGithubLogin(t *testing.T) {
	form, err := sweet.New().EnableCookieJar().FromURL("https://github.com/login").Select(".auth-form form").ExtractForm()
	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	form.Fields.Set("login", "user")
	form.Fields.Set("password", "pass")

	resp, err := form.SetEndpoint("https://github.com").Submit()
	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	userSessionFound := false

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "user_session" {
			userSessionFound = true
		}
	}

	if !userSessionFound {
		t.Error("Expected response to have cookie user_session")
	}
}
