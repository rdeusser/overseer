package buildspec

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		File     string
		Expected *Spec
		Err      bool
	}{
		{
			"virtualspec",
			&Spec{
				Hosts: []string{
					"hello.qa.local",
					"lol.qa.local",
					"with1234.qa.local",
					"nope.qa.local",
					"sometimes@#$@#%123135.qa.local",
				},
				MACs: nil,
			},
			false,
		},
		{
			"physicalspec",
			&Spec{
				Hosts: []string{
					"hello.qa.local",
					"lol.qa.local",
					"with1234.qa.local",
					"nope.qa.local",
					"sometimes@#$@#%123135.qa.local",
				},
				MACs: []string{
					"1C:29:DF:E5:AA:B5",
					"52:65:06:7A:C5:C8",
					"37:25:61:C8:B5:9C",
					"19:62:AD:A7:92:BA",
					"E5:CF:60:13:C2:3E",
				},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Logf("Starting parse on: %s", tt.File)

		path, err := filepath.Abs(filepath.Join("./test-fixtures", tt.File))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		actual, err := ParseFile(path)
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.Expected) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tt.File, actual, tt.Expected)
		}
	}
}
