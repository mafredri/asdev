package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/mafredri/asdev/apkg"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
)

func updateAPK(ctx context.Context, verbose bool, devt *devtool.DevTools, errc chan<- error, apps []App, apk *apkg.File, opt options) (err error) {
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

	if len(opt.categories) == 0 {
		opt.categories = app.Categories
	}

	setTags := func() error { return nil }
	if len(opt.tags) > 0 {
		setTags = func() error {
			return setInputValue(ctx, c, "#tags_en_US", strings.Join(opt.tags, " "))
		}
	}

	descb, err := readAll(apk.Description())
	if err != nil {
		return err
	}

	chlogb, err := readAll(apk.Changelog())
	if err != nil {
		return err
	}

	betaStatus := app.Beta || opt.beta

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

		func() error { return setCategories(ctx, c, opt.categories) },
		func() error { return setInputValue(ctx, c, "#name_en_US", conf.General.Name) },
		func() error { return setInputValue(ctx, c, "#description_en_US", string(descb)) },
		func() error { return setInputValue(ctx, c, "#changes_en_US", string(chlogb)) },
		setTags,
		// Beta status is decided by #radio (Yes) and #radio2 (No).
		func() error { return setCheckboxOrRadio(ctx, c, "#radio", betaStatus) },
		func() error { return setCheckboxOrRadio(ctx, c, "#radio2", !betaStatus) },
		func() error { return submitForm(ctx, c, `document.getElementById('mainform').submit()`) },
	} {
		if err = fn(); err != nil {
			return err
		}
	}

	return nil
}

func readAll(r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	return ioutil.ReadAll(r)
}
