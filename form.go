package sweet

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
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
	return nil, nil
}

func (form *Form) Set(field, value string) *Form {
	form.Fields.Set(field, value)
	return form
}

type Extractor struct {
	client *http.Client
	doc    *goquery.Document
	sel    string
	err    error
}

func FromReader(r io.Reader) *Extractor {
	doc, err := goquery.NewDocumentFromReader(r)
	return &Extractor{
		doc: doc,
		err: err,
		client: http.DefaultClient,
	}
}

func (extractor *Extractor) WithClient(client *http.Client) *Extractor {
	extractor.client = client
	return extractor
}

func (extractor *Extractor) Select(selector string) *Extractor {
	extractor.sel = selector
	return extractor
}

func (extractor *Extractor) ExtractForm() (*Form, error) {
	if extractor.err != nil {
		return nil, extractor.err
	}

	selectedForms := extractor.doc.Find(extractor.sel)

	if selectedForms.Size() == 0 {
		return nil, ErrNotFound{extractor.sel}
	}

	formAction, _ := selectedForms.Attr("action")
	formMethod, _ := selectedForms.Attr("method")

	selectedInputs := selectedForms.Find("input[type=text]")

	if selectedInputs.Size() == 0 {
		return nil, ErrEmptyForm{extractor.sel}
	}

	form := NewForm()
	form.Method = formMethod
	form.Action = formAction
	form.client = extractor.client

	selectedInputs.Each(func(_ int, input *goquery.Selection) {
		name, _ := input.Attr("name")
		value, _ := input.Attr("value")
		form.Fields.Set(name, value)
	})

	return form, nil
}
