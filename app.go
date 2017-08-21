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
	Update     *App      `json:"update"`
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
	const appTableRows = '#product-table:first-of-type tr';
	const appSelect = [
		{name: 'package', sel: 'td:nth-child(2)'},
		{name: 'name', sel: 'td:nth-child(3)'},
		{name: 'arch', sel: 'td:nth-child(4)'},
		{name: 'version', sel: 'td:nth-child(6)'},
		{name: 'categories', sel: 'td:nth-child(7)'},
		{name: 'lastUpdate', sel: 'td:nth-child(9)'},
		{name: 'status', sel: 'td:nth-child(10)'},
		{name: 'id', sel: 'a[href*="app/updateApp?id="]'},
	];
	const appSelectNames = appSelect.map(r => r.name);
	const appSelectCols = appSelect.map(r => r.sel).join(',');

	const updateTableRows = '#product-table:nth-of-type(2) tr';
	const updateSelect = [
		{name: 'package', sel: 'td:nth-child(2)'},
		{name: 'name', sel: 'td:nth-child(3)'},
		{name: 'arch', sel: 'td:nth-child(4)'},
		{name: 'version', sel: 'td:nth-child(5)'},
		{name: 'categories', sel: 'td:nth-child(6)'},
		{name: 'lastUpdate', sel: 'td:nth-child(8)'},
		{name: 'status', sel: 'td:nth-child(9)'},
	];
	const updateSelectNames = updateSelect.map(r => r.name);
	const updateSelectCols = updateSelect.map(r => r.sel).join(',');

	const zipObject = (entries, cols) => Object.assign({}, ...entries.map((v, i) => ({[cols[i]]: v})));

	function parseApp(el, select, names) {
		const app = zipObject(columnValues(el, select), names);
		app.categories = app.categories.split('\n').map(s => s.trim()).filter(s => s);
		// Convert to Go time format (UTC).
		app.lastUpdate = app.lastUpdate.replace(' ', 'T') + 'Z';
		app.updated = !!el.querySelector('img[src*="images/forms/update.png"]');
		app.beta = !!el.querySelector('img[src*="images/beta.jpg"]');
		return app;
	}

	function setUpdatedApp(apps, update) {
		const app = apps.find(app => app.updated && app.package === update.package && app.arch === update.arch);
		if (!app) {
			throw new Error('could not find app: ' + update.package);
		}
		app.update = update;
	}

	function columnValues(parent, sel) {
		return Array.from(parent.querySelectorAll(sel))
			.map(el => {
				if (el.href) {
					return parseInt(el.href.replace(/.*id=/, ''), 10);
				}
				if (el.src) {
					return true;
				}
				return el.textContent.trim();
			});
	}

	const apps = Array.from(document.querySelectorAll(appTableRows))
		.slice(1)
		.map(td => parseApp(td, appSelectCols, appSelectNames));

	// Decorate app data with upated apps.
	Array.from(document.querySelectorAll(updateTableRows))
		.slice(1)
		.map(td => parseApp(td, updateSelectCols, updateSelectNames))
		.forEach(update => setUpdatedApp(apps, update));

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
		if a.Update != nil {
			table.Append([]string{
				"",
				a.Update.Name,
				a.Update.Arch,
				a.Update.Version,
				strings.Join(a.Update.Categories, ", "),
				a.Update.LastUpdate.Format("2006-01-02 15:04:05"),
				a.Update.Status,
			})
		}
	}

	table.Render()
}
