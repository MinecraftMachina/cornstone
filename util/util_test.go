package util

import (
	"strings"
	"testing"
)

func TestSafePath(t *testing.T) {
	var tests = []struct {
		unsafePath string
	}{
		{"/tmp/base/etc/ssh/config"},
		{"/a/../../.."},
		{"../../../etc/ssh/config"},
		{"../../../etc/../../../ssh/config"},
		{"../etc/ssh/config"},
		{"/etc/../../../ssh/config"},
		{"C:\\Test"},
		{"C:\\Test\\..\\..\\..\\.."},
		{"..\\..\\Test"},
		{"C:\\..\\..\\Test"},
		{"..\\..\\ASD\\..\\..\\Test"},
		{"..\\..\\Test"},
	}

	basePath := "/the/safe"

	for _, test := range tests {
		if strings.Contains(test.unsafePath, "safe") {
			t.Error("unsafePath not allowed to contain 'safe'")
		}
		result := SafeJoin(basePath, test.unsafePath)
		if !strings.Contains(result, "safe") {
			t.Error("Failed case:", test.unsafePath, ", with result: "+result)
		} else {
			t.Log("Passed case:", test.unsafePath, ", with result: "+result)
		}
	}
}
