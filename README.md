# Sweet

Simple package for programmatically creating web login sessions based on cookies. 

# Example

The following example will log into a Twitter account by:

1. Visiting the login page generating the initial session cookies which are held within a cookie jar for future requests;
2. Collecting all defined input fields within the login form (including eventual CSRF tokens such as `authenticity_token`);
3. Posting the form data reusing cookies from the previous request;
4. Finally returning the session information upon a successful login. 

```Go
func TestTwitterLogin(t *testing.T) {
	form, err := sweet.FromURL("https://twitter.com/login").ExtractForm("form.signin")
	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	form.Fields.Set("session[username_or_email]", "user")
	form.Fields.Set("session[password]", "pass")

	_, session, err := form.Submit()
	if err != nil {
		t.Errorf("Expected nil, got", err)
	}

	if _, authTokenFound := session["auth_token"]; !authTokenFound {
		t.Error("Expected response to have cookie auth_token")
	}
}
```