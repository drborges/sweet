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

func TestExtractForm(t *testing.T) {
	form, err := sweet.FromReader(strings.NewReader(page)).Select("#first-form").ExtractForm()

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

func TestExtractFormErrNotFound(t *testing.T) {
	form, err := sweet.FromReader(strings.NewReader(page)).Select("#not-existent").ExtractForm()

	expectedErr := sweet.ErrNotFound{"#not-existent"}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected %v, got %v", expectedErr, err.Error())
	}

	if form != nil {
		t.Errorf("Expected %v, got %v", nil, form)
	}
}

func TestExtractFormErrEmptyForm(t *testing.T) {
	form, err := sweet.FromReader(strings.NewReader(page)).Select("#second-form").ExtractForm()

	expectedErr := sweet.ErrEmptyForm{"#second-form"}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected %v, got %v", expectedErr, err.Error())
	}

	if form != nil {
		t.Errorf("Expected %v, got %v", nil, form)
	}
}

func TestFormSubmit(t *testing.T) {
	form := sweet.NewForm()
	form.Action = "/session"
	form.Method = "POST"

	form.Submit()
}
