package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func handleTokenInput(content *fyne.Container, label *widget.Label, cfg *config, onComplete func()) {
	tokenEntry := widget.NewPasswordEntry()
	tokenEntry.SetPlaceHolder("token")

	saveButton := widget.NewButton("Save", func() {
		if tokenEntry.Text == "" {
			return
		}
		cfg.Token = tokenEntry.Text
		content.Objects = []fyne.CanvasObject{label}
		content.Refresh()
		onComplete()
	})

	saveButton.Importance = widget.HighImportance

	form := container.NewVBox(
		widget.NewLabelWithStyle("Creating config", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Enter your Discord token to connect the selfbot:"),
		container.NewPadded(tokenEntry),
	)

	paddedContent := container.NewPadded(
		container.NewVBox(
			form,
			container.NewHBox(
				layout.NewSpacer(),
				saveButton,
				layout.NewSpacer(),
			),
		),
	)

	content.Objects = []fyne.CanvasObject{paddedContent}
	content.Refresh()
}

func handlePrefixInput(content *fyne.Container, label *widget.Label, cfg *config, onComplete func()) {
	prefixEntry := widget.NewEntry()
	prefixEntry.SetPlaceHolder("\\")
	prefixEntry.Resize(fyne.NewSize(100, 36))

	saveButton := widget.NewButton("Save", func() {
		cfg.Prefix = "\\"
		if prefixEntry.Text != "" {
			cfg.Prefix = prefixEntry.Text
		}
		content.Objects = []fyne.CanvasObject{label}
		content.Refresh()
		onComplete()
	})

	saveButton.Importance = widget.HighImportance
	saveButton.Resize(fyne.NewSize(120, 36))

	form := container.NewVBox(
		widget.NewLabelWithStyle("Command Prefix", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Choose a prefix for bot commands (default is \\):"),
		container.NewPadded(prefixEntry),
	)

	paddedContent := container.NewPadded(
		container.NewVBox(
			form,
			container.NewHBox(
				layout.NewSpacer(),
				saveButton,
				layout.NewSpacer(),
			),
		),
	)

	content.Objects = []fyne.CanvasObject{paddedContent}
	content.Refresh()
}
