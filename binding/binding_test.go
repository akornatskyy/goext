package binding_test

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/akornatskyy/goext/binding"
)

func TestBindNotPtr(t *testing.T) {
	var m struct{}
	for _, tt := range []interface{}{nil, m} {
		if err := binding.Bind(tt, nil); err == nil {
			t.Error()
		}
	}
}

func TestBindNoBinding(t *testing.T) {
	var m struct {
		Test string
	}
	for _, tt := range []interface{}{&m} {
		if err := binding.Bind(tt, nil); err != nil {
			t.Error()
		}
	}
}

func TestBindNoValues(t *testing.T) {
	var m struct {
		Test string `binding:"test"`
	}

	if err := binding.Bind(&m, map[string][]string{}); err != nil {
		t.Error()
	}

	if err := binding.Bind(&m, map[string][]string{"test": {}}); err != nil {
		t.Error()
	}
}

func TestBindString(t *testing.T) {
	var m struct {
		Test string `binding:"test"`
	}
	for _, tt := range []string{"", " ", "test", " x", "x ", " x "} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		if m.Test != tt {
			t.Errorf("got %q; expected %q", m.Test, tt)
		}
	}
}

func TestBindInt(t *testing.T) {
	var m struct {
		Test int `binding:"test"`
	}
	for _, tt := range []string{"0", "1", "-1", "-1000", "1000"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		expected, _ := strconv.Atoi(tt)
		if m.Test != expected {
			t.Errorf("got %d; expected %s", m.Test, tt)
		}
	}

	for _, tt := range []string{"", "x", "1x", "x1", "123412312312313123131"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err == nil {
			t.Error()
		}
	}
}

func TestBindUint(t *testing.T) {
	var m struct {
		Test uint `binding:"test"`
	}
	for _, tt := range []string{"0", "1", "5", "1000", "1000000"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		expected, _ := strconv.ParseUint(tt, 10, 0)
		if m.Test != uint(expected) {
			t.Errorf("got %d; expected %s", m.Test, tt)
		}
	}

	for _, tt := range []string{"", "x", "1x", "x1", "-1", "99112312312313123131"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err == nil {
			t.Error()
		}
	}
}

func TestBindBool(t *testing.T) {
	var m struct {
		Test bool `binding:"test"`
	}
	for _, tt := range []string{"0", "1", "t", "f"} {
		b, _ := strconv.ParseBool(tt)
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		if m.Test != b {
			t.Errorf("got %t; expected %t for %q", m.Test, b, tt)
		}
	}

	for _, tt := range []string{"", " ", "x", "11", "no", "yes"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err == nil {
			t.Error()
		}
	}
}

func TestBindDuration(t *testing.T) {
	var m struct {
		Test time.Duration `binding:"test"`
	}
	for _, tt := range []string{"12s", "5m6s", "23h", "3605s"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		expected, _ := time.ParseDuration(tt)
		if m.Test != expected {
			t.Errorf("got %v; expected %v", m.Test, expected)
		}
	}

	for _, tt := range []string{"", " ", "x", "2019"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err == nil {
			t.Error()
		}
	}
}

func TestBindTime(t *testing.T) {
	var m struct {
		Test time.Time `binding:"test"`
	}
	tt := "2019-03-29T9:38:40Z"
	if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
		t.Error()
	}

	expected, _ := time.Parse(time.RFC3339, tt)
	if m.Test != expected {
		t.Errorf("got %v; expected %v", m.Test, expected)
	}
}

func TestBindTimeFailLoc(t *testing.T) {
	var m struct {
		Test time.Time `binding:"test" loc:"X"`
	}
	if err := binding.Bind(&m, map[string][]string{"test": {""}}); err == nil {
		t.Error()
	}
}

func TestBindTimeWithLoc(t *testing.T) {
	var m struct {
		Test time.Time `binding:"test" loc:"EET"`
	}
	tt := "2019-03-29T9:38:40Z"
	if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
		t.Error()
	}
}

func TestBindTimeWithLayout(t *testing.T) {
	var m struct {
		Test time.Time `binding:"test" layout:"2006-01-02"`
	}
	for _, tt := range []string{"2019-03-23", "2019-03-29"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err != nil {
			t.Error()
		}

		expected, _ := time.Parse("2006-01-02", tt)
		if m.Test != expected {
			t.Errorf("got %v; expected %v", m.Test, expected)
		}
	}

	for _, tt := range []string{"", "x", "2019", "2019-03", "2019-01-99"} {
		if err := binding.Bind(&m, map[string][]string{"test": {tt}}); err == nil {
			t.Error()
		}
	}
}

func TestBindSliceError(t *testing.T) {
	var m struct {
		Test []int `binding:"test"`
	}

	if err := binding.Bind(&m, map[string][]string{"test": {"1", "2x", "3"}}); err == nil {
		t.Error()
	}
}

func TestBindSlice(t *testing.T) {
	var m struct {
		Test []string `binding:"test"`
	}
	for _, tt := range [][]string{{"x"}, {"1", "2", "3"}} {
		if err := binding.Bind(&m, map[string][]string{"test": tt}); err != nil {
			t.Error()
		}

		if !reflect.DeepEqual(&m.Test, &tt) {
			t.Errorf("no match")
		}
	}
}

// Benchmarks

type sample struct {
	Query    string        `binding:"q"`
	Page     int           `binding:"page"`
	Size     uint          `binding:"size"`
	Ok       bool          `binding:"ok"`
	Duration time.Duration `binding:"duration"`
	From     time.Time     `binding:"from" layout:"2006-01-02" loc:"UTC"`
	Colors   []string      `binding:"colors"`
	Numbers  []int         `binding:"numbers"`
}

var values = map[string][]string{
	"q":        {"test"},
	"page":     {"1"},
	"size":     {"20"},
	"ok":       {"1"},
	"duration": {"4h30m45s"},
	"from":     {"2019-03-23"},
	"colors":   {"yellow", "blue"},
	"numbers":  {"1", "5", "10", "-20"},
}

func TestBind(t *testing.T) {
	var s sample

	err := binding.Bind(&s, values)

	if err != nil {
		t.Errorf("unexpected error")
	}
	duration, _ := time.ParseDuration("4h30m45s")
	from, _ := time.Parse("2006-01-02", "2019-03-23")
	expected := sample{
		Query:    "test",
		Page:     1,
		Size:     20,
		Ok:       true,
		Duration: duration,
		From:     from,
		Colors:   []string{"yellow", "blue"},
		Numbers:  []int{1, 5, 10, -20},
	}
	if !reflect.DeepEqual(&s, &expected) {
		t.Errorf("no match, %+v ; expected %+v", &s, &expected)
	}
}

func BenchmarkBind(b *testing.B) {
	var s sample
	var err error
	for i := 0; i < b.N; i++ {
		err = binding.Bind(&s, values)
	}
	b.StopTimer()

	if err != nil {
		b.Errorf("unexpected error")
	}
}
