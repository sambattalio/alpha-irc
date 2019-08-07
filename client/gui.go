package client

import (
	"github.com/awesome-gocui/gocui"
)

func (c *Client) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	for channel := range c.channels {
		if v, err := g.SetView(channel, maxX / 6 + 1, 0, maxX - 1, maxY - 4, 0); err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}

			v.Wrap = true
			v.Autoscroll = true

		}
	}

	if v, err := g.SetView("channels", 0, 0, maxX / 6, (maxY - 4) / 2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
                        return err
                }

                v.Wrap = true
	}

	if v, err := g.SetView("users", 0, (maxY - 4) / 2 + 1, maxX / 6, maxY - 4, 0); err != nil {
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

