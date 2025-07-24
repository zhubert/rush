package interpreter

import (
	"testing"
	"time"
)

// Test Time namespace methods
func TestTimeNamespaceMethods(t *testing.T) {
	// Test Time.now()
	timeNamespace := &TimeNamespace{}
	nowResult := applyTimeNamespaceMethod(timeNamespace, "now")
	timeVal, ok := nowResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from Time.now(), got %T", nowResult)
	}
	if timeVal.Location != "Local" {
		t.Errorf("expected Time.now() location to be Local, got %s", timeVal.Location)
	}
	
	// Test Time.parse()
	parseResult := applyTimeNamespaceMethod(timeNamespace, "parse", &String{Value: "2024-01-15 14:30:00"})
	parsedTime, ok := parseResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from Time.parse(), got %T", parseResult)
	}
	
	// Test Time.new()
	newResult := applyTimeNamespaceMethod(timeNamespace, "new", 
		&Integer{Value: 2024}, &Integer{Value: 1}, &Integer{Value: 15},
		&Integer{Value: 14}, &Integer{Value: 30}, &Integer{Value: 0})
	newTime, ok := newResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from Time.new(), got %T", newResult)
	}
	if newTime.Location != "Local" {
		t.Errorf("expected Time.new() location to be Local, got %s", newTime.Location)
	}
}

// Test Duration namespace methods
func TestDurationNamespaceMethods(t *testing.T) {
	durationNamespace := &DurationNamespace{}
	
	// Test Duration.seconds()
	secResult := applyDurationNamespaceMethod(durationNamespace, "seconds", &Integer{Value: 30})
	secDur, ok := secResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from Duration.seconds(), got %T", secResult)
	}
	expected := int64(30 * int(time.Second))
	if secDur.Value != expected {
		t.Errorf("expected Duration.seconds(30) to be %d nanoseconds, got %d", expected, secDur.Value)
	}
	
	// Test Duration.minutes()
	minResult := applyDurationNamespaceMethod(durationNamespace, "minutes", &Integer{Value: 5})
	minDur, ok := minResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from Duration.minutes(), got %T", minResult)
	}
	expected = int64(5 * int(time.Minute))
	if minDur.Value != expected {
		t.Errorf("expected Duration.minutes(5) to be %d nanoseconds, got %d", expected, minDur.Value)
	}
	
	// Test Duration.hours()
	hourResult := applyDurationNamespaceMethod(durationNamespace, "hours", &Integer{Value: 2})
	hourDur, ok := hourResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from Duration.hours(), got %T", hourResult)
	}
	expected = int64(2 * int(time.Hour))
	if hourDur.Value != expected {
		t.Errorf("expected Duration.hours(2) to be %d nanoseconds, got %d", expected, hourDur.Value)
	}
	
	// Test Duration.days()
	dayResult := applyDurationNamespaceMethod(durationNamespace, "days", &Integer{Value: 1})
	dayDur, ok := dayResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from Duration.days(), got %T", dayResult)
	}
	expected = int64(24 * int(time.Hour))
	if dayDur.Value != expected {
		t.Errorf("expected Duration.days(1) to be %d nanoseconds, got %d", expected, dayDur.Value)
	}
	
	// Test Duration.parse()
	parseResult := applyDurationNamespaceMethod(durationNamespace, "parse", &String{Value: "2h30m15s"})
	parsedDur, ok := parseResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from Duration.parse(), got %T", parseResult)
	}
	expectedDuration, _ := time.ParseDuration("2h30m15s")
	if parsedDur.Value != int64(expectedDuration) {
		t.Errorf("expected Duration.parse('2h30m15s') to be %d nanoseconds, got %d", int64(expectedDuration), parsedDur.Value)
	}
}

