package main

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	g "xabbo.b7c.io/goearth"
	"xabbo.b7c.io/goearth/shockwave/out"
)

type UI struct {
	app                 fyne.App
	window              fyne.Window
	messageList         *widget.List
	startButton         *widget.Button
	stopButton          *widget.Button
	clearButton         *widget.Button
	intervalLabel       *widget.Label
	stopOnTradeCheckbox *widget.Check
}

type AutoTyper struct {
	ext          *g.Ext
	intervalSec  int
	messages     []string
	ticker       *time.Ticker
	stopTickerCh chan bool
	mutex        sync.Mutex
	ui           *UI
}

func NewAutoTyper(ext *g.Ext) *AutoTyper {
	return &AutoTyper{
		ext:         ext,
		intervalSec: 120, // Default interval to 120 seconds
		ui: &UI{
			app: app.New(),
		},
	}
}

func (at *AutoTyper) Start() {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if len(at.messages) == 0 {
		fmt.Println("No messages to send.")
		return
	}

	if at.ticker == nil {
		if at.ext == nil {
			fmt.Println("Extension not initialized.")
			return
		}

		at.sendMessage(0)

		at.ticker = time.NewTicker(time.Duration(at.intervalSec) * time.Second)
		at.stopTickerCh = make(chan bool)
		go at.runTickerLoop()

		at.updateUI()
	}
}

func (at *AutoTyper) Stop() {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if at.stopTickerCh != nil {
		at.stopTickerCh <- true
		at.ticker.Stop()
		at.ticker = nil
		at.stopTickerCh = nil

		at.updateUI()
		fmt.Println("AutoTyper stopped.")
	}
}

func (at *AutoTyper) SetInterval(seconds int) {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if seconds < 1 {
		seconds = 1
	}
	at.intervalSec = seconds
	at.updateUI()
}

func (at *AutoTyper) AddMessage(message string) error {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	at.messages = append(at.messages, message)
	at.updateMessageList() // Refresh the message list
	return nil
}

func (at *AutoTyper) GetMessages() []string {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	return at.messages
}

func (at *AutoTyper) ClearMessages() {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	at.messages = []string{}
	at.updateMessageList() // Refresh the message list
}

func (at *AutoTyper) sendMessage(index int) {
	at.ext.Send(out.SHOUT, at.messages[index])
	fmt.Println(at.messages[index])
}

func (at *AutoTyper) runTickerLoop() {
	index := 1
	for {
		select {
		case <-at.ticker.C:
			if index >= len(at.messages) {
				index = 0
			}

			at.sendMessage(index)
			index++

		case <-at.stopTickerCh:
			return
		}
	}
}

func (at *AutoTyper) Run() {
	ui := at.ui
	ui.window = ui.app.NewWindow("AutoTyper by Nanobyte")
	ui.window.Resize(fyne.NewSize(500, 500)) // Increased window size

	// Initialize the message list
	ui.messageList = widget.NewList(
		func() int {
			return len(at.messages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(at.messages[id])
		},
	)

	// Initialize buttons
	ui.startButton = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {
		at.Start()
	})

	ui.stopButton = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		at.Stop()
	})

	ui.clearButton = widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		at.ClearMessages()
	})

	ui.intervalLabel = widget.NewLabel(fmt.Sprintf("Interval: %d seconds", at.intervalSec))

	// Build UI sections
	messageSection := at.buildMessagesSection()
	controlSection := at.buildControlSection()
	settingsSection := at.buildSettingsSection()

	// Layout all sections
	content := container.NewVBox(
		messageSection,
		controlSection,
		settingsSection, // at.intervalLabel is now part of the settingsSection
	)

	// Apply padding
	paddedContent := container.New(layout.NewPaddedLayout(), content)

	ui.window.SetContent(paddedContent)
	icon, _ := fyne.LoadResourceFromPath("./app_icon.ico")
	ui.app.SetIcon(icon)
	// Initial UI update
	at.updateUI()

	ui.window.ShowAndRun()
}

