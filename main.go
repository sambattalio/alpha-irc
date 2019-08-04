package main

import (
	"fmt"
	"github.com/sambattalio/alpha-irc/client"
	"github.com/awesome-gocui/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		fmt.Println(err)
	}
	defer g.Close()

	c := &client.Client{
		Gui: g,
	}

	u := &client.User{
		Server: "chat.freenode.net:6667",
		Nick: "student069client",
		User: "student069",
		Name: "sbattali",
	}

	if c.Connect(u) != nil {
		fmt.Println("Error initializing client")
	}

	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		fmt.Println(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, switchView); err != nil {
		fmt.Println(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, c.GetInput); err != nil {
		fmt.Println(err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		fmt.Println(err)
	}
	fmt.Println("quitting")
}

func switchView(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
        maxX, maxY := g.Size()
        if v, err := g.SetView("stream", maxX / 6 + 1, 0, maxX - 1, maxY - 4, 0); err != nil {
                if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
                v.Autoscroll = true

        }

	if v, err := g.SetView("channels", 0, 0, maxX / 6, (maxY - 4) / 2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
	}

	if v, err := g.SetView("users", 0, (maxY - 4) / 2, maxX / 6, maxY - 4, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
	}

	if v, err := g.SetView("input", 0, maxY - 3, maxX - 1, maxY - 1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }
		v.Wrap = true
		v.Editable = true
		if _, err := g.SetCurrentView("input"); err != nil {
                        return err
                }
	}

        return nil
}
