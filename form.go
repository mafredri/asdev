package main

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
)

func submitForm(ctx context.Context, c *cdp.Client, expression string) error {
	domContentEventFired, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContentEventFired.Close()

	_, err = c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
		Expression: expression,
	})
	if err != nil {
		return err
	}

	_, err = domContentEventFired.Recv()
	return err
}

func setFormInputFiles(ctx context.Context, c *cdp.Client, selector string, files ...string) error {
	doc, err := c.DOM.GetDocument(ctx, dom.NewGetDocumentArgs())
	if err != nil {
		return err
	}

	fileInput, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
		NodeID:   doc.Root.NodeID,
		Selector: selector,
	})
	if err != nil {
		return err
	}

	err = c.DOM.SetFileInputFiles(ctx, &dom.SetFileInputFilesArgs{
		NodeID: &fileInput.NodeID,
		Files:  files,
	})
	return err
}
