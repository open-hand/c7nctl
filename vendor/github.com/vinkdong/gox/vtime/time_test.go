package vtime

import (
	"testing"
	"time"
	"fmt"
)

const (
	TestTimeMS     = 1514461867000
	TestTimeNS     = 1514461867000000000
	TestTimeS      = 1514461867
	TestTimeYear   = 2017
	TestTimeSecond = 7
	TestTimeMinute = 51
	TestTimeString = "2017-12-28 19:51:07"
	TestTimeFormat = "2006-01-02 15:04:05"
	TestTimeTZ     = "Asia/Shanghai"
)

func checkTimeIn(t *testing.T, time time.Time) {
	if time.UTC().Year() != TestTimeYear {
		t.Errorf("Expect Time year is %d, But got %d", TestTimeYear, time.UTC().Year())
	}
	if time.Second() != TestTimeSecond {
		t.Errorf("Expect Time second is %d, But got %d", TestTimeSecond, time.UTC().Second())
	}
	if time.Minute() != TestTimeMinute {
		t.Errorf("Expect Time minute is %d, But got %d", TestTimeMinute, time.UTC().Minute())
	}
	if time.UnixNano() != TestTimeNS {
		t.Errorf("Expect Timestamp is %d, but got %d", TestTimeNS, time.UnixNano())
	}
}

func TestParserTimestampMs(t *testing.T) {
	var tsMs int64
	tsMs = TestTimeMS
	time := ParserTimestampMs(tsMs)
	checkTimeIn(t, time)
}

func TestParserTimestampNs(t *testing.T) {
	var tsNs int64
	tsNs = TestTimeNS
	time := ParserTimestampNs(tsNs)
	checkTimeIn(t, time)
}

func TestParserTimestampS(t *testing.T) {
	var tsS int64
	tsS = TestTimeS
	time := ParserTimestampS(tsS)
	checkTimeIn(t, time)
}

func TestParserVTime(t *testing.T) {
	vt := Time{
		Format: "2006-01-02 15:04:05",
		Value:  "2017-12-28 19:51:07",
		TZ: TestTimeTZ,
	}
	time, err := vt.Parser()
	if err!=nil{
		t.Fail()
	}
	checkTimeIn(t,time)
}

func TestFromTime(t *testing.T) {
	vt := Time{
		Format: TestTimeFormat,
		TZ: TestTimeTZ,
	}
	var tsS int64
	tsS = TestTimeS
	time := ParserTimestampS(tsS)
	vt.FromTime(time)
	if vt.Value != TestTimeString {
		t.Errorf("expect time value is %s but got %s", TestTimeString, vt.Value)
	}
}

func TestTimeTransfer(t *testing.T) {
	to := &Time{
		Format: TestTimeFormat,
		TZ:     TestTimeTZ,
	}
	from := &Time{
		Format: "timestamp",
		Value:  fmt.Sprintf("%d", TestTimeMS),
		Unit:   "ms",
	}
	from.Transfer(to)
	if to.Value != TestTimeString {
		t.Errorf("expect time value is %s but got %s", TestTimeString, to.Value)
	}
}

func TestTime_FromRelativeTime(t *testing.T) {
	tm := &Time{
		Format: TestTimeFormat,
		TZ:     TestTimeTZ,
	}
	tm.FromRelativeTime("now+1h")
	if tm.Time.Hour() != time.Now().UTC().Hour()+1 {
		t.Errorf("expect time hour is %d but got %d", time.Now().UTC().Hour()+1, tm.Time.Hour(), )
	}
}