package superlist

import (
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
			var style lipgloss.Style
			if m.cursor == idx {
				style = hoveredItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			var itemName string
			if m.Focused() && m.cursor == idx {
				itemName = m.input.View()
			} else {
				itemName = m.renderItem(i, i.Style(style))
			}

			lines = append(lines,
				i.Style(style).
					Width(m.width).
					Render(itemName),
			)

			idx += 1
		}
	}

	return lines
}
