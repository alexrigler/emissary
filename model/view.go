package model

import (
	"bytes"
	"html/template"

	"github.com/benpate/derp"
)

// View is an individual HTML template that can render a part of a stream
type View struct {
	Label       string             `json:"label"`       // Human-friendly label of this view
	Permissions []string           `json:"permissions"` // List of roles/users who can view this view
	HTML        string             `json:"html"`        // Raw HTML to render
	Compiled    *template.Template `json:"-"`           // Parsed HTML template to render (by merging with Stream dataset)
}

// Execute executes this template on the provided data.  It maintains a cache of the compiled template
func (v *View) Execute(data interface{}) (string, *derp.Error) {

	// If this view has already been compiled, then return the compiled version
	if v.Compiled == nil {

		result, err := template.New("").Parse(v.HTML)

		if err != nil {
			return "", derp.Wrap(err, "model.View.Template", "Unable to parse template HTML")
		}

		// Save the value into this view
		v.Compiled = result
	}

	var buffer bytes.Buffer

	if err := v.Compiled.Execute(&buffer, data); err != nil {
		return "", derp.Wrap(err, "Model.View.Template", "Error executing template", v.HTML, data)
	}

	// Return to caller.  TRUE means that the object has been changed.
	return buffer.String(), nil
}