// Test TimeZone namespace methods
func TestTimeZoneNamespaceMethods(t *testing.T) {
	tzNamespace := &TimeZoneNamespace{}
	
	// Test TimeZone.utc()
	utcResult := applyTimeZoneNamespaceMethod(tzNamespace, "utc")
	utcTz, ok := utcResult.(*TimeZone)
	if !ok {
		t.Fatalf("expected TimeZone result from TimeZone.utc(), got %T", utcResult)
	}
	if utcTz.Name != "UTC" {
		t.Errorf("expected TimeZone.utc() name to be UTC, got %s", utcTz.Name)
	}
	if utcTz.Offset != 0 {
		t.Errorf("expected TimeZone.utc() offset to be 0, got %d", utcTz.Offset)
	}
	
	// Test TimeZone.local()
	localResult := applyTimeZoneNamespaceMethod(tzNamespace, "local")
	localTz, ok := localResult.(*TimeZone)
	if !ok {
		t.Fatalf("expected TimeZone result from TimeZone.local(), got %T", localResult)
	}
	if localTz.Name != "Local" {
		t.Errorf("expected TimeZone.local() name to be Local, got %s", localTz.Name)
	}
}

// Test Time instance methods
func TestTimeInstanceMethods(t *testing.T) {
	// Create a test time: 2024-01-15 14:30:00
	testTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.Local)
	timeObj := &Time{
		Value:    testTime.UnixNano(),
		Location: "Local",
	}
	
	// Test format method
	timeMethod := &TimeMethod{Time: timeObj, Method: "format"}
	formatResult := applyTimeMethod(timeMethod, []Value{&String{Value: "2006-01-02 15:04:05"}}, nil)
	formatted, ok := formatResult.(*String)
	if !ok {
		t.Fatalf("expected String result from time.format(), got %T", formatResult)
	}
	expected := "2024-01-15 14:30:00"
	if formatted.Value != expected {
		t.Errorf("expected time.format() to return %s, got %s", expected, formatted.Value)
	}
	
	// Test year method
	timeMethod.Method = "year"
	yearResult := applyTimeMethod(timeMethod, []Value{}, nil)
	year, ok := yearResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.year(), got %T", yearResult)
	}
	if year.Value != 2024 {
		t.Errorf("expected time.year() to return 2024, got %d", year.Value)
	}
	
	// Test month method
	timeMethod.Method = "month"
	monthResult := applyTimeMethod(timeMethod, []Value{}, nil)
	month, ok := monthResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.month(), got %T", monthResult)
	}
	if month.Value != 1 {
		t.Errorf("expected time.month() to return 1, got %d", month.Value)
	}
	
	// Test day method
	timeMethod.Method = "day"
	dayResult := applyTimeMethod(timeMethod, []Value{}, nil)
	day, ok := dayResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.day(), got %T", dayResult)
	}
	if day.Value != 15 {
		t.Errorf("expected time.day() to return 15, got %d", day.Value)
	}
	
	// Test hour method
	timeMethod.Method = "hour"
	hourResult := applyTimeMethod(timeMethod, []Value{}, nil)
	hour, ok := hourResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.hour(), got %T", hourResult)
	}
	if hour.Value != 14 {
		t.Errorf("expected time.hour() to return 14, got %d", hour.Value)
	}
	
	// Test minute method
	timeMethod.Method = "minute"
	minuteResult := applyTimeMethod(timeMethod, []Value{}, nil)
	minute, ok := minuteResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.minute(), got %T", minuteResult)
	}
	if minute.Value != 30 {
		t.Errorf("expected time.minute() to return 30, got %d", minute.Value)
	}
	
	// Test second method
	timeMethod.Method = "second"
	secondResult := applyTimeMethod(timeMethod, []Value{}, nil)
	second, ok := secondResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from time.second(), got %T", secondResult)
	}
	if second.Value != 0 {
		t.Errorf("expected time.second() to return 0, got %d", second.Value)
	}
}

