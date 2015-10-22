package sweet

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Extractor struct {
	client *http.Client
	doc    *goquery.Document
	err    error
}

func New() *Extractor {
	extractor := &Extractor{
		client: http.DefaultClient,
	}
	return extractor.EnableCookieJar()
}

func NewWithClient(client *http.Client) *Extractor {
	extractor := New()
	extractor.client = client
	return extractor
}

func FromURL(url string) *Extractor {
	return New().FromURL(url)
}

func FromReader(r io.Reader) *Extractor {
	return New().FromReader(r)
}

func FromString(html string) *Extractor {
	return New().FromReader(strings.NewReader(html))
}

func (extractor *Extractor) EnableCookieJar() *Extractor {
	jar, _ := cookiejar.New(nil)
	extractor.client.Jar = jar
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

func (extractor *Extractor) SelectForm(sel string) (*Form, error) {
	if extractor.err != nil {
		return nil, extractor.err
	}

	selectedForms := extractor.doc.Find(sel)

	if selectedForms.Size() == 0 {
		return nil, ErrNotFound{sel}
	}

	formAction, _ := selectedForms.Attr("action")
	formMethod, _ := selectedForms.Attr("method")

	textInputsSel := selectedForms.Find("input[type=text]")
	hiddenInputsSel := selectedForms.Find("input[type=hidden]")

	if textInputsSel.Size() == 0 && hiddenInputsSel.Size() == 0 {
		return nil, ErrEmptyForm{sel}
	}

	form := NewForm()
	form.Method = strings.ToUpper(formMethod)
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
