package main

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/pkg/errors"
)

func login(ctx context.Context, c *cdp.Client, username, password string) error {
	domLoadTimeout := 10 * time.Second
	err := navigate(ctx, c.Page, "http://developer.asustor.com/user/login", domLoadTimeout)
	if err != nil {
		return err
	}

	doc, err := c.DOM.GetDocument(ctx, dom.NewGetDocumentArgs())
	if err != nil {
		return err
	}

	input, err := c.DOM.QuerySelectorAll(ctx, dom.NewQuerySelectorAllArgs(doc.Root.NodeID, "#username, #password"))
	if err != nil {
		return err
	}

	if len(input.NodeIDs) != 2 {
		return errors.New("could not find #username and #password input")
	}

	// NOTE: When verbose mode is enabled, username and password will be
	// output to the terminal (via logging) in plain text.
	setValueArgs := []*dom.SetAttributeValueArgs{
		dom.NewSetAttributeValueArgs(input.NodeIDs[0], "value", username),
		dom.NewSetAttributeValueArgs(input.NodeIDs[1], "value", password),
	}

	for _, args := range setValueArgs {
		if err = c.DOM.SetAttributeValue(ctx, args); err != nil {
			return err
		}
	}

	return submitForm(ctx, c, `document.querySelector('input[value="Login"]').click();`)
}