// Test Duration instance methods
func TestDurationInstanceMethods(t *testing.T) {
	// Create a test duration: 2h30m15s500ms
	testDuration := 2*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond
	durObj := &Duration{Value: int64(testDuration)}
	
	// Test total_seconds method
	durMethod := &DurationMethod{Duration: durObj, Method: "total_seconds"}
	totalSecsResult := applyDurationMethod(durMethod, []Value{}, nil)
	totalSecs, ok := totalSecsResult.(*Float)
	if !ok {
		t.Fatalf("expected Float result from duration.total_seconds(), got %T", totalSecsResult)
	}
	expectedSecs := testDuration.Seconds()
	if totalSecs.Value != expectedSecs {
		t.Errorf("expected duration.total_seconds() to return %f, got %f", expectedSecs, totalSecs.Value)
	}
	
	// Test total_minutes method
	durMethod.Method = "total_minutes"
	totalMinsResult := applyDurationMethod(durMethod, []Value{}, nil)
	totalMins, ok := totalMinsResult.(*Float)
	if !ok {
		t.Fatalf("expected Float result from duration.total_minutes(), got %T", totalMinsResult)
	}
	expectedMins := testDuration.Minutes()
	if totalMins.Value != expectedMins {
		t.Errorf("expected duration.total_minutes() to return %f, got %f", expectedMins, totalMins.Value)
	}
	
	// Test total_hours method
	durMethod.Method = "total_hours"
	totalHoursResult := applyDurationMethod(durMethod, []Value{}, nil)
	totalHours, ok := totalHoursResult.(*Float)
	if !ok {
		t.Fatalf("expected Float result from duration.total_hours(), got %T", totalHoursResult)
	}
	expectedHours := testDuration.Hours()
	if totalHours.Value != expectedHours {
		t.Errorf("expected duration.total_hours() to return %f, got %f", expectedHours, totalHours.Value)
	}
	
	// Test hours component method
	durMethod.Method = "hours"
	hoursResult := applyDurationMethod(durMethod, []Value{}, nil)
	hours, ok := hoursResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from duration.hours(), got %T", hoursResult)
	}
	if hours.Value != 2 {
		t.Errorf("expected duration.hours() to return 2, got %d", hours.Value)
	}
	
	// Test minutes component method
	durMethod.Method = "minutes"
	minutesResult := applyDurationMethod(durMethod, []Value{}, nil)
	minutes, ok := minutesResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from duration.minutes(), got %T", minutesResult)
	}
	if minutes.Value != 30 {
		t.Errorf("expected duration.minutes() to return 30, got %d", minutes.Value)
	}
	
	// Test seconds component method
	durMethod.Method = "seconds"
	secondsResult := applyDurationMethod(durMethod, []Value{}, nil)
	seconds, ok := secondsResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from duration.seconds(), got %T", secondsResult)
	}
	if seconds.Value != 15 {
		t.Errorf("expected duration.seconds() to return 15, got %d", seconds.Value)
	}
	
	// Test milliseconds component method
	durMethod.Method = "milliseconds"
	millisecondsResult := applyDurationMethod(durMethod, []Value{}, nil)
	milliseconds, ok := millisecondsResult.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from duration.milliseconds(), got %T", millisecondsResult)
	}
	if milliseconds.Value != 500 {
		t.Errorf("expected duration.milliseconds() to return 500, got %d", milliseconds.Value)
	}
}

// Test Duration arithmetic methods
func TestDurationArithmetic(t *testing.T) {
	dur1 := &Duration{Value: int64(2 * time.Hour)}
	dur2 := &Duration{Value: int64(30 * time.Minute)}
	
	// Test add method
	durMethod := &DurationMethod{Duration: dur1, Method: "add"}
	addResult := applyDurationMethod(durMethod, []Value{dur2}, nil)
	addedDur, ok := addResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from duration.add(), got %T", addResult)
	}
	expected := int64(2*time.Hour + 30*time.Minute)
	if addedDur.Value != expected {
		t.Errorf("expected duration.add() to return %d nanoseconds, got %d", expected, addedDur.Value)
	}
	
	// Test subtract method
	durMethod.Method = "subtract"
	subResult := applyDurationMethod(durMethod, []Value{dur2}, nil)
	subtractedDur, ok := subResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from duration.subtract(), got %T", subResult)
	}
	expected = int64(2*time.Hour - 30*time.Minute)
	if subtractedDur.Value != expected {
		t.Errorf("expected duration.subtract() to return %d nanoseconds, got %d", expected, subtractedDur.Value)
	}
	
	// Test multiply method
	durMethod.Method = "multiply"
	mulResult := applyDurationMethod(durMethod, []Value{&Integer{Value: 3}}, nil)
	multipliedDur, ok := mulResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from duration.multiply(), got %T", mulResult)
	}
	expected = int64(6 * time.Hour)
	if multipliedDur.Value != expected {
		t.Errorf("expected duration.multiply(3) to return %d nanoseconds, got %d", expected, multipliedDur.Value)
	}
	
	// Test divide method
	durMethod.Method = "divide"
	divResult := applyDurationMethod(durMethod, []Value{&Integer{Value: 2}}, nil)
	dividedDur, ok := divResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from duration.divide(), got %T", divResult)
	}
	expected = int64(1 * time.Hour)
	if dividedDur.Value != expected {
		t.Errorf("expected duration.divide(2) to return %d nanoseconds, got %d", expected, dividedDur.Value)
	}
}

