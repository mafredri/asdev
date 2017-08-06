package apkg

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	changelogText = "1.0.0:\n- Initial release\n"
	descText      = "Test Package\n\nThis is a test.\n"
)

func TestOpen(t *testing.T) {
	apks := []string{
		"test_1.0.0_arm.apk",
		"test_1.0.0_i386.apk",
		"test_1.0.0_x86-64.apk",
	}
	for _, apk := range apks {
		f, err := Open(path.Join("testdata", apk))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
	}
}

func TestConfig(t *testing.T) {
	apks := []struct {
		arch string
		apk  string
	}{
		{"arm", "test_1.0.0_arm.apk"},
		{"i386", "test_1.0.0_i386.apk"},
		{"x86-64", "test_1.0.0_x86-64.apk"},
	}
	for _, tt := range apks {
		f, err := Open(path.Join("testdata", tt.apk))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		c, err := f.Config()
		if err != nil {
			t.Error(err)
		}

		if c.General.Package != "test" {
			t.Errorf("Config.General.Package: got %v, want %v", c.General.Package, "test")
		}
		if c.General.Architecture != tt.arch {
			t.Errorf("Config.General.Architecture: got %v, want %v", c.General.Architecture, tt.arch)
		}
	}
}

func TestChangelog(t *testing.T) {
	apks := []string{
		"test_1.0.0_arm.apk",
		"test_1.0.0_i386.apk",
		"test_1.0.0_x86-64.apk",
	}
	for _, apk := range apks {
		f, err := Open(path.Join("testdata", apk))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		chlog := f.Changelog()
		b, err := ioutil.ReadAll(chlog)
		if err != nil {
			t.Error(err)
		}
		defer chlog.Close()

		if diff := cmp.Diff(string(b), changelogText); diff != "" {
			t.Errorf("Changlog output differs: (-got +want)\n%s", diff)
		}
	}
}

func TestDescription(t *testing.T) {
	apks := []string{
		"test_1.0.0_arm.apk",
		"test_1.0.0_i386.apk",
		"test_1.0.0_x86-64.apk",
	}
	for _, apk := range apks {
		f, err := Open(path.Join("testdata", apk))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		desc := f.Description()
		b, err := ioutil.ReadAll(desc)
		if err != nil {
			t.Error(err)
		}
		defer desc.Close()

		if diff := cmp.Diff(string(b), descText); diff != "" {
			t.Errorf("Changlog output differs: (-got +want)\n%s", diff)
		}
	}
}
