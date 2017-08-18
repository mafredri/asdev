package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/olekukonko/tablewriter"
)

// App represents an app in the ASUSTOR Developer Corner.
type App struct {
	ID         int       `json:"id"`
	Package    string    `json:"package"`
	Name       string    `json:"name"`
	Arch       string    `json:"arch"`
	Version    string    `json:"version"`
	Categories []string  `json:"categories"`
	LastUpdate time.Time `json:"lastUpdate"`
	Status     string    `json:"status"`
	Beta       bool      `json:"beta"`
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
	const selectColumns = 'td:nth-child(2), td:nth-child(3), td:nth-child(4), td:nth-child(6), td:nth-child(7), td:nth-child(9), td:nth-child(10), a[href*="app/updateApp?id="]';
	const colNames = ['package', 'name', 'arch', 'version', 'categories', 'lastUpdate', 'status', 'id'];
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
			// Convert to Go time format (UTC).
			obj.lastUpdate = obj.lastUpdate.replace(' ', 'T') + 'Z';
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

func drawAppTable(apps ...App) {
	tableHeader := []string{"Package", "Name", "Arch", "Version", "Categories", "Last update", "Status"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeader)
	table.SetAutoWrapText(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if width, _, err := terminal.GetSize(int(os.Stdin.Fd())); err == nil {
		table.SetColWidth(width / len(tableHeader))
	}

	for _, a := range apps {
		table.Append([]string{
			a.Package,
			a.Name,
			a.Arch,
			a.Version,
			strings.Join(a.Categories, ", "),
			a.LastUpdate.Format("2006-01-02 15:04:05"),
			a.Status,
		})
	}

	table.Render()
}
