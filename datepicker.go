package datepicker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	date     time.Time
	monthMap func(date time.Time) Month
	Keys     KeyMap
	Help     help.Model
	config   Config
	Colors   Colors
}

type (
	Month []Week
	Week  [7]int
)

func InitModel(config Config, colors Colors) *Model {
	m := &Model{
		date:   config.StartAt,
		Keys:   Keys,
		Help:   help.New(),
		config: config,
		Colors: colors,
	}
	m.monthMap = cachedMonthMaps(m.config.FirstWeekdayIsMo)
	return m
}

func (m *Model) CurrentValue() string {
	return m.date.Format(m.config.OutputFormat)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		case key.Matches(msg, m.Keys.Today):
			m.date = time.Now()
		case key.Matches(msg, m.Keys.Left):
			m.date = m.date.AddDate(0, 0, -1)
		case key.Matches(msg, m.Keys.Right):
			m.date = m.date.AddDate(0, 0, 1)
		case key.Matches(msg, m.Keys.Down):
			m.date = m.date.AddDate(0, 0, 7)
		case key.Matches(msg, m.Keys.Up):
			m.date = m.date.AddDate(0, 0, -7)
		case key.Matches(msg, m.Keys.WeekStart):
			m.date = time.Date(m.date.Year(), m.date.Month(), m.monthMap(m.date)[m.week()].firstDay(), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.Keys.WeekEnd):

			m.date = time.Date(m.date.Year(), m.date.Month(), m.monthMap(m.date)[m.week()].lastDay(), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.Keys.MonthStart):
			m.date = time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.Keys.MonthEnd):
			m.date = time.Date(m.date.Year(), m.date.Month(), daysInMonth(m.date.Year(), m.date.Month()), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.Keys.MonthPrev):
			m.date = m.date.AddDate(0, -1, 0)
		case key.Matches(msg, m.Keys.MonthNext):
			m.date = m.date.AddDate(0, 1, 0)
		case key.Matches(msg, m.Keys.YearPrev):
			m.date = m.date.AddDate(-1, 0, 0)
		case key.Matches(msg, m.Keys.YearNext):
			m.date = m.date.AddDate(1, 0, 0)
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var weekLegend string
	if m.config.FirstWeekdayIsMo {
		weekLegend = " Mo Tu We Th Fr Sa Su"
	} else {
		weekLegend = " Su Mo Tu We Th Fr Sa"
	}

	monthYearTitle := fmt.Sprintf("%s %d", m.date.Month(), m.date.Year())
	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		Render(
			strings.Repeat(" ", (len(weekLegend)-len(monthYearTitle))/2+1)+monthYearTitle+"\n"+weekLegend,
		) + "\n"

	notCurrentMonthdays := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("7"))

	for weekK, week := range m.monthMap(m.date) {
		for k, day := range week {
			if day == 0 {
				if weekK == 0 {
					prevMonthMap := m.monthMap(m.date.AddDate(0, -1, 0))
					s += notCurrentMonthdays.Render(fmt.Sprintf(" %2d", prevMonthMap[len(prevMonthMap)-1][k]))
				} else {
					nextMonthMap := m.monthMap(m.date.AddDate(0, +1, 0))
					s += notCurrentMonthdays.Render(fmt.Sprintf(" %2d", nextMonthMap[0][k]))
				}
				continue
			}

			today := day == time.Now().Day() && m.date.Month() == time.Now().Month() && m.date.Year() == time.Now().Year()
			weekend := k >= 5
			if !m.config.FirstWeekdayIsMo {
				weekend = k == 0 || k == 6
			}
			focused := day == m.date.Day()
			style := lipgloss.NewStyle()

			curDate := time.Date(m.date.Year(), m.date.Month(), day, m.date.Hour(), m.date.Minute(), m.date.Second(), m.date.Nanosecond(), m.date.Location()).Format("2006/01/02")

			if c, exists := m.Colors[curDate]; exists {
				if focused {
					style = style.Background(lipgloss.Color(c)).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color(c))
				}
			} else if today {
				if focused {
					style = style.Background(lipgloss.Color("9")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("9"))
				}
			} else if weekend {
				if focused {
					style = style.Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("4"))
				}
			} else {
				if focused {
					style = style.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("3"))
				}
			}
			s += style.Render(fmt.Sprintf(" %2d", day))
		}

		s += "\n"
	}

	if len(m.monthMap(m.date)) == 4 {
		s += "\n\n"
	} else if len(m.monthMap(m.date)) == 5 {
		s += "\n"
	}

	if !m.config.HideHelp {
		s += m.Help.View(m.Keys)
	}

	return s
}

func (week Week) firstDay() int {
	for _, day := range week {
		if day != 0 {
			return day
		}
	}
	return 0
}

func (week Week) lastDay() int {
	for i := 6; i >= 0; i-- {
		if week[i] != 0 {
			return week[i]
		}
	}
	return 0
}

func (m Model) week() int {
	firstDay := time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstDayN := firstDay.Weekday()
	if firstDayN == 0 {
		firstDayN = 7
	}
	offset := 1
	if m.config.FirstWeekdayIsMo {
		offset = 2
	}
	return (m.date.Day() + int(firstDayN) - offset) / 7
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func firstDayOfMonth(year int, month time.Month, firstWeekDayIsMonday bool) int {
	firstDay := int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())
	if firstWeekDayIsMonday {
		firstDay = (firstDay + 6) % 7
	}
	return firstDay
}

func cachedMonthMaps(firstWeekdayIsMo bool) func(date time.Time) Month {
	cache := make(map[string]Month)
	return func(date time.Time) Month {
		key := strconv.Itoa(date.Year() + int(date.Month()))
		if _, exists := cache[key]; !exists {
			daysInMonth := daysInMonth(date.Year(), date.Month())
			startDay := firstDayOfMonth(date.Year(), date.Month(), firstWeekdayIsMo)

			monthMap := make(Month, 0)
			week := Week{}
			dayCounter := 1

			// Fill the first week with leading zeros
			for i := range startDay {
				week[i] = 0
			}

			// Fill the days of the month
			for dayCounter <= daysInMonth {
				week[startDay] = dayCounter
				dayCounter++
				startDay++

				// If the week is full, add it to the weeks slice and reset
				if startDay == 7 {
					monthMap = append(monthMap, week)
					week = Week{}
					startDay = 0
				}
			}

			// Add the last week if it has any days
			if startDay > 0 {
				monthMap = append(monthMap, week)
			}
			cache[key] = monthMap
		}
		return cache[key]
	}
}
