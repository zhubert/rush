package interpreter

import (
	"testing"
	"time"
)

// Test Time namespace builtin functionality
func TestTimeNamespaceBuiltin(t *testing.T) {
	// Test Time namespace builtin exists
	timeNamespaceFunc := builtins["Time"]
	if timeNamespaceFunc == nil {
		t.Fatal("Time builtin function not found")
	}

	// Get the Time namespace
	timeNamespace := timeNamespaceFunc.Fn()
	if timeNamespace == nil {
		t.Fatal("Time namespace returned nil")
	}
	
	_, ok := timeNamespace.(*TimeNamespace)
	if !ok {
		t.Fatalf("expected TimeNamespace, got %T", timeNamespace)
	}
}

// Test Duration namespace builtin functionality
func TestDurationNamespaceBuiltin(t *testing.T) {
	// Test Duration namespace builtin exists
	durationNamespaceFunc := builtins["Duration"]
	if durationNamespaceFunc == nil {
		t.Fatal("Duration builtin function not found")
	}

	// Get the Duration namespace
	durationNamespace := durationNamespaceFunc.Fn()
	if durationNamespace == nil {
		t.Fatal("Duration namespace returned nil")
	}
	
	_, ok := durationNamespace.(*DurationNamespace)
	if !ok {
		t.Fatalf("expected DurationNamespace, got %T", durationNamespace)
	}
}

// Test TimeZone namespace builtin functionality
func TestTimeZoneNamespaceBuiltin(t *testing.T) {
	// Test TimeZone namespace builtin exists
	timezoneNamespaceFunc := builtins["TimeZone"]
	if timezoneNamespaceFunc == nil {
		t.Fatal("TimeZone builtin function not found")
	}

	// Get the TimeZone namespace
	timezoneNamespace := timezoneNamespaceFunc.Fn()
	if timezoneNamespace == nil {
		t.Fatal("TimeZone namespace returned nil")
	}
	
	_, ok := timezoneNamespace.(*TimeZoneNamespace)
	if !ok {
		t.Fatalf("expected TimeZoneNamespace, got %T", timezoneNamespace)
	}
}

// Test Time.now() functionality
func TestTimeNow(t *testing.T) {
	timeNamespace := &TimeNamespace{}
	
	// Test Time.now()
	result := applyTimeNamespaceMethod(timeNamespace, "now")
	timeObj, ok := result.(*Time)
	if !ok {
		t.Fatalf("expected Time object, got %T", result)
	}
	
	// Verify the time is recent (within last minute)
	now := time.Now()
	if now.Sub(timeObj.Value) > time.Minute {
		t.Fatal("Time.now() returned time too far in the past")
	}
}

// Test Time.parse() functionality
func TestTimeParse(t *testing.T) {
	timeNamespace := &TimeNamespace{}
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"2024-01-15 14:30:00", true},
		{"2024-01-15T14:30:00Z", true},
		{"2024-01-15", true},
		{"14:30:00", true},
		{"invalid time", false},
	}
	
	for _, test := range tests {
		result := applyTimeNamespaceMethod(timeNamespace, "parse", &String{Value: test.input})
		
		if test.expected {
			timeObj, ok := result.(*Time)
			if !ok {
				t.Fatalf("expected Time object for input %s, got %T", test.input, result)
			}
			if timeObj.Value.IsZero() {
				t.Fatalf("parsed time should not be zero for input %s", test.input)
			}
		} else {
			if !isError(result) {
				t.Fatalf("expected error for invalid input %s, got %T", test.input, result)
			}
		}
	}
}

