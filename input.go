package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func handleTokenInput(content *fyne.Container, label *widget.Label, cfg *config, onComplete func()) {
	tokenEntry := widget.NewPasswordEntry()
	saveButton := widget.NewButton("Save", func() {
		if tokenEntry.Text == "" {
			fmt.Println("invalid token")
			return
		}
		cfg.Token = tokenEntry.Text
		content.Objects = []fyne.CanvasObject{label}
		content.Refresh()
		onComplete()
	})

	// Fixed button width so it doesnt grow
	saveButtonContainer := container.NewCenter(saveButton)
	saveButtonContainer.Resize(fyne.NewSize(100, saveButton.MinSize().Height))

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(
			centeredLabel("Enter your token:"),
			layout.NewSpacer(),
			tokenEntry,
			layout.NewSpacer(),
			saveButtonContainer,
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

	// Fixed button width so it doesnt grow
	saveButtonContainer := container.NewCenter(saveButton)
	saveButtonContainer.Resize(fyne.NewSize(100, saveButton.MinSize().Height))

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(
			centeredLabel("Please enter your prefix, default is \"\\\" (e.g. \\help):"),
			layout.NewSpacer(),
			prefixEntry,
			layout.NewSpacer(),
			saveButtonContainer,
		),
	}
	content.Refresh()
}
