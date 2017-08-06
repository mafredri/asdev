package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
)

// App represents an app in the ASUSTOR Developer Corner.
type App struct {
	ID         int      `json:"id"`
	Package    string   `json:"package"`
	Arch       string   `json:"arch"`
	Categories []string `json:"categories"`
	Status     string   `json:"status"`
	Beta       bool     `json:"beta"`
}

// HasCategory returns true if the app has the category,
// uses case insensitive comparison.
func (a *App) HasCategory(cat string) bool {
	cat = strings.ToLower(cat)
	for _, c := range a.Categories {
		if cat == strings.ToLower(c) {
			return true
		}
	}
	return false
}

// UpdateURL returns the URL for updating the App.
func (a *App) UpdateURL() string {
	return fmt.Sprintf("http://developer.asustor.com/app/updateApp?id=%d", a.ID)
}

type appSlice []App

func (as appSlice) Find(name, arch string) (a App) {
	for _, app := range as {
		if app.Package == name && app.Arch == arch {
			return app
		}
	}
	return a
}

// Parse application data from the table.
const jsParseApps = `
	const selectRows = '#product-table:first-of-type tr';
	const selectColumns = 'td:nth-child(2), td:nth-child(4), td:nth-child(7), td:nth-child(9), td:nth-child(10), a[href*="app/updateApp?id="]';
	const colNames = ['package', 'arch', 'categories', 'lastUpdate', 'status', 'id'];
	const zipObject = (entries, cols) => Object.assign({}, ...entries.map((v, i) => ({[cols[i]]: v})));

	const apps = Array.from(document.querySelectorAll(selectRows))
		.slice(1)
		.map(td => {
			const cols = Array.from(td.querySelectorAll(selectColumns))
				.map(el => {
					if (el.href) {
						return parseInt(el.href.replace(/.*id=/, ''), 10);
					}
					return el.textContent.trim();
				});
			const obj = zipObject(cols, colNames);
			obj.categories = obj.categories.split('\n').map(s => s.trim()).filter(s => s);
			obj.beta = !!td.querySelector('img[src*="beta"]');
			return obj;
		});

	// Filter out the older entries if there are duplicates.
	apps.filter(a => !apps.find(b => b !== a && a.package === b.package && a.arch === b.arch && a.lastUpdate < b.lastUpdate));
`

func getApps(ctx context.Context, c *cdp.Client) ([]App, error) {
	err := navigate(ctx, c.Page, "http://developer.asustor.com/app/mgt", 10*time.Second)
	if err != nil {
		return nil, err
	}

	evalArgs := runtime.NewEvaluateArgs(jsParseApps).SetReturnByValue(true)
	eval, err := c.Runtime.Evaluate(ctx, evalArgs)
	if err != nil {
		return nil, err
	}

	var apps []App
	err = json.Unmarshal(eval.Result.Value, &apps)
	return apps, err
}
