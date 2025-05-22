package generator

import (
	"context"
	"encoding/binary"

	"github.com/ddddddO/packemon"
	"github.com/rivo/tview"
)

func (g *generator) icmpv6Form() *tview.Form {
	icmpv6Form := tview.NewForm().
		AddTextView("ICMPv6", "This section generates the ICMPv6 packet.\nSupports Echo Request and other ICMPv6 message types.", 60, 4, true, false).
		AddDropDown("Type", []string{
			"Echo Request (128)",
			"Router Solicitation (133)",
			"Neighbor Solicitation (135)",
		}, 0, func(option string, optionIndex int) {
			switch optionIndex {
			case 0:
				g.sender.packets.icmpv6.Type = packemon.ICMPv6_TYPE_ECHO_REQUEST
			case 1:
				g.sender.packets.icmpv6.Type = packemon.ICMPv6_TYPE_ROUTER_SOLICITATION
			case 2:
				g.sender.packets.icmpv6.Type = packemon.ICMPv6_TYPE_NEIGHBOR_SOLICITATION
			}
		}).
		AddInputField("Code(hex)", DEFAULT_ICMPv6_CODE, 4, func(textToCheck string, lastChar rune) bool {
			if len(textToCheck) < 4 {
				return true
			} else if len(textToCheck) > 4 {
				return false
			}

			b, err := strHexToUint8(textToCheck)
			if err != nil {
				return false
			}
			g.sender.packets.icmpv6.Code = uint8(b)

			return true
		}, nil).
		AddInputField("Identifier(hex)", DEFAULT_ICMPv6_IDENTIFIER, 6, func(textToCheck string, lastChar rune) bool {
			if len(textToCheck) < 6 {
				return true
			} else if len(textToCheck) > 6 {
				return false
			}

			b, err := packemon.StrHexToBytes2(textToCheck)
			if err != nil {
				return false
			}
			if g.sender.packets.icmpv6Echo == nil {
				g.sender.packets.icmpv6Echo = &packemon.ICMPv6Echo{}
			}
			g.sender.packets.icmpv6Echo.Identifier = binary.BigEndian.Uint16(b)

			return true
		}, nil).
		AddInputField("Sequence Number(hex)", DEFAULT_ICMPv6_SEQUENCE, 6, func(textToCheck string, lastChar rune) bool {
			if len(textToCheck) < 6 {
				return true
			} else if len(textToCheck) > 6 {
				return false
			}

			b, err := packemon.StrHexToBytes2(textToCheck)
			if err != nil {
				return false
			}
			if g.sender.packets.icmpv6Echo == nil {
				g.sender.packets.icmpv6Echo = &packemon.ICMPv6Echo{}
			}
			g.sender.packets.icmpv6Echo.SequenceNumber = binary.BigEndian.Uint16(b)

			return true
		}, nil).
		AddButton("Send!", func() {
			if err := g.sender.sendLayer4IPv6(context.TODO()); err != nil {
				g.addErrPage(err)
			}
		}).
		AddButton("Quit", func() {
			g.app.Stop()
		})

	return icmpv6Form
}
