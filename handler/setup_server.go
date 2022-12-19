package handler

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/EmissarySocial/emissary/render"
	"github.com/EmissarySocial/emissary/server"
	"github.com/benpate/derp"
	"github.com/benpate/form"
	"github.com/benpate/rosetta/maps"
	"github.com/benpate/table"
	"github.com/labstack/echo/v4"
)

// SetupPageGet generates simple template pages for the setup server, based on the Templates and ID provided.
func SetupPageGet(factory *server.Factory, templates *template.Template, templateID string) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		config := factory.Config()
		useWrapper := (ctx.Request().Header.Get("HX-Request") != "true")

		header := ctx.Response().Header()
		header.Set("Content-Type", "text/html")
		header.Set("Cache-Control", "no-cache")

		if useWrapper {
			if err := templates.ExecuteTemplate(ctx.Response().Writer, "_header.html", config); err != nil {
				derp.Report(render.WrapInlineError(ctx, derp.Wrap(err, "setup.getIndex", "Error rendering index page")))
			}
		}

		if err := templates.ExecuteTemplate(ctx.Response().Writer, templateID, config); err != nil {
			derp.Report(render.WrapInlineError(ctx, derp.Wrap(err, "setup.getIndex", "Error rendering index page")))
		}

		if useWrapper {
			if err := templates.ExecuteTemplate(ctx.Response().Writer, "_footer.html", config); err != nil {
				derp.Report(render.WrapInlineError(ctx, derp.Wrap(err, "setup.getIndex", "Error rendering index page")))
			}
		}

		return nil
	}
}

// SetupServerGet generates a form for the setup app.
func SetupServerGet(factory *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		// Data schema and UI schema
		config := factory.Config()
		schema := config.Schema()
		section := ctx.Param("section")
		uri := "/server/" + section

		// Find the correct form for this section (or fail)
		element, asTable, err := getSetupForm(section)

		if err != nil {
			return derp.Wrap(err, "setup.serverTable", "Invalid table name")
		}

		// Write Table-formatted forms.
		if asTable {
			widget := table.New(&schema, &element, &config, section, factory.Icons(), uri)
			return widget.Draw(ctx.Request().URL, ctx.Response().Writer)
		}

		// Fall through to single form
		widget := form.New(schema, element)
		result, err := widget.Editor(&config, nil)

		if err != nil {
			return derp.Wrap(err, "setup.serverTable", "Error creating form")
		}

		// Return the form
		return ctx.HTML(http.StatusOK, render.WrapForm(uri, result, "cancel-button:hide"))
	}
}

// SetupServerPost saves the form data to the config file.
func SetupServerPost(factory *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		data := maps.Map{}

		if err := ctx.Bind(&data); err != nil {
			return render.WrapInlineError(ctx, derp.Wrap(err, "setup.serverPost", "Error parsing form data"))
		}

		// Data schema and UI schema
		config := factory.Config()
		schema := config.Schema()
		section := ctx.Param("section")
		uri := "/server/" + section

		// Find the correct form for this section (or fail)
		element, asTable, err := getSetupForm(section)

		if err != nil {
			return render.WrapInlineError(ctx, derp.Wrap(err, "setup.serverTable", "Invalid table name"))
		}

		// Write Table-formatted forms.
		if asTable {
			widget := table.New(&schema, &element, &config, section, factory.Icons(), uri)

			// Apply the changes to the configuration
			if err := widget.Do(ctx.Request().URL, data); err != nil {
				return render.WrapInlineError(ctx, derp.Wrap(err, "setup.serverTable", "Error saving form data"))
			}

			// Try to save the configuration to the persistent storage
			if err := factory.UpdateConfig(config); err != nil {
				return render.WrapInlineError(ctx, derp.Wrap(err, "setup.postServer", "Internal error saving config.  Try again later."))
			}

			// Redraw the table
			return widget.DrawView(ctx.Response().Writer)
		}

		// Fall through to single form
		form := form.New(schema, element)

		// Apply the changes to the configuration
		if err := form.SetAll(&config, data, nil); err != nil {
			return render.WrapInlineError(ctx, derp.Wrap(err, "setup.serverPost", "Error saving form data", data))
		}

		// Try to save the configuration to the persistent storage
		if err := factory.UpdateConfig(config); err != nil {
			return render.WrapInlineError(ctx, derp.Wrap(err, "setup.postServer", "Internal error saving config.  Try again later."))
		}

		// Success!
		return render.WrapInlineSuccess(ctx, "Record Updated")
	}
}

// getSetupForm generates the different form layouts to use on the setup/server page.
func getSetupForm(name string) (form.Element, bool, error) {

	switch name {

	case "templates":
		return form.Element{
			Type: "layout-vertical",
			Children: []form.Element{
				{Type: "select", Label: "Adapter", Path: "adapter"},
				{Type: "text", Label: "Location", Path: "location", Options: maps.Map{"column-width": "100%"}},
			},
		}, true, nil

	case "emails":
		return form.Element{
			Type: "layout-vertical",
			Children: []form.Element{
				{Type: "select", Label: "Adapter", Path: "adapter"},
				{Type: "text", Label: "Location", Path: "location", Options: maps.Map{"column-width": "100%"}},
			},
		}, true, nil

	case "layouts":
		return form.Element{
			Type: "layout-vertical",
			Children: []form.Element{
				{Type: "select", Label: "Adapter", Path: "adapter"},
				{Type: "text", Label: "Location", Path: "location", Options: maps.Map{"column-width": "100&"}},
			},
		}, true, nil

	case "attachments":
		return form.Element{
			Type:        "layout-group",
			Description: "Readable/Writeable location where uploaded files (originals and thumbnails) are stored.",
			Children: []form.Element{
				{Type: "layout-vertical", Label: "Originals", Children: []form.Element{
					{Type: "select", Label: "Adapter", Path: "attachmentOriginals.adapter"},
					{Type: "text", Label: "Location", Path: "attachmentOriginals.location"},
				}},
				{Type: "layout-vertical", Label: "Thumbnails", Children: []form.Element{
					{Type: "select", Label: "Adapter", Path: "attachmentCache.adapter"},
					{Type: "text", Label: "Location", Path: "attachmentCache.location"},
				}},
			},
		}, false, nil

	case "certificates":
		return form.Element{
			Type:        "layout-vertical",
			Description: "Readable/Writeable location where SSL certificates are stored.",
			Children: []form.Element{
				{Type: "select", Label: "Adapter", Path: "certificates.adapter"},
				{Type: "text", Label: "Location", Path: "certificates.location"},
			},
		}, false, nil

	}

	return form.Element{}, false, derp.NewBadRequestError("handler.getSetupForm", "Invalid form name", name)
}
