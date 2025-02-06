package main

import "testing"

func TestPlatformSpecificClipboardImplementation(t *testing.T) {
	t.Run("when OS is darwin, clipboard is instance of DarwinClipboard", func(t *testing.T) {
		want := DarwinClipboard{}
		got := getPlatformClipboard("darwin")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("when OS is linux, clipboard is not supported", func(t *testing.T) {
		got := getPlatformClipboard("linux")

		if got != nil {
			t.Errorf("got %q, want nil", got)
		}
	})

	t.Run("when OS is windows, clipboard is an instance of WindowsClipboard", func(t *testing.T) {
		want := WindowsClipboard{}
		got := getPlatformClipboard("windows")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
