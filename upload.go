package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mafredri/asdev/apkg"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
)

func uploadAPK(ctx context.Context, verbose bool, devt *devtool.DevTools, errc chan<- error, apps []App, apk *apkg.File, opt options) (err error) {
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
		return errors.Wrap(err, "could not read apk config")
	}

	app := appSlice(apps).Find(conf.General.Package, conf.General.Architecture)
	if app.Valid() {
		return fmt.Errorf("could not upload %s (%s): app already exists", conf.General.Package, conf.General.Architecture)
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
		func() error { return c.Network.Enable(ctx, nil) },
		func() error { return c.Runtime.Enable(ctx) },

		func() error { return navigate(ctx, c.Page, "http://developer.asustor.com/app/upload", 10*time.Second) },
		func() error { return setFormInputFiles(ctx, c, "#appFile", absPath) },
		func() error { return setCheckboxOrRadio(ctx, c, `#is_beta`, false) },
		func() error { return submitForm(ctx, c, `document.getElementById('mainform').submit()`) },

		func() error { return setCategories(ctx, c, opt.categories) },
		func() error { return setInputValue(ctx, c, `#tags_en_US`, strings.Join(opt.tags, " ")) },
		func() error { return submitForm(ctx, c, `document.getElementById('mainform').submit()`) },
	} {
		if err = fn(); err != nil {
			return err
		}
	}

	/*
		Support Media Mode?

		<input type="hidden" name="media_mode" id="media_mode" value="0">
		<tr>
			<td align="center"><span class="star">*</span><strong>Media Mode</strong></td>
			<td class="addAppIcon">
				<input name="media_mode" id="no_media_mode" type="radio" value="0" checked="checked" chk="true" fieldname="Media Mode Description">  No media mode.<br>
				<input name="media_mode" id="is_media_mode" type="radio" value="1" chk="true" fieldname="Media Mode Description">  Some functions require media mode<br>
				<input name="media_mode" id="is_media_mode2" type="radio" value="2" chk="true" fieldname="Media Mode Description">  Be sure to media mode to use
			</td>
		</tr>
	*/

	// TODO: Implement icon/licence/submit for review.

	return nil
}

func setCategories(ctx context.Context, c *cdp.Client, cat Categories) error {
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

		if cat.Contains(text) {
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

	return nil
}

func setCheckboxOrRadio(ctx context.Context, c *cdp.Client, selector string, checked bool) error {
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		return err
	}

	sel, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
		NodeID:   doc.Root.NodeID,
		Selector: selector,
	})
	if err != nil {
		return err
	}

	if !checked {
		return c.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(sel.NodeID, "checked"))
	}
	return c.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(sel.NodeID, "checked", "checked"))
}

func setInputValue(ctx context.Context, c *cdp.Client, selector, value string) error {
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		return err
	}

	sel, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
		NodeID:   doc.Root.NodeID,
		Selector: selector,
	})
	if err != nil {
		return err
	}

	desc, err := c.DOM.DescribeNode(ctx, dom.NewDescribeNodeArgs().SetNodeID(sel.NodeID))
	if err != nil {
		return err
	}
	log.Println(selector, desc.Node.Name, desc.Node.NodeName, desc.Node.LocalName, desc.Node.NodeType, desc.Node.PseudoType, desc.Node)

	if desc.Node.NodeName == "TEXTAREA" {
		// TODO: Figure out how to set text content of textarea without eval.
		c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
			Expression: fmt.Sprintf(`document.querySelector(%q).textContent = %q`, selector, value),
		})
		return nil
	}

	return c.DOM.SetAttributeValue(ctx, &dom.SetAttributeValueArgs{
		NodeID: sel.NodeID,
		Name:   "value",
		Value:  value,
	})
}
