package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
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

type stringVar []string

func (s *stringVar) String() string {
	return fmt.Sprint(*s)
}

func (s *stringVar) Set(value string) error {
	for _, ss := range strings.Split(value, ",") {
		*s = append(*s, ss)
	}
	return nil
}

func main() {
	var (
		username   = flag.String("username", "", "Username (login)")
		password   = flag.String("password", "", "Password (login)")
		timeout    = flag.Duration("timeout", 5*time.Minute, "Timeout for package submission")
		verbose    = flag.Bool("v", false, "Verbose")
		browser    = flag.String("browser", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "Path to Chrome or Chromium binary")
		noHeadless = flag.Bool("no-headless", false, "Disable (Chrome) headless mode")
	)

	var apkVars stringVar
	flag.Var(&apkVars, "apk", "Package apk to submit")

	flag.Parse()

	if *username == "" {
		*username = os.Getenv("ASDEV_USERNAME")
	}
	if *password == "" {
		*password = os.Getenv("ASDEV_PASSWORD")
	}

	if *username == "" || *password == "" {
		fmt.Println("error: username or password is missing, use cli flag or set in environment")
		os.Exit(1)
	}

	var apks []*apkg.File
	for _, av := range apkVars {
		apk, err := apkg.Open(av)
		if err != nil {
			fmt.Printf("error: could open apk %q: %v\n", av, err)
			os.Exit(1)
		}
		defer apk.Close()
		apks = append(apks, apk)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	if err := run(ctx, *verbose, !*noHeadless, *browser, *username, *password, apks); err != nil {
		log.Fatal(err)
	}
}

var (
	optionTagRe = regexp.MustCompile("</?option( [^>]+)?>")
)

func run(ctx context.Context, verbose, headless bool, chromeBin string, username, password string, apks []*apkg.File) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tmpdir, err := ioutil.TempDir("", "asdev-chrome-userdata")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	chrome, err := startChrome(ctx, chromeBin, tmpdir, headless)
	if err != nil {
		return err
	}
	defer chrome.Close()

	devt := devtool.New(fmt.Sprintf("http://localhost:%d", chrome.port))
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		return err
	}

	var opts []rpcc.DialOption
	if verbose {
		opts = append(opts, newLogCodec("login"))
	}
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	err = login(ctx, c, username, password)
	if err != nil {
		return err
	}

	apps, err := getApps(ctx, c)
	if err != nil {
		return err
	}

	errc := make(chan chan error, len(apks))
	for _, apk := range apks {
		errc2 := make(chan error, 1)
		go upload(ctx, verbose, devt, errc2, apps, apk)
		errc <- errc2
	}
	close(errc)

	for e := range errc {
		select {
		case err = <-e:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

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

	tab, err := devt.Create(ctx)
	if err != nil {
		return err
	}

	conf, err := apk.Config()
	if err != nil {
		return errors.Wrap(err, "could read apk config")
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

	targetCrashed, err := c.Inspector.TargetCrashed(ctx)
	if err != nil {
		return nil
	}
	go func() {
		defer targetCrashed.Close()

		_, err := targetCrashed.Recv()
		if err != nil {
			if cdp.ErrorCause(err) != ctx.Err() {
				log.Printf("targetCrashed.Recv(): %v", err)
			}
			return
		}
		log.Println("Target crashed!")
		errc <- errors.New("target crashed")
		cancel()
	}()

	detached, err := c.Inspector.Detached(ctx)
	if err != nil {
		return nil
	}
	go func() {
		defer detached.Close()

		ev, err := detached.Recv()
		if err != nil {
			if cdp.ErrorCause(err) != ctx.Err() {
				log.Printf("detached.Recv(): %v", err)
			}
			return
		}
		log.Printf("Inspector detached: %v!", ev.Reason)
		errc <- fmt.Errorf("inspector detached: %v", ev.Reason)
		cancel()
	}()

	app := appSlice(apps).Find(conf.General.Package, conf.General.Architecture)

	absPath, err := filepath.Abs(apk.Path())
	if err != nil {
		return err
	}

	for _, fn := range []func() error{
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
		res, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{NodeID: o})
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
