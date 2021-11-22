package paack

import "time"

type ScheduleSlot struct {
	Start Schedule `json:"start"`
	End   Schedule `json:"end"`
}

func NewScheduleSlot(start time.Time, end time.Time) ScheduleSlot {
	return ScheduleSlot{
		Start: NewSchedule(start),
		End:   NewSchedule(end),
	}
}

type Schedule struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

func NewSchedule(time time.Time) Schedule {
	return Schedule{
		Date: time.Format("2006-01-02"),
		Time: time.Format("15:04:05"),
	}
}
