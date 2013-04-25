package flip

import (
	"github.com/jteeuwen/ircb/cmd"
	"github.com/jteeuwen/ircb/plugin"
	"github.com/jteeuwen/ircb/proto"

	"strings"
)

const (
	table = "┻━┻"
)

var flipTable = map[rune]rune{
	'a':  'ɐ',
	'b':  'q',
	'c':  'ɔ',
	'd':  'p',
	'e':  'ǝ',
	'f':  'ɟ',
	'g':  'ƃ',
	'h':  'ɥ',
	'i':  'ı',
	'j':  'ɾ',
	'k':  'ʞ',
	'l':  'ʃ',
	'm':  'ɯ',
	'n':  'u',
	'r':  'ɹ',
	't':  'ʇ',
	'v':  'ʌ',
	'w':  'ʍ',
	'y':  'ʎ',
	'.':  '˙',
	'[':  ']',
	'(':  ')',
	'{':  '}',
	'?':  '¿',
	'!':  '¡',
	'\'': ',',
	'<':  '>',
	'_':  '‾',
	'&':  '⅋',
	';':  '؛',
	'"':  '„',
}

type Plugin struct {
	*plugin.Base
}

func init() {
	for k, v := range flipTable {
		flipTable[v] = k
	}
	plugin.Register(New)
}

func New(profile string) plugin.Plugin {
	return &Plugin{Base: plugin.New(profile, "flip")}
}

func (p *Plugin) Load(client *proto.Client) error {
	err := p.Base.Load(client)
	if err != nil {
		return err
	}
	cmnd := &cmd.Command{
		Name: "flip", Description: "Memetastic unicode table and text flipper",
		Execute: flipCmd,
	}
	cmd.Register(cmnd)

	return nil
}

func flipCmd(cmd *cmd.Command, client *proto.Client, msg *proto.Message) {
	targ := msg.Receiver
	if !msg.FromChannel() {
		targ = msg.SenderName
	}
	text := cmd.Data
	flipped := ""
	if len(text) > 0 {
		flipped = flip(text)
	} else {
		flipped = table
	}
	client.PrivMsg(targ, "(ノಠ益ಠ)ノ彡 "+flipped)
}

func flip(str string) string {
	out := ""
	for _, char := range strings.ToLower(str) {
		outChar := char
		if flipChar, ok := flipTable[char]; ok {
			outChar = flipChar
		}
		out = string(outChar) + out
	}
	return out
}
