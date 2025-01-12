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
	spacer := layout.NewSpacer()

	saveButton := widget.NewButton("Save", func() {
		if tokenEntry.Text == "" {
			return
		}
		cfg.Token = tokenEntry.Text
		content.Objects = []fyne.CanvasObject{label}
		content.Refresh()
		onComplete()
	})

	buttonContainer := container.NewCenter(saveButton)
	buttonContainer.Resize(fyne.NewSize(120, saveButton.MinSize().Height))

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(
			centeredLabel("Config file not found, creating one..."),
			spacer,
			tokenEntry,
			spacer,
			buttonContainer,
		),
	}
	content.Refresh()
}

func handlePrefixInput(content *fyne.Container, label *widget.Label, cfg *config, onComplete func()) {
	prefixEntry := widget.NewEntry()

	saveButton := widget.NewButton("Save", func() {
		cfg.Prefix = "\\"
		if prefixEntry.Text != "" {
			cfg.Prefix = prefixEntry.Text
		}
		content.Objects = []fyne.CanvasObject{label}
		content.Refresh()
		onComplete()
	})

	buttonContainer := container.NewCenter(saveButton)
	buttonContainer.Resize(fyne.NewSize(120, saveButton.MinSize().Height))

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(
			centeredLabel("Configure Prefix (default: \\)"),
			layout.NewSpacer(),
			prefixEntry,
			layout.NewSpacer(),
			buttonContainer,
		),
	}
	content.Refresh()
}
