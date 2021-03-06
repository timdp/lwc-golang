package lwc

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

type formatTest struct {
	counts []uint64
	label  string
	cr     bool
	lf     bool
}

func (t *formatTest) expected() []string {
	result := make([]string, len(t.counts))
	for i, num := range t.counts {
		result[i] = fmt.Sprintf("%d", num)
	}
	if t.label != "" {
		result = append(result, t.label)
	}
	return result
}

func withWithout(b bool) string {
	if b {
		return "with"
	}
	return "without"
}

func tokenize(str string) []string {
	tokens := make([]string, 100)
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		tokens[count] = scanner.Text()
		count++
	}
	return tokens[0:count]
}

var formatTests = []formatTest{
	{
		[]uint64{42939},
		"",
		false,
		false,
	},
	{
		[]uint64{42, 2993},
		"bar",
		true,
		false,
	},
	{
		[]uint64{90210},
		"baz-quux",
		false,
		true,
	},
	{
		[]uint64{123, 4567, 899999},
		"/etc/passwd",
		true,
		true,
	},
}

func TestFormatCounts(t *testing.T) {
	for i, test := range formatTests {
		result := formatCounts(&test.counts, test.label, test.cr, test.lf).String()
		hasCr := strings.HasPrefix(result, "\r")
		if test.cr != hasCr {
			t.Errorf("Test #%d failed: expecting string %s CR prefix", i, withWithout(test.cr))
		}
		hasLf := strings.HasSuffix(result, "\n")
		if test.lf != hasLf {
			t.Errorf("Test #%d failed: expecting string %s LF suffix", i, withWithout(test.lf))
		}
		actual := tokenize(result)
		expected := test.expected()
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Test #%d failed: expecting %#v, got %#v", i, expected, actual)
		}
	}
}

func TestPrintNamedCounts(t *testing.T) {
	PrintCounts(&[]uint64{1, 2, 3}, "file", true, true)
	actual := string(lwcutil.FlushStdoutBuffer())
	expected := "\r       1        2        3 file\n"
	if expected != actual {
		t.Errorf("Expecting %#v, got %#v", expected, actual)
	}
}

func TestPrintStdinCounts(t *testing.T) {
	PrintCounts(&[]uint64{1, 2, 3}, "", true, true)
	actual := string(lwcutil.FlushStdoutBuffer())
	expected := "\r       1        2        3\n"
	if expected != actual {
		t.Errorf("Expecting %#v, got %#v", expected, actual)
	}
}
