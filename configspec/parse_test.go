package configspec

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
			"complete.conf",
			&Spec{
				Foreman: Foreman{
					Username: "admin",
					Password: "datpass",
				},
				Chef: Chef{
					Username:      "admin",
					Password:      "datpass",
					ChefServer:    "https://localhost/",
					ClientKey:     "~/.chef/client.pem",
					ValidationKey: "~/.chef/validation.pem",
				},
				Vsphere: Vsphere{
					Username: "admin",
					Password: "datpass",
				},
				Infoblox: Infoblox{
					Username: "admin",
					Password: "datpass",
				},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Logf("Testing parse on: %s", tt.File)

		path, err := filepath.Abs(filepath.Join("./test-fixtures", tt.File))
		if err != nil {
			t.Fatalf("fie: %s\n\n%s", tt.File, err)
			continue
		}

		actual, err := ParseFile(path)
		if (err != nil) != tt.Err {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.Expected) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tt.File, actual, tt.Expected)
		}
	}
}
