package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	devtoolsPortFile = "DevToolsActivePort"
)

type chromeProc struct {
	port int
	cmd  *exec.Cmd
}

func (c *chromeProc) Close() error {
	return c.cmd.Process.Kill()
}

func startChrome(ctx context.Context, name string, tmpdir string, headless bool) (*chromeProc, error) {
	args := []string{
		"--no-default-browser-check",
		"--no-first-run",
		"--user-data-dir=" + tmpdir,
		"--remote-debugging-port=0", // Automatic port selection by OS.
		"about:blank",               // No URL.
	}
	if headless {
		args = append(args, "--headless", "--disable-gpu")
	}
	cmd := exec.CommandContext(ctx, name, args...)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Give a reasonable amount of time for the debug port to become active.
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)

	go func() {
		defer cancel() // Cancel ping if process exits prematurely.
		cmd.Wait()
	}()

	c := &chromeProc{cmd: cmd}
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return nil, errors.Wrap(ctx.Err(), "timeout waiting for debugging port")
			}
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
			// Probe the DevTools port file for the port number chosen by OS.
			ps, err := ioutil.ReadFile(filepath.Join(tmpdir, devtoolsPortFile))
			if err != nil {
				continue
			}

			// Filter out the debugger URL and only parse the port number.
			ps = bytes.SplitN(ps, []byte{'\n'}, 2)[0]

			c.port, err = strconv.Atoi(string(ps))
			if err != nil {
				// Assume file does not exist / have content yet.
				continue
			}

			return c, nil
		}
	}
}