// Test Time.new() functionality
func TestTimeNew(t *testing.T) {
	timeNamespace := &TimeNamespace{}
	
	// Test with year, month, day only
	result := applyTimeNamespaceMethod(timeNamespace, "new", 
		&Integer{Value: 2024}, &Integer{Value: 1}, &Integer{Value: 15})
	timeObj, ok := result.(*Time)
	if !ok {
		t.Fatalf("expected Time object, got %T", result)
	}
	
	if timeObj.Value.Year() != 2024 || timeObj.Value.Month() != 1 || timeObj.Value.Day() != 15 {
		t.Fatal("Time.new() did not create correct date")
	}
	
	// Test with full date and time
	result = applyTimeNamespaceMethod(timeNamespace, "new",
		&Integer{Value: 2024}, &Integer{Value: 1}, &Integer{Value: 15},
		&Integer{Value: 14}, &Integer{Value: 30}, &Integer{Value: 45})
	timeObj, ok = result.(*Time)
	if !ok {
		t.Fatalf("expected Time object, got %T", result)
	}
	
	if timeObj.Value.Hour() != 14 || timeObj.Value.Minute() != 30 || timeObj.Value.Second() != 45 {
		t.Fatal("Time.new() did not create correct time")
	}
}

// Test Duration creation methods
func TestDurationCreation(t *testing.T) {
	durationNamespace := &DurationNamespace{}
	
	tests := []struct {
		method   string
		value    Value
		expected time.Duration
	}{
		{"seconds", &Integer{Value: 30}, 30 * time.Second},
		{"minutes", &Integer{Value: 5}, 5 * time.Minute},
		{"hours", &Integer{Value: 2}, 2 * time.Hour},
		{"days", &Integer{Value: 1}, 24 * time.Hour},
		{"seconds", &Float{Value: 1.5}, time.Duration(1.5 * float64(time.Second))},
	}
	
	for _, test := range tests {
		result := applyDurationNamespaceMethod(durationNamespace, test.method, test.value)
		durObj, ok := result.(*Duration)
		if !ok {
			t.Fatalf("expected Duration object for %s, got %T", test.method, result)
		}
		
		if durObj.Value != test.expected {
			t.Fatalf("expected duration %v for %s, got %v", test.expected, test.method, durObj.Value)
		}
	}
}

// Test Duration.parse() functionality
func TestDurationParse(t *testing.T) {
	durationNamespace := &DurationNamespace{}
	
	tests := []struct {
		input    string
		expected time.Duration
		valid    bool
	}{
		{"30s", 30 * time.Second, true},
		{"5m", 5 * time.Minute, true},
		{"2h", 2 * time.Hour, true},
		{"2h30m15s", 2*time.Hour + 30*time.Minute + 15*time.Second, true},
		{"invalid", 0, false},
	}
	
	for _, test := range tests {
		result := applyDurationNamespaceMethod(durationNamespace, "parse", &String{Value: test.input})
		
		if test.valid {
			durObj, ok := result.(*Duration)
			if !ok {
				t.Fatalf("expected Duration object for input %s, got %T", test.input, result)
			}
			if durObj.Value != test.expected {
				t.Fatalf("expected duration %v for input %s, got %v", test.expected, test.input, durObj.Value)
			}
		} else {
			if !isError(result) {
				t.Fatalf("expected error for invalid input %s, got %T", test.input, result)  
			}
		}
	}
}

// Test TimeZone creation methods
func TestTimeZoneCreation(t *testing.T) {
	timezoneNamespace := &TimeZoneNamespace{}
	
	// Test TimeZone.utc()
	result := applyTimeZoneNamespaceMethod(timezoneNamespace, "utc")
	tzObj, ok := result.(*TimeZone)
	if !ok {
		t.Fatalf("expected TimeZone object, got %T", result)
	}
	if tzObj.Location != time.UTC {
		t.Fatal("TimeZone.utc() did not return UTC timezone")
	}
	
	// Test TimeZone.local()
	result = applyTimeZoneNamespaceMethod(timezoneNamespace, "local")
	tzObj, ok = result.(*TimeZone)
	if !ok {
		t.Fatalf("expected TimeZone object, got %T", result)
	}
	if tzObj.Location != time.Local {
		t.Fatal("TimeZone.local() did not return local timezone")
	}
	
	// Test TimeZone.parse()
	result = applyTimeZoneNamespaceMethod(timezoneNamespace, "parse", &String{Value: "UTC"})
	tzObj, ok = result.(*TimeZone)
	if !ok {
		t.Fatalf("expected TimeZone object, got %T", result)
	}
	if tzObj.Location.String() != "UTC" {
		t.Fatal("TimeZone.parse('UTC') did not return UTC timezone")
	}
}

