package main

import "testing"

func TestPlatformSpecificClipboardImplementation(t *testing.T) {
	t.Run("when OS is darwin", func(t *testing.T) {
		want := DarwinClipboard{}
		got := getPlatformClipboard("darwin")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("when OS is linux", func(t *testing.T) {
		got := getPlatformClipboard("linux")

		if got != nil {
			t.Errorf("got %q, want nil", got)
		}
	})

	t.Run("when OS is windows", func(t *testing.T) {
		got := getPlatformClipboard("windows")

		if got != nil {
			t.Errorf("got %q, want nil", got)
		}
	})
}
