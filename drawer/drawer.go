package drawer

import (
	"errors"
	"fmt"
	"github.com/inhies/go-bytesize"
	"log"
	"main/drawer_templates"
	"main/drawer_theme"
	"main/snapshot"
	"main/state_providers/battery_state"
	"main/state_providers/brightness_state"
	"main/state_providers/cpu_temp"
	"main/state_providers/keyboard_layout"
	"main/state_providers/network_connection_state"
	"main/state_providers/network_stat"
	"main/state_providers/notifications_state"
	"main/state_providers/volume_state"
	"main/state_providers/weather_state"
	"main/util"
	"time"
)

var batteryStatusMap = map[string]string{
	"Full":         drawer_templates.BatPartFull,
	"Discharging":  drawer_templates.BatPartDischarging,
	"Charging":     drawer_templates.BatPartCharging,
	"Not charging": drawer_templates.BatPartNotCharging,
	"Unknown":      drawer_templates.BatPartUnknown,
	"Empty":        "Empty",
	"Idle":         "Idle",
}

type Drawer struct {
	s       *snapshot.DwmBarStatsSnapshot
	_v      string
	t       *drawer_theme.Theme
	c       *util.Config
	checker *util.ErrorChecker
}

func (d *Drawer) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.redraw()
		}
	}
}

func (d *Drawer) redraw() {
	d._v = ""
	d.drawNetworkStat(d.s.NetworkStat)
	d.drawCpu(d.s.Cpu)
	d.drawTemp(d.s.CpuTemp)
	d.drawMemory(d.s.Mem)
	d.drawBrightness(d.s.BrightnessState)
	d.drawVolume(d.s.VolumeState)
	d.drawNetworkState(d.s.NetworkState)
	d.drawPowerState(d.s.BatteryState)
	d.drawKeyboardLayout(d.s.KeyboardLayout)
	d.drawNotificationsDisabled(d.s.NotificationsState)
	d.drawWeather(d.s.WeatherState)
	d.drawClock(d.s.NowDateTime)
	d.print()
}

func NewDwmBarDrawer(
	theme *drawer_theme.Theme,
	snapshot *snapshot.DwmBarStatsSnapshot,
	config *util.Config,
	checker *util.ErrorChecker,
) *Drawer {
	return &Drawer{t: theme, s: snapshot, c: config, checker: checker}
}

func (d *Drawer) add(string string) *Drawer {
	d._v = d._v + string
	return d
}

func (d *Drawer) drawWeather(stats weather_state.Stats) {
	if d.c.NoWeatherState {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.Weather,
		d.t.Black,
		d.t.Black,
		d.t.Green,
		drawer_templates.WeatherEmoji(stats.Code),
		stats.Temperature,
	)

	d.add(result)
}

func (d *Drawer) drawNetworkStat(stats network_stat.Stats) {
	if d.c.NoNetworkStats {
		return
	}

	resultP1 := fmt.Sprintf(
		drawer_templates.NetworkStatsDownload,
		d.t.Magenta,
		d.t.Magenta,
		d.t.Black,
		d.t.Black,
		d.t.Magenta,
		bytesize.New(float64(stats.RxBytes)).String(),
	)

	resultP2 := fmt.Sprintf(
		drawer_templates.NetworkStatsUpload,
		d.t.Magenta,
		d.t.Black,
		d.t.Black,
		d.t.Magenta,
		bytesize.New(float64(stats.TxBytes)).String(),
	)

	result := resultP1 + resultP2

	d.add(result)
}

func (d *Drawer) drawCpu(cpu int) {
	if d.c.NoCpu {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.CpuData,
		d.t.Cyan,
		d.t.Black,
		d.t.Cyan,
		cpu,
	) + "% "
	d.add(result)
}

func (d *Drawer) drawTemp(temp cpu_temp.Stats) {
	if d.c.NoTemp {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.CpuTemp,
		d.t.Red,
		d.t.Black,
		d.t.Red,
		temp.Temperature,
	)
	d.add(result)
}

func (d *Drawer) drawMemory(memory uint64) {
	if d.c.NoMemory {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.MemoryUsage,
		d.t.Magenta,
		d.t.Magenta,
		d.t.Black,
		d.t.Magenta,
		bytesize.New(float64(memory)).String(),
	)
	d.add(result)
}

func (d *Drawer) drawBrightness(stats brightness_state.Stats) {
	if d.c.NoBrightness {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.Brightness,
		d.t.Yellow,
		d.t.Yellow,
		d.t.Black,
		d.t.Yellow,
		stats.Brightness,
	)
	d.add(result)
}

