package enums

//WeekDay
type WeekDay int32

//Types verify status of business
const (
	WeekDayUnknown WeekDay = iota
	WeekDaySaturday
	WeekDaySunday
	WeekDayMonday
	WeekDayTuesday
	WeekDayWednesday
	WeekDayThursday
	WeekDayFriday
)

//WeekDays list of Enum
var WeekDays = [...]string{
	"Unknown",
	"Saturday",
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
}

// String() function will return the english name
// that we want out constant Day be recognized as
func (weekDay WeekDay) String() string {
	return WeekDays[weekDay]
}
