package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/bwmarrin/discordgo"
)

type Runner interface {
	run() error
}

type UIMode struct {
	cfg       *config
	dg        *discordgo.Session
	window    fyne.Window
	label     *widget.Label
	loginTime time.Time
}

type CLIMode struct {
	cfg       *config
	dg        *discordgo.Session
	loginTime time.Time
}

func main() {
	noUI := len(os.Args) > 1 && (os.Args[1] == "--noui" || os.Args[1] == "--no-ui")

	var cfg config
	var runner Runner

	if noUI {
		runner = &CLIMode{cfg: &cfg}
	} else {
		runner = NewUI(&cfg)
	}

	if err := runner.run(); err != nil {
		log.Fatal(err)
	}
}

func NewUI(cfg *config) *UIMode {
	a := app.New()
	w := a.NewWindow("self got")
	setupWindow(w)

	return &UIMode{
		cfg:    cfg,
		window: w,
		label:  centeredLabel("Trying to find the config file..."),
	}
}

func (u *UIMode) run() error {
	content := container.NewCenter(u.label)
	u.window.SetContent(content)

	go func() {
		err := u.cfg.loadConfig()
		switch {
		case os.IsNotExist(err):
			fmt.Println("config file not found, creating one")
			u.label.SetText("Config file not found, creating one...")

			tokenDone := make(chan struct{})
			prefixDone := make(chan struct{})

			handleTokenInput(content, u.label, u.cfg, func() {
				handlePrefixInput(content, u.label, u.cfg, func() {
					err := u.cfg.createConfig()
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

		if !u.cfg.isValid() {
			log.Fatal("config file is invalid")
			return
		}

		fmt.Println("testing token...")
		u.label.SetText("Testing token...")

		u.dg, err = createSession(u.cfg)
		if err != nil {
			log.Fatal(err)
		}

		u.label.SetText("Token is valid, connecting to Discord...")

		if err := connectToWS(u.dg); err != nil {
			log.Fatal(err)
		}

		u.loginTime = time.Now()
		u.label.SetText("Connected, close this window to exit")

		prepareCommands()
		u.dg.AddHandler(messageCreateWrapper(u.cfg))
		finalWindow(u.window)
	}()

	u.window.ShowAndRun()
	return nil
}

func (c *CLIMode) run() error {
	fmt.Println("running in no-ui mode")

	var err error
	if err = c.cfg.loadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	if !c.cfg.isValid() {
		return fmt.Errorf("invalid config")
	}

	c.dg, err = createSession(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	if err = connectToWS(c.dg); err != nil {
		return fmt.Errorf("failed to connect to websocket: %v", err)
	}

	c.loginTime = time.Now()
	prepareCommands()
	c.dg.AddHandler(messageCreateWrapper(c.cfg))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	return nil
}

func setupWindow(w fyne.Window) {
	w.SetMaster()
	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)
	w.CenterOnScreen()
}

func finalWindow(w fyne.Window) {
	commandNames := getCommandNames()

	// Create and populate a list with the command names
	list := widget.NewList(
		func() int {
			return len(commandNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(commandNames[i])
		},
	)

	helpLabel := widget.NewLabel("Select a command to see its help")

	// Update the help text when a command is selected
	list.OnSelected = func(id widget.ListItemID) {
		helpLabel.SetText(commMap[commandNames[id]].Help)
	}

	// Horizontal split to hold the list and the help text
	split := container.NewHSplit(list, helpLabel)
	split.Offset = 0.3

	w.SetContent(split)
}

func centeredLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{})
}