func (d *Drawer) drawVolume(stats volume_state.Stats) {
	if d.c.NoVolume {
		return
	}

	if stats.Muted {
		result := fmt.Sprintf(
			drawer_templates.VolumeMuted,
			d.t.Magenta,
			d.t.Magenta,
			d.t.Black,
			d.t.Magenta,
		)
		d.add(result)
		return
	}

	var icon string
	if stats.Volume > 70 {
		icon = drawer_templates.VolumePartGT70
	} else if stats.Volume > 30 {
		icon = drawer_templates.VolumePartGT30
	} else if stats.Volume > 0 {
		icon = drawer_templates.VolumePartGT0
	} else {
		icon = drawer_templates.VolumePart0
	}

	result := fmt.Sprintf(
		drawer_templates.Volume,
		d.t.Magenta,
		d.t.Magenta,
		d.t.Black,
		d.t.Magenta,
		icon,
		stats.Volume,
	) + "% "
	d.add(result)
}

func (d *Drawer) drawNetworkState(stats network_connection_state.Stats) {
	if d.c.NoNetworkState {
		return
	}

	if stats.WiredConnected {
		result := fmt.Sprintf(
			drawer_templates.NetworkConnectionWired,
			d.t.Green,
			d.t.Green,
			d.t.Black,
			stats.WiredInterfaceName,
		)
		d.add(result)
		return
	}

	if stats.WirelessConnected {
		result := fmt.Sprintf(
			drawer_templates.NetworkConnectionWireless,
			d.t.Green,
			d.t.Black,
			stats.WirelessConnectionName,
		)
		d.add(result)
		return
	}
}

func (d *Drawer) drawPowerState(s battery_state.Stats) {
	if d.c.NoPowerState {
		return
	}

	status, exists := batteryStatusMap[s.State]
	if !exists {
		status = drawer_templates.BatPartUndefined
	}

	warn := ""
	if s.Percent <= 25 {
		warn = drawer_templates.BatPartWarningSymbol
	}

	statusTemplate := fmt.Sprintf(drawer_templates.BatPartStatus, d.t.Orange, d.t.Black, d.t.Orange, status)
	warningTemplate := fmt.Sprintf(drawer_templates.BatPartWarning, d.t.Black, d.t.Orange, warn)

	value := fmt.Sprintf(drawer_templates.Bat, statusTemplate, warningTemplate, s.Percent)

	d.add(value)
}

func (d *Drawer) drawKeyboardLayout(stats keyboard_layout.Stats) {
	if d.c.NoKeyboardLayout {
		return
	}

	result := fmt.Sprintf(
		drawer_templates.KeyboardLayout,
		d.t.Green,
		d.t.Green,
		d.t.Black,
		d.t.Green,
		stats.Lang,
	)
	d.add(result)
}

func (d *Drawer) drawClock(clockTime time.Time) {
	var (
		clockMonth   string
		clockWeekDay string
	)

	switch d.c.Lang {
	case "ru":
		clockMonth = drawer_templates.GetClockMonthRu(clockTime.Month())
		clockWeekDay = drawer_templates.GetClockWeekDayRu(clockTime.Weekday())
	case "en":
		clockMonth = drawer_templates.GetClockMonthEn(clockTime.Month())
		clockWeekDay = drawer_templates.GetClockWeekDayEn(clockTime.Weekday())
	default:
		str := "Invalid language(use default 'ru'): " + d.c.Lang
		d.checker.ErrorFound(errors.New(str))

		clockMonth = drawer_templates.GetClockMonthRu(clockTime.Month())
		clockWeekDay = drawer_templates.GetClockWeekDayRu(clockTime.Weekday())
	}

	date := fmt.Sprintf(
		drawer_templates.Clock,
		d.t.Blue,
		d.t.Black,
		d.t.Blue,
		fmt.Sprintf(
			"%s %d %s %02d:%02d",
			clockWeekDay,
			clockTime.Day(),
			clockMonth,
			clockTime.Hour(),
			clockTime.Minute(),
		),
	)
	d.add(date)
}

func (d *Drawer) blinkOneSecond() int64 {
	return time.Now().Unix() % 2 // 0 or 1
}

func (d *Drawer) drawNotificationsDisabled(state notifications_state.Stats) {
	if d.c.NoNotificationsState {
		return
	}

	if state.IsDisabled {
		var color string

		if d.c.EnableNotificationsStateBgBlinking && 0 == d.blinkOneSecond() {
			color = d.t.Orange
		} else {
			color = d.t.Cyan
		}

		date := fmt.Sprintf(
			drawer_templates.NotificationsDisabled,
			color,
			d.t.Black,
			color,
		)
		d.add(date)
	}
}

func (d *Drawer) print() {
	_, err := util.ExecCmd("xsetroot", "-name", d._v)
	if err != nil {
		log.Println(
			fmt.Sprintf(
				"Error in xsetroot -name: %s",
				err.Error(),
			),
		)
	}
}
