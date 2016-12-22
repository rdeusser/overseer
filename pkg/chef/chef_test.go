package chef

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNewHost(t *testing.T) {
	cases := []struct {
		File       string
		ChefServer string
	}{
		{"testkey1", "https://localhost/"},
		{"testkey2", "https://localhost/"},
		{"testkey1", "https://localhost/"},
		{"testkey2", "https://localhost/"},
		{"testkey1", "https://localhost/"},
		{"testkey2", "https://localhost/"},
		{"testkey1", "https://localhost/"},
		{"testkey2", "https://localhost/"},
		{"testkey1", "https://localhost/"},
		{"testkey2", "https://localhost/"},
	}

	for _, tt := range cases {
		path, err := filepath.Abs(filepath.Join("./test-fixtures", tt.File))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		key, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		_, err = NewClient(string(key), tt.ChefServer)
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}
	}
}

func TestReadValidationKey(t *testing.T) {}
