package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/bwmarrin/discordgo"
)

var (
	cfg       config
	dg        *discordgo.Session
	loginTime time.Time
)

func main() {
	a := app.New()
	w := a.NewWindow("seld got")
	setupWindow(w)

	fmt.Println("trying to find the config file...")
	label := centeredLabel("Trying to find the config file...")
	content := container.NewCenter(label)
	w.SetContent(content)

	go func() {
		err := cfg.loadConfig()
		switch {
		case os.IsNotExist(err):
			// create config
			fmt.Println("config file not found, creating one")
			label.SetText("Config file not found, creating one...")

			tokenDone := make(chan struct{})
			prefixDone := make(chan struct{})

			handleTokenInput(content, label, func() {
				handlePrefixInput(content, label, func() {
					err := cfg.createConfig()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("config file created")
					close(prefixDone)
				})
				close(tokenDone)
			})

			<-tokenDone
			<-prefixDone
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Println("config file found")
		}

		if !cfg.isValid() {
			log.Fatal("config file is invalid")
			return
		}

		fmt.Println("testing token...")
		label.SetText("Testing token...")
		if err = testToken(); err != nil {
			log.Fatal(err)
		}
		label.SetText("Token is valid, connecting to Discord...")

		loginTime = time.Now()
		label.SetText("Connected")
		connectToDiscord()

		prepareCommands()

		dg.AddHandler(handleMessageCreate)
	}()

	w.ShowAndRun()
}

func setupWindow(w fyne.Window) {
	w.SetMaster()
	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)
	w.CenterOnScreen()
}

func handleTokenInput(content *fyne.Container, label *widget.Label, onComplete func()) {
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
			centeredLabel("Please enter your token:"),
			layout.NewSpacer(),
			tokenEntry,
			layout.NewSpacer(),
			saveButtonContainer,
		),
	}
	content.Refresh()
}

func handlePrefixInput(content *fyne.Container, label *widget.Label, onComplete func()) {
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

func centeredLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{})
}