// Test Time instance methods
func TestTimeInstanceMethods(t *testing.T) {
	// Create a test time
	testTime := &Time{Value: time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)}
	
	// Test format method
	result := applyTimeMethod(&TimeMethod{Time: testTime, Method: "format"}, 
		[]Value{&String{Value: "2006-01-02 15:04:05"}}, nil)
	strResult, ok := result.(*String)
	if !ok {
		t.Fatalf("expected String result from format, got %T", result)
	}
	if strResult.Value != "2024-01-15 14:30:45" {
		t.Fatalf("expected '2024-01-15 14:30:45', got '%s'", strResult.Value)
	}
	
	// Test format_iso method
	result = applyTimeMethod(&TimeMethod{Time: testTime, Method: "format_iso"}, []Value{}, nil)
	strResult, ok = result.(*String)
	if !ok {
		t.Fatalf("expected String result from format_iso, got %T", result)
	}
	expected := "2024-01-15T14:30:45Z"
	if strResult.Value != expected {
		t.Fatalf("expected '%s', got '%s'", expected, strResult.Value)
	}
	
	// Test to_utc method
	result = applyTimeMethod(&TimeMethod{Time: testTime, Method: "to_utc"}, []Value{}, nil)
	timeResult, ok := result.(*Time)
	if !ok {
		t.Fatalf("expected Time result from to_utc, got %T", result)
	}
	if timeResult.Value.Location() != time.UTC {
		t.Fatal("to_utc did not return UTC time")
	}
}

// Test Duration instance methods
func TestDurationInstanceMethods(t *testing.T) {
	// Create test durations
	dur1 := &Duration{Value: 2*time.Hour + 30*time.Minute}
	dur2 := &Duration{Value: 1 * time.Hour}
	
	// Test add method
	result := applyDurationMethod(&DurationMethod{Duration: dur1, Method: "add"}, 
		[]Value{dur2}, nil)
	durResult, ok := result.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from add, got %T", result)
	}
	expected := 3*time.Hour + 30*time.Minute
	if durResult.Value != expected {
		t.Fatalf("expected duration %v, got %v", expected, durResult.Value)
	}
	
	// Test subtract method
	result = applyDurationMethod(&DurationMethod{Duration: dur1, Method: "subtract"}, 
		[]Value{dur2}, nil)
	durResult, ok = result.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from subtract, got %T", result)
	}
	expected = 1*time.Hour + 30*time.Minute
	if durResult.Value != expected {
		t.Fatalf("expected duration %v, got %v", expected, durResult.Value)
	}
	
	// Test multiply method
	result = applyDurationMethod(&DurationMethod{Duration: dur2, Method: "multiply"}, 
		[]Value{&Integer{Value: 2}}, nil)
	durResult, ok = result.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from multiply, got %T", result)
	}
	expected = 2 * time.Hour
	if durResult.Value != expected {
		t.Fatalf("expected duration %v, got %v", expected, durResult.Value)
	}
	
	// Test abs method with negative duration
	negativeDur := &Duration{Value: -1 * time.Hour}
	result = applyDurationMethod(&DurationMethod{Duration: negativeDur, Method: "abs"}, 
		[]Value{}, nil)
	durResult, ok = result.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from abs, got %T", result)
	}
	if durResult.Value != time.Hour {
		t.Fatalf("expected positive duration, got %v", durResult.Value)
	}
}

// Test Time comparison methods
func TestTimeComparisons(t *testing.T) {
	time1 := &Time{Value: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)}
	time2 := &Time{Value: time.Date(2024, 1, 15, 15, 30, 0, 0, time.UTC)} // 1 hour later
	
	// Test is_before?
	result := applyTimeMethod(&TimeMethod{Time: time1, Method: "is_before?"}, 
		[]Value{time2}, nil)
	boolResult, ok := result.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from is_before?, got %T", result)
	}
	if !boolResult.Value {
		t.Fatal("time1 should be before time2")
	}
	
	// Test is_after?
	result = applyTimeMethod(&TimeMethod{Time: time2, Method: "is_after?"}, 
		[]Value{time1}, nil)
	boolResult, ok = result.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from is_after?, got %T", result)
	}
	if !boolResult.Value {
		t.Fatal("time2 should be after time1")
	}
	
	// Test is_equal?
	result = applyTimeMethod(&TimeMethod{Time: time1, Method: "is_equal?"}, 
		[]Value{time1}, nil)
	boolResult, ok = result.(*Boolean)
	if !ok {
		t.Fatalf("expected Boolean result from is_equal?, got %T", result)
	}
	if !boolResult.Value {
		t.Fatal("time1 should be equal to itself")
	}
}

