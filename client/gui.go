package client

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (c *Client) Layout(_ *gocui.Gui) error {
	maxX, maxY := c.gui.Size()
	for channel := range c.channels {
		if v, err := c.gui.SetView(channel, maxX / 6 + 1, 0, maxX - 1, maxY - 4, 0); err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}

			v.Wrap = true
			v.Autoscroll = true

		}
	}

	if v, err := c.gui.SetView("channels", 0, 0, maxX / 6, (maxY - 4) / 2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
	}

	if v, err := c.gui.SetView("users", 0, (maxY - 4) / 2 + 1, maxX / 6, maxY - 4, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
	}

	if v, err := c.gui.SetView("input", 0, maxY - 3, maxX - 1, maxY - 1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }
		v.Wrap = true
		v.Editable = true
		if _, err := c.gui.SetCurrentView("input"); err != nil {
                        return err
                }
	}

	return nil
}

func (c *Client) setKeybindings() {
	if err := c.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
                fmt.Println(err)
        }

        if err := c.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, switchView); err != nil {
                fmt.Println(err)
        }

        if err := c.gui.SetKeybinding("channels", gocui.KeyEnter, gocui.ModNone, c.setChannelView); err != nil {
                fmt.Println(err)
        }

        if err := c.gui.SetKeybinding("channels", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
                fmt.Println(err)
        }

        if err := c.gui.SetKeybinding("channels", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
                fmt.Println(err)
        }

	if err := c.gui.SetKeybinding("channels", gocui.KeyBackspace, gocui.ModNone, c.deleteChannel); err != nil {
		fmt.Println(err)
	}

	if err := c.gui.SetKeybinding("channels", gocui.KeyBackspace2, gocui.ModNone, c.deleteChannel); err != nil {
		fmt.Println(err)
	}

        if err := c.gui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, c.GetInput); err != nil {
                fmt.Println(err)
        }
}

func (c *Client) startGui() error {
        if err := c.gui.MainLoop(); err != nil && !gocui.IsQuit(err) {
                return err
        }
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
        cx, cy := v.Cursor()

        if channel, _ := v.Line(cy + 1); channel == "" {
                return nil
        }

        if err := v.SetCursor(cx, cy + 1); err != nil {
                return err
        }
        return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
        cx, cy := v.Cursor();
        if cy == 0 {
                return nil
        }
        if err := v.SetCursor(cx, cy - 1); err != nil {
                return err
        }
        return nil
}

func switchView(g *gocui.Gui, v *gocui.View) error {
        g.SetCurrentView("channels")
        return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()
        return gocui.ErrQuit
}

