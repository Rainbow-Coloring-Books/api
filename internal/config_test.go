package internal_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"rainbowcoloringbooks/internal"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		configYAML  string
		expectError bool
	}{
		{
			name: "valid config file",
			configYAML: `
db_user: test_user
db_password: test_password
db_name: test_db
`,
			expectError: false,
		},
		{
			name: "invalid config file",
			configYAML: `
db_user: test_user
db_password test_password
db_name: test_db
`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := ioutil.TempFile("", "config-*.yaml")
			assert.NoError(t, err)

			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString(tc.configYAML)
			assert.NoError(t, err)

			err = tempFile.Close()
			assert.NoError(t, err)

			config, err := internal.LoadConfig(tempFile.Name())

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test_user", config.DBUser)
				assert.Equal(t, "test_password", config.DBPassword)
				assert.Equal(t, "test_db", config.DBName)
			}
		})
	}
}
