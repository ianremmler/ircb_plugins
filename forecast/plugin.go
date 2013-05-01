package forecast

import (
	"github.com/jteeuwen/ircb/cmd"
	"github.com/jteeuwen/ircb/plugin"
	"github.com/jteeuwen/ircb/proto"

	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	locUrlFmt = "http://nominatim.openstreetmap.org/search?format=json&limit=1&q=%s"
)

type Plugin struct {
	*plugin.Base
}

type location struct {
	Name string `json:"display_name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

func init() { plugin.Register(New) }

func New(profile string) plugin.Plugin {
	return &Plugin{Base: plugin.New(profile, "forecast")}
}

func (p *Plugin) Load(client *proto.Client) error {
	err := p.Base.Load(client)
	if err != nil {
		return err
	}
	cmnd := &cmd.Command{
		Name:        "forecast",
		Description: "Print a pretty forecast for the given location",
		Execute:     forecastCmd,
	}
	cmnd.Params = []cmd.Param{
		{Name: "loc", Description: "Location to get forecast for", Pattern: cmd.RegAny},
	}
	cmd.Register(cmnd)

	return nil
}

func forecastCmd(cmd *cmd.Command, client *proto.Client, msg *proto.Message) {
	loc := regexp.MustCompile("\\s+").ReplaceAllLiteralString(cmd.Data, "+")
	targ := msg.Receiver
	if !msg.FromChannel() {
		targ = msg.SenderName
	}
	const locErr = "I had a problem finding that location."
	const fcErr = "I had a problem getting a forecast for that location."

	resp, err := http.Get(fmt.Sprintf(locUrlFmt, loc))
	if err != nil {
		client.PrivMsg(targ, locErr)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		client.PrivMsg(targ, locErr)
		return
	}
	dec := json.NewDecoder(resp.Body)
	locs := []location{}
	dec.Decode(&locs)
	if len(locs) == 0 {
		client.PrivMsg(targ, locErr)
		return
	}
	fc, err := forecast(locs[0])
	if err != nil {
		client.PrivMsg(targ, fcErr)
		return
	}
	for _, line := range strings.Split(fc, "\n") {
		client.PrivMsg(targ, line)
	}
}
