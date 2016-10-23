package aero

import (
	"io"
	"strings"
	"time"

	"github.com/Joker/jade"
	"github.com/fatih/color"
	cache "github.com/patrickmn/go-cache"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasttemplate"
)

func init() {
	jade.PrettyOutput = false
}

// Template ...
type Template struct {
	Script      *otto.Script
	Code        string
	raw         string
	renderCache *cache.Cache
}

// NewTemplate ...
func NewTemplate(file string) *Template {
	template := new(Template)

	raw, _ := jade.ParseFile(file)
	raw = strings.TrimSpace(raw)
	raw = strings.Replace(raw, "{{ ", "{{", -1)
	raw = strings.Replace(raw, " }}", "}}", -1)
	raw = strings.Replace(raw, "\n", " ", -1)
	template.raw = raw

	t, _ := fasttemplate.NewTemplate(raw, "{{", "}}")

	code := "html = '" + t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		if tag == "end" {
			return w.Write([]byte("'; }\nhtml += '"))
		}

		if strings.HasPrefix(tag, "if ") {
			return w.Write([]byte("';\nif(" + tag[3:] + ") { html += '"))
		}

		if tag == "else" {
			return w.Write([]byte("';\nelse { html += '"))
		}

		return w.Write([]byte("';\nhtml += (" + tag + ");\nhtml += '"))
	}) + "';"

	code = strings.Replace(code, "html += '';", "", -1)
	template.Code = code

	compiler := otto.New()
	script, err := compiler.Compile(file, code)

	if err != nil {
		color.Red(err.Error())
	}

	template.Script = script

	template.renderCache = cache.New(5*time.Minute, 1*time.Minute)

	return template
}

// Render renders the template with the given parameters and returns the resulting string.
func (template *Template) Render(params map[string]interface{}) string {
	renderJobs <- renderJob{
		template,
		params,
	}
	return <-renderResults
}