// Test Duration validation methods  
func TestDurationValidation(t *testing.T) {
	positiveDur := &Duration{Value: int64(2 * time.Hour)}
	negativeDur := &Duration{Value: int64(-2 * time.Hour)}
	zeroDur := &Duration{Value: 0}
	
	// Test is_positive? method
	durMethod := &DurationMethod{Duration: positiveDur, Method: "is_positive?"}
	posResult := applyDurationMethod(durMethod, []Value{}, nil)
	isPositive, ok := posResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from duration.is_positive?(), got %T", posResult)
	}
	if !isPositive.Value {
		t.Errorf("expected positive duration.is_positive?() to return true")
	}
	
	// Test is_negative? method
	durMethod = &DurationMethod{Duration: negativeDur, Method: "is_negative?"}
	negResult := applyDurationMethod(durMethod, []Value{}, nil)
	isNegative, ok := negResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from duration.is_negative?(), got %T", negResult)
	}
	if !isNegative.Value {
		t.Errorf("expected negative duration.is_negative?() to return true")
	}
	
	// Test is_zero? method
	durMethod = &DurationMethod{Duration: zeroDur, Method: "is_zero?"}
	zeroResult := applyDurationMethod(durMethod, []Value{}, nil)
	isZero, ok := zeroResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from duration.is_zero?(), got %T", zeroResult)
	}
	if !isZero.Value {
		t.Errorf("expected zero duration.is_zero?() to return true")
	}
	
	// Test abs method on negative duration
	durMethod = &DurationMethod{Duration: negativeDur, Method: "abs"}
	absResult := applyDurationMethod(durMethod, []Value{}, nil)
	absDur, ok := absResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from duration.abs(), got %T", absResult)
	}
	if absDur.Value != int64(2*time.Hour) {
		t.Errorf("expected abs of negative duration to be positive, got %d", absDur.Value)
	}
}

// Test Time arithmetic and comparison methods
func TestTimeArithmetic(t *testing.T) {
	// Create two test times
	baseTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.Local)
	timeObj1 := &Time{Value: baseTime.UnixNano(), Location: "Local"}
	
	laterTime := baseTime.Add(2 * time.Hour)
	timeObj2 := &Time{Value: laterTime.UnixNano(), Location: "Local"}
	
	duration := &Duration{Value: int64(30 * time.Minute)}
	
	// Test add_duration method
	timeMethod := &TimeMethod{Time: timeObj1, Method: "add_duration"}
	addResult := applyTimeMethod(timeMethod, []Value{duration}, nil)
	newTime, ok := addResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from time.add_duration(), got %T", addResult)
	}
	expected := baseTime.Add(30 * time.Minute)
	if newTime.Value != expected.UnixNano() {
		t.Errorf("expected time.add_duration() to return %d, got %d", expected.UnixNano(), newTime.Value)
	}
	
	// Test subtract_duration method
	timeMethod.Method = "subtract_duration"
	subResult := applyTimeMethod(timeMethod, []Value{duration}, nil)
	earlierTime, ok := subResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from time.subtract_duration(), got %T", subResult)
	}
	expected = baseTime.Add(-30 * time.Minute)
	if earlierTime.Value != expected.UnixNano() {
		t.Errorf("expected time.subtract_duration() to return %d, got %d", expected.UnixNano(), earlierTime.Value)
	}
	
	// Test difference method
	timeMethod.Method = "difference"
	diffResult := applyTimeMethod(timeMethod, []Value{timeObj2}, nil)
	diff, ok := diffResult.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from time.difference(), got %T", diffResult)
	}
	expectedDiff := int64(-2 * time.Hour) // timeObj1 is 2 hours before timeObj2
	if diff.Value != expectedDiff {
		t.Errorf("expected time.difference() to return %d, got %d", expectedDiff, diff.Value)
	}
	
	// Test is_before? method
	timeMethod.Method = "is_before?"
	beforeResult := applyTimeMethod(timeMethod, []Value{timeObj2}, nil)
	isBefore, ok := beforeResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from time.is_before?(), got %T", beforeResult)
	}
	if !isBefore.Value {
		t.Errorf("expected time.is_before?() to return true for earlier time")
	}
	
	// Test is_after? method
	timeMethod.Method = "is_after?"
	afterResult := applyTimeMethod(timeMethod, []Value{timeObj2}, nil)
	isAfter, ok := afterResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from time.is_after?(), got %T", afterResult)
	}
	if isAfter.Value {
		t.Errorf("expected time.is_after?() to return false for earlier time")
	}
	
	// Test is_equal? method
	timeMethod.Method = "is_equal?"
	equalResult := applyTimeMethod(timeMethod, []Value{timeObj1}, nil)
	isEqual, ok := equalResult.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from time.is_equal?(), got %T", equalResult)
	}
	if !isEqual.Value {
		t.Errorf("expected time.is_equal?() to return true for same time")
	}
}

