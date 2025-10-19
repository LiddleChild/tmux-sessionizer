package superlist

import (
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/colors"
	"github.com/LiddleChild/tmux-sessionpane/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

type previewInfo struct {
	// first line in bound
	TopBound int

	// last line in bound
	BottomBound int

	// global cursor offset
	CursorOffset int

	// rendering height
	Height int
}

func (m Model) renderFilter() string {
	return m.filter.View()
}

func (m Model) renderItem(item Item, style lipgloss.Style) string {
	switch item := item.(type) {
	case *filteredItem:
		var (
			builder      strings.Builder
			matchesIndex int
		)

		for i, r := range item.Label() {
			if matchesIndex < len(item.matches) && i == item.matches[matchesIndex] {
				builder.WriteString(style.
					Foreground(colors.Yellow).
					Render(string(r)))

				matchesIndex += 1
			} else {
				builder.WriteString(style.Render(string(r)))
			}
		}

		return builder.String()

	default:
		return item.Label()
	}
}

// TODO: find a way to refactor previewList() and render()
func (m Model) previewList() previewInfo {
	var (
		info previewInfo

		itemIdx          = 0
		currentRenderLen = 0
	)

	for _, g := range m.GetGroupIter() {
		// group name
		currentRenderLen += 1

		for range g.Items {
			if m.cursor == itemIdx {
				info.CursorOffset = currentRenderLen
			}

			// item
			currentRenderLen += 1
			itemIdx += 1
		}
	}

	info.Height = currentRenderLen
	info.TopBound = utils.Clamp(m.yOffset, 0, currentRenderLen)
	info.BottomBound = utils.Clamp(m.yOffset+m.listHeight-1, 0, currentRenderLen)

	return info
}

func (m Model) renderList() []string {
	var (
		idx   = 0
		lines = []string{}
	)

	for _, g := range m.filteredGroups {
		if len(g.Items) == 0 {
			continue
		}

		lines = append(lines, groupNameStyle.Render(g.Name))

		for _, i := range g.Items {
			var (
				style      lipgloss.Style
				isSelected = m.cursor == idx
			)
			if isSelected {
				style = hoveredItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			var itemName string
			if m.FocusedComponent() == FocusedComponentItem && isSelected {
				itemName = m.input.View()
			} else {
				itemName = m.renderItem(i, i.Style(style))
			}

			style = i.Style(style).
				BorderStyle(lipgloss.OuterHalfBlockBorder()).
				BorderBackground(style.GetBackground()).
				BorderForeground(style.GetForeground()).
				BorderLeft(isSelected)

			style = style.
				PaddingLeft(2 - style.GetHorizontalFrameSize())

			lines = append(lines,
				style.
					Width(m.width-style.GetHorizontalFrameSize()).
					Render(itemName),
			)

			idx += 1
		}
	}

	return lines
}