func (at *AutoTyper) buildMessagesSection() fyne.CanvasObject {
	ui := at.ui
	// Messages list
	messageListLabel := widget.NewLabel("Messages:")

	// Text field for adding messages
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Enter your message")

	// Add button to add a message to the AutoTyper
	addButton := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		if err := at.AddMessage(messageEntry.Text); err == nil {
			messageEntry.SetText("")
		} else {
			fmt.Println(err)
		}
	})

	// Set minimum size for the message list container
	messageListContainer := container.NewVScroll(ui.messageList)
	messageListContainer.SetMinSize(fyne.NewSize(400, 50)) // Set the height to 600

	// Create the message section as a bordered card with extra padding
	messageCard := widget.NewCard("Messages", "", container.NewVBox(
		messageListLabel,
		container.New(layout.NewBorderLayout(nil, nil, nil, addButton), messageEntry, addButton),
		messageListContainer,
	))

	paddedMessageCard := container.New(layout.NewPaddedLayout(), messageCard)

	return paddedMessageCard
}

func (at *AutoTyper) buildControlSection() fyne.CanvasObject {
	ui := at.ui
	// Create the control section as a bordered card with extra padding
	controlCard := widget.NewCard("Controls", "", container.NewVBox(
		container.NewHBox(ui.startButton, ui.stopButton, ui.clearButton),
	))

	paddedControlCard := container.New(layout.NewPaddedLayout(), controlCard)

	return paddedControlCard
}

func (at *AutoTyper) buildSettingsSection() fyne.CanvasObject {
	ui := at.ui
	// Interval entry and Set button
	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("Set interval (seconds)")

	setIntervalButton := widget.NewButtonWithIcon("Set Interval", theme.SettingsIcon(), func() {
		interval, err := time.ParseDuration(intervalEntry.Text + "s")
		if err == nil {
			at.SetInterval(int(interval.Seconds()))
			intervalEntry.SetText("")
		} else {
			fmt.Println("Invalid interval")
		}
	})

	// Create a checkbox for stopping on trade completion
	ui.stopOnTradeCheckbox = widget.NewCheck("Automatically Stop when Trade Completed", func(value bool) {
		// Handle checkbox change if needed
	})

	// Create a container for the interval controls (entry, button, and label)
	intervalContainer := container.NewVBox(
		container.New(layout.NewBorderLayout(nil, nil, nil, setIntervalButton), intervalEntry, setIntervalButton),
		ui.intervalLabel, // Place the interval label within the interval container
	)

	// Create the settings section as a bordered card with extra padding
	settingsCard := widget.NewCard("Settings", "", container.NewVBox(
		intervalContainer,      // Add the interval container with the label inside
		ui.stopOnTradeCheckbox, // Add the checkbox to the settings card
	))

	paddedSettingsCard := container.New(layout.NewPaddedLayout(), settingsCard)

	return paddedSettingsCard
}

func (at *AutoTyper) onTradeCompleted() {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	if at.ui.stopOnTradeCheckbox.Checked { // Check if the checkbox is checked
		if at.ticker != nil {
			at.Stop()
		}
	}
}

func (at *AutoTyper) updateMessageList() {
	at.ui.messageList.Refresh()
}

func (at *AutoTyper) updateUI() {
	ui := at.ui
	isRunning := at.ticker != nil

	// Update interval label
	ui.intervalLabel.SetText(fmt.Sprintf("Interval: %d seconds", at.intervalSec))

	// Update button states
	if isRunning {
		ui.startButton.Disable()
		ui.stopButton.Enable()
		ui.clearButton.Disable()
	} else {
		ui.startButton.Enable()
		ui.stopButton.Disable()
		ui.clearButton.Enable()
	}

	// Refresh the UI layout
	ui.window.Content().Refresh()
}