// Test Time and Duration integration
func TestTimeAndDurationIntegration(t *testing.T) {
	baseTime := &Time{Value: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)}
	duration := &Duration{Value: 2 * time.Hour}
	
	// Test add_duration
	result := applyTimeMethod(&TimeMethod{Time: baseTime, Method: "add_duration"}, 
		[]Value{duration}, nil)
	timeResult, ok := result.(*Time)
	if !ok {
		t.Fatalf("expected Time result from add_duration, got %T", result)
	}
	
	expected := time.Date(2024, 1, 15, 16, 30, 0, 0, time.UTC)
	if !timeResult.Value.Equal(expected) {
		t.Fatalf("expected time %v, got %v", expected, timeResult.Value)
	}
	
	// Test subtract_duration
	result = applyTimeMethod(&TimeMethod{Time: baseTime, Method: "subtract_duration"}, 
		[]Value{duration}, nil)
	timeResult, ok = result.(*Time)
	if !ok {
		t.Fatalf("expected Time result from subtract_duration, got %T", result)
	}
	
	expected = time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	if !timeResult.Value.Equal(expected) {
		t.Fatalf("expected time %v, got %v", expected, timeResult.Value)
	}
	
	// Test difference
	laterTime := &Time{Value: time.Date(2024, 1, 15, 16, 30, 0, 0, time.UTC)}
	result = applyTimeMethod(&TimeMethod{Time: laterTime, Method: "difference"}, 
		[]Value{baseTime}, nil)
	durResult, ok := result.(*Duration)
	if !ok {
		t.Fatalf("expected Duration result from difference, got %T", result)
	}
	
	if durResult.Value != 2*time.Hour {
		t.Fatalf("expected 2 hour difference, got %v", durResult.Value)
	}
}

// Test TimeZone methods
func TestTimeZoneMethods(t *testing.T) {
	utcTz := &TimeZone{Location: time.UTC}
	
	// Test offset method
	result := applyTimeZoneMethod(&TimeZoneMethod{TimeZone: utcTz, Method: "offset"}, 
		[]Value{}, nil)
	intResult, ok := result.(*Integer)
	if !ok {
		t.Fatalf("expected Integer result from offset, got %T", result)
	}
	if intResult.Value != 0 {
		t.Fatalf("expected UTC offset of 0, got %d", intResult.Value)
	}
	
	// Test abbreviation method
	result = applyTimeZoneMethod(&TimeZoneMethod{TimeZone: utcTz, Method: "abbreviation"}, 
		[]Value{}, nil)
	strResult, ok := result.(*String)
	if !ok {
		t.Fatalf("expected String result from abbreviation, got %T", result)
	}
	if strResult.Value != "UTC" {
		t.Fatalf("expected UTC abbreviation, got %s", strResult.Value)
	}
}

// Test error cases
func TestTimeErrorCases(t *testing.T) {
	timeNamespace := &TimeNamespace{}
	
	// Test Time.now() with arguments (should fail)
	result := applyTimeNamespaceMethod(timeNamespace, "now", &String{Value: "invalid"})
	if !isError(result) {
		t.Fatal("Time.now() with arguments should return error")
	}
	
	// Test Time.parse() with invalid argument type
	result = applyTimeNamespaceMethod(timeNamespace, "parse", &Integer{Value: 42})
	if !isError(result) {
		t.Fatal("Time.parse() with integer should return error")
	}
	
	// Test Duration creation with invalid argument type
	durationNamespace := &DurationNamespace{}
	result = applyDurationNamespaceMethod(durationNamespace, "seconds", &String{Value: "invalid"})
	if !isError(result) {
		t.Fatal("Duration.seconds() with string should return error")
	}
}