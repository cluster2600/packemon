package monitor

import (
	"github.com/ddddddO/packemon"
	"github.com/ddddddO/packemon/internal/tui"
	"github.com/rivo/tview"
)

type ICMPv6 struct {
	*packemon.ICMPv6
}

func (*ICMPv6) rows() int {
	return 16
}

func (*ICMPv6) columns() int {
	return 30
}

func (i *ICMPv6) viewTable() *tview.Table {
	table := tview.NewTable().SetBorders(false)
	table.Box = tview.NewBox().SetBorder(true).SetTitle(" ICMPv6 Header ").SetTitleAlign(tview.AlignLeft).SetBorderPadding(1, 1, 1, 1)

	// Get the type name based on the ICMPv6 type value
	typeName := getICMPv6TypeName(i.Type)
	
	table.SetCell(0, 0, tui.TableCellTitle("Type"))
	table.SetCell(0, 1, tui.TableCellContent("%d (%s)", i.Type, typeName))

	table.SetCell(1, 0, tui.TableCellTitle("Code"))
	table.SetCell(1, 1, tui.TableCellContent("%d", i.Code))

	table.SetCell(2, 0, tui.TableCellTitle("Checksum"))
	table.SetCell(2, 1, tui.TableCellContent("0x%04x", i.Checksum))

	// For Echo Request/Reply, parse and display the additional fields
	if i.Type == packemon.ICMPv6_TYPE_ECHO_REQUEST || i.Type == packemon.ICMPv6_TYPE_ECHO_REPLY {
		// Extract Echo fields if we have enough data
		if len(i.MessageBody) >= 4 {
			identifier := uint16(i.MessageBody[0])<<8 | uint16(i.MessageBody[1])
			sequenceNumber := uint16(i.MessageBody[2])<<8 | uint16(i.MessageBody[3])
			
			table.SetCell(3, 0, tui.TableCellTitle("Identifier"))
			table.SetCell(3, 1, tui.TableCellContent("%d (0x%04x)", identifier, identifier))

			table.SetCell(4, 0, tui.TableCellTitle("Sequence Number"))
			table.SetCell(4, 1, tui.TableCellContent("%d (0x%04x)", sequenceNumber, sequenceNumber))
			
			// If there's data beyond the identifier and sequence number
			if len(i.MessageBody) > 4 {
				viewHexadecimalDump(table, 5, "Echo Data", i.MessageBody[4:])
			}
		} else {
			viewHexadecimalDump(table, 3, "Message Body", i.MessageBody)
		}
	} else if i.Type == packemon.ICMPv6_TYPE_ROUTER_ADVERTISEMENT {
		// Display Router Advertisement specific fields
		table.SetCell(3, 0, tui.TableCellTitle("Router Advertisement"))
		viewHexadecimalDump(table, 4, "Message Body", i.MessageBody)
	} else if i.Type == packemon.ICMPv6_TYPE_NEIGHBOR_SOLICITATION || i.Type == packemon.ICMPv6_TYPE_NEIGHBOR_ADVERTISEMENT {
		// Display Neighbor Discovery specific fields
		table.SetCell(3, 0, tui.TableCellTitle("Neighbor Discovery"))
		viewHexadecimalDump(table, 4, "Message Body", i.MessageBody)
	} else {
		// Generic display for other ICMPv6 message types
		viewHexadecimalDump(table, 3, "Message Body", i.MessageBody)
	}

	return table
}

// getICMPv6TypeName returns a human-readable name for ICMPv6 message types
func getICMPv6TypeName(typ uint8) string {
	switch typ {
	case packemon.ICMPv6_TYPE_DESTINATION_UNREACHABLE:
		return "Destination Unreachable"
	case packemon.ICMPv6_TYPE_PACKET_TOO_BIG:
		return "Packet Too Big"
	case packemon.ICMPv6_TYPE_TIME_EXCEEDED:
		return "Time Exceeded"
	case packemon.ICMPv6_TYPE_PARAMETER_PROBLEM:
		return "Parameter Problem"
	case packemon.ICMPv6_TYPE_ECHO_REQUEST:
		return "Echo Request"
	case packemon.ICMPv6_TYPE_ECHO_REPLY:
		return "Echo Reply"
	case packemon.ICMPv6_TYPE_ROUTER_SOLICITATION:
		return "Router Solicitation"
	case packemon.ICMPv6_TYPE_ROUTER_ADVERTISEMENT:
		return "Router Advertisement"
	case packemon.ICMPv6_TYPE_NEIGHBOR_SOLICITATION:
		return "Neighbor Solicitation"
	case packemon.ICMPv6_TYPE_NEIGHBOR_ADVERTISEMENT:
		return "Neighbor Advertisement"
	case packemon.ICMPv6_TYPE_REDIRECT:
		return "Redirect"
	default:
		return "Unknown"
	}
}
