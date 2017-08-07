package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"

	"github.com/mafredri/asdev/apkg"

	"github.com/alecthomas/kingpin"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
)

const (
	defaultBrowser = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
)

// Pull in latest categories from App Central.
//go:generate go run cmd/catgen/main.go main category.go

func main() {
	var (
		username = kingpin.Flag("username", "Username (for login)").Short('u').Envar("ASDEV_USERNAME").String()
		password = kingpin.Flag("password", "Password (for login)").Short('p').Envar("ASDEV_PASSWORD").String()
		browser  = kingpin.Flag("browser", "Path to Chrome or Chromium executable").
				Default(defaultBrowser).Envar("ASDEV_BROWSER").String()
		noHeadless = kingpin.Flag("no-headless", "Disable (Chrome) headless mode").Bool()
		timeout    = kingpin.Flag("timeout", "Command timeout").Default("10m").Duration()
		verbose    = kingpin.Flag("verbose", "Verbose mode").Short('v').Bool()

		show           = kingpin.Command("show", "Show additional information")
		showCategories = show.Command("categories", "Show all available categories")

		update     = kingpin.Command("update", "Update apps by uploading one or multiple APK(s)")
		updateAPKs = update.Arg("APKs", "APK(s) to update").Required().ExistingFiles()
	)

	// Provide help via short flag as well.
	kingpin.HelpFlag.Short('h')

	switch kingpin.Parse() {
	case showCategories.FullCommand():
		maxlen := 0
		for _, c := range categories {
			if len(c) > maxlen {
				maxlen = len(c)
			}
		}
		format := fmt.Sprintf("  %%-%ds(%%s)\n", maxlen+1)
		fmt.Printf("Available categories:\n\n")
		for _, c := range categories {
			fmt.Printf(format, c, category(c).Name())
		}
	case update.FullCommand():
		if *username == "" || *password == "" {
			fmt.Println("error: username or password is missing, use cli flag or set in environment")
			os.Exit(1)
		}

		var apks []*apkg.File
		for _, av := range *updateAPKs {
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

func abortOnDetachOrCrash(ctx context.Context, ic cdp.Inspector, abort func(err error)) error {
	targetCrashed, err := ic.TargetCrashed(ctx)
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
		abort(errors.New("target crashed"))
	}()

	detached, err := ic.Detached(ctx)
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
		abort(fmt.Errorf("inspector detached: %v", ev.Reason))
	}()

	return nil
}
