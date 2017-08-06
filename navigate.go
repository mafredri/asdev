package main

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
)

// navigate to the URL and wait for DOMContentEventFired. An error is
// returned if timeout happens before DOMContentEventFired.
func navigate(ctx context.Context, pc cdp.Page, url string, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Make sure nothing is loading.
	err = pc.StopLoading(ctx)
	if err != nil {
		return err
	}

	err = pc.Enable(ctx)
	if err != nil {
		return err
	}

	// Open client for DOMContentEventFired to block until DOM has fully loaded.
	domContentEventFired, err := pc.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContentEventFired.Close()

	_, err = pc.Navigate(ctx, page.NewNavigateArgs(url))
	if err != nil {
		return err
	}

	_, err = domContentEventFired.Recv()
	if err != nil {
		return err
	}

	return nil
}
