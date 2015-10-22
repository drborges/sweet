package sweet

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

type Extractor struct {
	client *http.Client
	doc    *goquery.Document
	sel    string
	err    error
}

func New() *Extractor {
	return &Extractor{
		client: http.DefaultClient,
	}
}

func NewWithClient(client *http.Client) *Extractor {
	extractor := New()
	extractor.client = client
	return extractor
}

func (extractor *Extractor) FromReader(r io.Reader) *Extractor {
	extractor.doc, extractor.err = goquery.NewDocumentFromReader(r)
	return extractor
}

func (extractor *Extractor) FromURL(url string) *Extractor {
	resp, err := extractor.client.Get(url)
	if err != nil {
		extractor.err = err
		return extractor
	}

	extractor.doc, extractor.err = goquery.NewDocumentFromReader(resp.Body)
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

	textInputsSel := selectedForms.Find("input[type=text]")
	hiddenInputsSel := selectedForms.Find("input[type=hidden]")

	if textInputsSel.Size() == 0 && hiddenInputsSel.Size() == 0 {
		return nil, ErrEmptyForm{extractor.sel}
	}

	form := NewForm()
	form.Method = formMethod
	form.Action = formAction
	form.client = extractor.client

	textInputsSel.Each(func(_ int, input *goquery.Selection) {
		name, _ := input.Attr("name")
		value, _ := input.Attr("value")
		form.Fields.Set(name, value)
	})

	hiddenInputsSel.Each(func(_ int, input *goquery.Selection) {
		name, _ := input.Attr("name")
		value, _ := input.Attr("value")
		form.Fields.Set(name, value)
	})

	return form, nil
}