// Test timezone conversion methods
func TestTimeZoneConversion(t *testing.T) {
	// Create a test time in local timezone
	baseTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.Local)
	timeObj := &Time{Value: baseTime.UnixNano(), Location: "Local"}
	
	// Test to_utc method
	timeMethod := &TimeMethod{Time: timeObj, Method: "to_utc"}
	utcResult := applyTimeMethod(timeMethod, []Value{}, nil)
	utcTime, ok := utcResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from time.to_utc(), got %T", utcResult)
	}
	if utcTime.Location != "UTC" {
		t.Errorf("expected time.to_utc() location to be UTC, got %s", utcTime.Location)
	}
	
	// Test to_local method
	timeMethod.Method = "to_local"
	localResult := applyTimeMethod(timeMethod, []Value{}, nil)
	localTime, ok := localResult.(*Time)
	if !ok {
		t.Fatalf("expected Time result from time.to_local(), got %T", localResult)
	}
	if localTime.Location != "Local" {
		t.Errorf("expected time.to_local() location to be Local, got %s", localTime.Location)
	}
}

// Test TimeZone methods
func TestTimeZoneMethods(t *testing.T) {
	utcTz := &TimeZone{Name: "UTC", Offset: 0}
	
	// Test abbreviation method
	tzMethod := &TimeZoneMethod{TimeZone: utcTz, Method: "abbreviation"}
	abbrevResult := applyTimeZoneMethod(tzMethod, []Value{}, nil)
	abbrev, ok := abbrevResult.(*String)
	if !ok {
		t.Fatalf("expected String result from timezone.abbreviation(), got %T", abbrevResult)
	}
	// UTC abbreviation should be "UTC"
	if abbrev.Value != "UTC" {
		t.Errorf("expected UTC timezone abbreviation to be 'UTC', got '%s'", abbrev.Value)
	}
}

// Test error cases
func TestTimeModuleErrors(t *testing.T) {
	timeNamespace := &TimeNamespace{}
	
	// Test wrong number of arguments for Time.now()
	nowResult := applyTimeNamespaceMethod(timeNamespace, "now", &String{Value: "extra"})
	if !isError(nowResult) {
		t.Errorf("expected error for Time.now() with arguments, got %T", nowResult)
	}
	
	// Test invalid time string for Time.parse()
	parseResult := applyTimeNamespaceMethod(timeNamespace, "parse", &String{Value: "invalid time"})
	if !isError(parseResult) {
		t.Errorf("expected error for Time.parse() with invalid string, got %T", parseResult)
	}
	
	// Test wrong argument type for Time.parse()
	parseResult = applyTimeNamespaceMethod(timeNamespace, "parse", &Integer{Value: 42})
	if !isError(parseResult) {
		t.Errorf("expected error for Time.parse() with non-string argument, got %T", parseResult)
	}
	
	durationNamespace := &DurationNamespace{}
	
	// Test wrong argument type for Duration.seconds()
	secResult := applyDurationNamespaceMethod(durationNamespace, "seconds", &String{Value: "not a number"})
	if !isError(secResult) {
		t.Errorf("expected error for Duration.seconds() with string argument, got %T", secResult)
	}
	
	// Test invalid duration string for Duration.parse()
	parseResult = applyDurationNamespaceMethod(durationNamespace, "parse", &String{Value: "invalid duration"})
	if !isError(parseResult) {
		t.Errorf("expected error for Duration.parse() with invalid string, got %T", parseResult)
	}
}

