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

type options struct {
	categories Categories
	tags       []string
	beta       bool
	icon       string
}

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

		list = kingpin.Command("list", "List and show status of all apps and app updates")

		update     = kingpin.Command("update", "Update apps by uploading one or more APK(s)")
		updateCats = update.Flag("category", "(NOT IMPLEMENTED) Change categorie(s)").Short('c').Enums(categories...)
		updateTags = update.Flag("tag", "(NOT IMPLEMENTED) Change tag(s)").Short('t').HintOptions("multimedia", "web").Strings()
		updateBeta = update.Flag("beta", "(NOT IMPLEMENTED) Beta app").Short('b').Bool()
		updateIcon = update.Flag("icon", "(NOT IMPLEMENTED) Change icon (256x256)").Short('i').ExistingFile()
		updateAPKs = update.Arg("APKs", "APK(s) to update").Required().ExistingFiles()

		create     = kingpin.Command("create", "(NOT IMPLEMENTED) Submit a new application by uploading one or more APK(s)")
		createCats = create.Flag("category", "Set categorie(s)").Short('c').Required().Enums(categories...)
		createTags = create.Flag("tag", "Set tag(s)").Short('t').HintOptions("multimedia", "web").Required().Strings()
		createBeta = create.Flag("beta", "Set app to beta status").Short('b').Bool()
		createIcon = create.Flag("icon", "Set icon (256x256)").Short('i').ExistingFile()
		createAPKs = create.Arg("APKs", "APK(s) to create").Required().ExistingFiles()
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

	case list.FullCommand():
		checkLogin(*username, *password)

		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cancel()
		}()

		do := func(devt *devtool.DevTools, apps ...App) error {
			drawAppTable(apps...)
			return nil
		}
		err := run(ctx, *verbose, !*noHeadless, *browser, *username, *password, nil, do)
		if err != nil {
			log.Fatal(err)
		}

	case create.FullCommand():
		checkLogin(*username, *password)

		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cancel()
		}()

		apks := openAPKs(*createAPKs...)

		log.Println(*createCats, *createTags, *createBeta, *createIcon, *createAPKs, apks)
		fmt.Println("create is not implemented yet!")

		do := func(devt *devtool.DevTools, apps ...App) error {
			errc := make(chan chan error, len(apks))
			for _, apk := range apks {
				errc2 := make(chan error, 1)
				go uploadAPK(ctx, *verbose, devt, errc2, apps, apk, options{
					categories: parseCategories(*createCats...),
					tags:       *createTags,
					beta:       *createBeta,
					icon:       *createIcon,
				})
				defer apk.Close()
				errc <- errc2
			}
			close(errc)

			for e := range errc {
				select {
				case err := <-e:
					if err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		}

		if err := run(ctx, *verbose, !*noHeadless, *browser, *username, *password, apks, do); err != nil {
			log.Fatal(err)
		}

	case update.FullCommand():
		checkLogin(*username, *password)

		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cancel()
		}()

		apks := openAPKs(*updateAPKs...)

		do := func(devt *devtool.DevTools, apps ...App) error {
			errc := make(chan chan error, len(apks))
			for _, apk := range apks {
				errc2 := make(chan error, 1)
				go updateAPK(ctx, *verbose, devt, errc2, apps, apk, options{
					categories: parseCategories(*updateCats...),
					tags:       *updateTags,
					beta:       *updateBeta,
					icon:       *updateIcon,
				})
				defer apk.Close()
				errc <- errc2
			}
			close(errc)

			for e := range errc {
				select {
				case err := <-e:
					if err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		}

		if err := run(ctx, *verbose, !*noHeadless, *browser, *username, *password, apks, do); err != nil {
			log.Fatal(err)
		}
	}
}

func parseCategories(c ...string) []string {
	var names []string
	for _, s := range c {
		names = append(names, category(s).Name())
	}
	return names
}

func checkLogin(username, password string) {
	if username == "" || password == "" {
		fmt.Println("error: username or password is missing, use cli flag or set in environment")
		os.Exit(1)
	}
}

func openAPKs(name ...string) []*apkg.File {
	var apks []*apkg.File
	for _, av := range name {
		apk, err := apkg.Open(av)
		if err != nil {
			fmt.Printf("error: could open apk %q: %v\n", av, err)
			os.Exit(1)
		}
		apks = append(apks, apk)
	}
	return apks
}

var (
	optionTagRe = regexp.MustCompile("</?option( [^>]+)?>")
)

func run(ctx context.Context, verbose, headless bool, chromeBin string, username, password string, apks []*apkg.File, do func(*devtool.DevTools, ...App) error) error {
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

	return do(devt, apps...)
}

func abortOnDetachOrCrash(ctx context.Context, ic cdp.Inspector, abort func(err error)) error {
	targetCrashed, err := ic.TargetCrashed(ctx)
	if err != nil {
		return nil
	}

	detached, err := ic.Detached(ctx)
	if err != nil {
		return nil
	}

	go func() {
		defer targetCrashed.Close()
		defer detached.Close()

		select {
		case <-targetCrashed.Ready():
			_, err := targetCrashed.Recv()
			if err != nil {
				if cdp.ErrorCause(err) != ctx.Err() {
					log.Printf("targetCrashed.Recv(): %v", err)
				}
				return
			}
			abort(errors.New("target crashed"))

		case <-detached.Ready():
			ev, err := detached.Recv()
			if err != nil {
				if cdp.ErrorCause(err) != ctx.Err() {
					log.Printf("detached.Recv(): %v", err)
				}
				return
			}
			abort(fmt.Errorf("inspector detached: %v", ev.Reason))
		}
	}()

	return nil
}
