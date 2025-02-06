package main

type ClipboardError string

const ErrClipboardNotSupported = ClipboardError("clipboard not supported")

func (e ClipboardError) Error() string {
	return string(e)
}

type Clipboard interface {
	// Copies the given text to the system clipboard
	Copy(text string) error
}

func getPlatformClipboard(osName string) Clipboard {
	switch osName {
	case "darwin":
		return DarwinClipboard{}

	case "windows":
		return WindowsClipboard{}

	default:
		return nil
	}
}
