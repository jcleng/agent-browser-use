package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestExtensionFlagArgs(t *testing.T) {
	// Off: no extension flags at all, so chromedp's --disable-extensions stays in place.
	if got := extensionFlagArgs(nil); got != nil {
		t.Fatalf("extensionFlagArgs(nil) = %v, want nil", got)
	}

	a := t.TempDir()
	b := t.TempDir()
	got := extensionFlagArgs([]string{a, b})
	want := [][2]string{
		// Without --disable-extensions-except, Chromium ignores --load-extension.
		{"disable-extensions-except", a + "," + b},
		{"load-extension", a + "," + b},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("extensionFlagArgs = %v, want %v", got, want)
	}
}

func TestExtensionDirs(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	notADir := filepath.Join(a, "manifest.json")
	if err := os.WriteFile(notADir, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Blank entries and folders that do not exist (or are files) are dropped.
	got := extensionDirs(" " + a + " , ," + b + ", " + filepath.Join(a, "missing") + "," + notADir + " ")
	if !slices.Equal(got, []string{a, b}) {
		t.Fatalf("extensionDirs = %v, want %v", got, []string{a, b})
	}
	if dirs := extensionDirs(""); dirs != nil {
		t.Errorf("extensionDirs(\"\") = %v, want nil", dirs)
	}
}

func TestExtensionDirsRelative(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "ext"), 0o755); err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(wd) })

	// Chrome is launched from the app folder, so a relative folder must become absolute.
	got := extensionDirs("ext")
	if len(got) != 1 || !filepath.IsAbs(got[0]) || !strings.HasSuffix(got[0], "ext") {
		t.Fatalf("extensionDirs(\"ext\") = %v, want an absolute path to ext", got)
	}
}

func TestExtensionsEnabled(t *testing.T) {
	// Default: off, whatever the folder variable says (images set the folder).
	for _, env := range []string{"", "0", "false", "no", "off", "banana"} {
		if extensionsEnabled(env, false, false) {
			t.Errorf("extensionsEnabled(%q, no flag) = true, want false", env)
		}
	}
	// AGENT_BROWSER_USE_EXTENSIONS turns them on.
	for _, env := range []string{"1", "true", "TRUE", "yes", "on", " On "} {
		if !extensionsEnabled(env, false, false) {
			t.Errorf("extensionsEnabled(%q, no flag) = false, want true", env)
		}
	}
	// An explicit -extensions wins over the environment, in both directions.
	if extensionsEnabled("1", true, false) {
		t.Error("-extensions=false must win over the env var")
	}
	if !extensionsEnabled("0", true, true) {
		t.Error("-extensions=true must win over the env var")
	}
}
