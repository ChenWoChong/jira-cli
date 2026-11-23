package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TextData is the data to be shown in text layout.
type TextData string

// Text is the text view layout.
type Text struct {
	screen  *Screen
	painter *tview.TextView
}

// NewText constructs a new text view layout.
func NewText() *Text {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	tv := Text{screen: NewScreen()}
	tv.init()

	return &tv
}

// Render renders the text layout.
func (tv *Text) Render(td TextData) error {
	tv.painter.SetText(tview.TranslateANSI(string(td)))
	return tv.screen.Paint(tv.painter)
}

func (tv *Text) init() {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	view.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			tv.screen.Stop()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			tv.screen.Stop()
			return nil
		}
		return event
	})

	tv.painter = view
}

// ViewText displays text in the built-in scrollable viewer.
func ViewText(out string) error {
	return NewText().Render(TextData(out))
}
