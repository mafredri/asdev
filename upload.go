package main

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/mafredri/asdev/apkg"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
)

func upload(ctx context.Context, verbose bool, devt *devtool.DevTools, errc chan<- error, apps []App, apk *apkg.File) (err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	defer func() {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		errc <- err
	}()

	conf, err := apk.Config()
	if err != nil {
		return errors.Wrap(err, "could read apk config")
	}

	app := appSlice(apps).Find(conf.General.Package, conf.General.Architecture)
	if app.Update != nil {
		return fmt.Errorf("could not update %s %s (%s): version %s is %s", app.Package, app.Version, app.Arch, app.Update.Version, app.Update.Status)
	}

	tab, err := devt.Create(ctx)
	if err != nil {
		return err
	}

	var opts []rpcc.DialOption
	if verbose {
		opts = append(opts, newLogCodec(conf.General.Package+":"+conf.General.Architecture))
	}
	conn, err := rpcc.DialContext(ctx, tab.WebSocketDebuggerURL, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	absPath, err := filepath.Abs(apk.Path())
	if err != nil {
		return err
	}

	for _, fn := range []func() error{
		func() error {
			return abortOnDetachOrCrash(ctx, c.Inspector, func(err error) {
				errc <- err
				cancel()
			})
		},
		func() error { return c.Inspector.Enable(ctx) },
		func() error { return c.DOM.Enable(ctx) },
		func() error { return c.Runtime.Enable(ctx) },
		func() error { return navigate(ctx, c.Page, app.UpdateURL(), 10*time.Second) },
		func() error { return setFormInputFiles(ctx, c, "#appFile", absPath) },
		func() error { return submitForm(ctx, c, `document.getElementById('mainform').submit()`) },
	} {
		if err = fn(); err != nil {
			return err
		}
	}

	// Fetch the document root for querying elements.
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		return err
	}

	// Fetch all product category options so that we can inherit them from
	// the previous version that was published.
	catOpts, err := c.DOM.QuerySelectorAll(ctx, &dom.QuerySelectorAllArgs{
		NodeID:   doc.Root.NodeID,
		Selector: "#category_id option",
	})
	if err != nil {
		return err
	}

	for _, o := range catOpts.NodeIDs {
		res, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{NodeID: &o})
		if err != nil {
			return err
		}

		// TODO: Cleanup this limitation of the Chrome Debuggin Protocol.
		// It only gives us OuterHTML.
		text := html.UnescapeString(optionTagRe.ReplaceAllString(res.OuterHTML, ""))

		if app.HasCategory(text) {
			err = c.DOM.SetAttributeValue(ctx, &dom.SetAttributeValueArgs{
				NodeID: o,
				Name:   "selected",
				Value:  "selected",
			})
			if err != nil {
				return err
			}
		}
	}

	descr := apk.Description()
	defer descr.Close()
	descb, err := ioutil.ReadAll(descr)
	if err != nil {
		return err
	}

	chlogr := apk.Changelog()
	defer chlogr.Close()
	chlogb, err := ioutil.ReadAll(chlogr)
	if err != nil {
		return err
	}

	for _, set := range []struct {
		id       string
		value    string
		textarea bool
	}{
		{id: "#name_en_US", value: conf.General.Name, textarea: false},
		// {id: "#tags_en_US", value: "", textarea: false},
		{id: "#description_en_US", value: string(descb), textarea: true},
		{id: "#changes_en_US", value: string(chlogb), textarea: true},
	} {
		if set.textarea {
			// TODO: Figure out how to set text content of textarea without eval.
			c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
				Expression: fmt.Sprintf(`document.querySelector(%q).textContent = %q`, set.id, set.value),
			})
			continue
		}

		sel, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
			NodeID:   doc.Root.NodeID,
			Selector: set.id,
		})
		if err != nil {
			return err
		}

		err = c.DOM.SetAttributeValue(ctx, &dom.SetAttributeValueArgs{
			NodeID: sel.NodeID,
			Name:   "value",
			Value:  set.value,
		})
		if err != nil {
			return err
		}
		continue
	}

	// Inherit the beta status for this app.
	if app.Beta {
		// Beta status is decided by #radio (Yes) and #radio2 (No).
		sel, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
			NodeID:   doc.Root.NodeID,
			Selector: "#radio",
		})
		if err != nil {
			return err
		}
		err = c.DOM.SetAttributeValue(ctx, &dom.SetAttributeValueArgs{
			NodeID: sel.NodeID,
			Name:   "checked",
			Value:  "checked",
		})
		if err != nil {
			return err
		}
	}

	err = submitForm(ctx, c, `document.getElementById('mainform').submit()`)
	if err != nil {
		return err
	}

	return nil
}