// Test value type information
func TestTimeValueTypes(t *testing.T) {
	timeObj := &Time{Value: time.Now().UnixNano(), Location: "Local"}
	durObj := &Duration{Value: int64(2 * time.Hour)}
	tzObj := &TimeZone{Name: "UTC", Offset: 0}
	
	// Test Time value type
	if timeObj.Type() != TIME_VALUE {
		t.Errorf("expected Time.Type() to return TIME_VALUE, got %s", timeObj.Type())
	}
	
	// Test Duration value type
	if durObj.Type() != DURATION_VALUE {
		t.Errorf("expected Duration.Type() to return DURATION_VALUE, got %s", durObj.Type())
	}
	
	// Test TimeZone value type
	if tzObj.Type() != TIMEZONE_VALUE {
		t.Errorf("expected TimeZone.Type() to return TIMEZONE_VALUE, got %s", tzObj.Type())
	}
	
	// Test namespace types
	timeNamespace := &TimeNamespace{}
	if timeNamespace.Type() != TIME_NAMESPACE_VALUE {
		t.Errorf("expected TimeNamespace.Type() to return TIME_NAMESPACE_VALUE, got %s", timeNamespace.Type())
	}
	
	durationNamespace := &DurationNamespace{}
	if durationNamespace.Type() != DURATION_NAMESPACE_VALUE {
		t.Errorf("expected DurationNamespace.Type() to return DURATION_NAMESPACE_VALUE, got %s", durationNamespace.Type())
	}
	
	tzNamespace := &TimeZoneNamespace{}
	if tzNamespace.Type() != TIMEZONE_NAMESPACE_VALUE {
		t.Errorf("expected TimeZoneNamespace.Type() to return TIMEZONE_NAMESPACE_VALUE, got %s", tzNamespace.Type())
	}
	
	// Test method types
	timeMethod := &TimeMethod{Time: timeObj, Method: "format"}
	if timeMethod.Type() != TIME_METHOD_VALUE {
		t.Errorf("expected TimeMethod.Type() to return TIME_METHOD_VALUE, got %s", timeMethod.Type())
	}
	
	durMethod := &DurationMethod{Duration: durObj, Method: "total_seconds"}
	if durMethod.Type() != DURATION_METHOD_VALUE {
		t.Errorf("expected DurationMethod.Type() to return DURATION_METHOD_VALUE, got %s", durMethod.Type())
	}
	
	tzMethod := &TimeZoneMethod{TimeZone: tzObj, Method: "abbreviation"}
	if tzMethod.Type() != TIMEZONE_METHOD_VALUE {
		t.Errorf("expected TimeZoneMethod.Type() to return TIMEZONE_METHOD_VALUE, got %s", tzMethod.Type())
	}
}

// Test Inspect methods for debugging output
func TestTimeInspectMethods(t *testing.T) {
	timeObj := &Time{Value: time.Date(2024, 1, 15, 14, 30, 0, 0, time.Local).UnixNano(), Location: "Local"}
	durObj := &Duration{Value: int64(2 * time.Hour)}
	tzObj := &TimeZone{Name: "UTC", Offset: 0}
	
	// Test Time inspect
	timeInspect := timeObj.Inspect()
	if timeInspect == "" {
		t.Errorf("expected Time.Inspect() to return non-empty string")
	}
	
	// Test Duration inspect
	durInspect := durObj.Inspect()
	if durInspect == "" {
		t.Errorf("expected Duration.Inspect() to return non-empty string")
	}
	
	// Test TimeZone inspect
	tzInspect := tzObj.Inspect()
	if tzInspect == "" {
		t.Errorf("expected TimeZone.Inspect() to return non-empty string")
	}
}

// Helper function to check if a value is an error
func isError(val Value) bool {
	if val == nil {
		return false
	}
	_, ok := val.(*Error)
	return ok
}