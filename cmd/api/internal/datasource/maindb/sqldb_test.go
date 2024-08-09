package maindb_test

import (
	"go-playground/cmd/api/internal/datasource/maindb"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDB(t *testing.T) {
	tests := map[string]struct {
		envMap      map[string]string
		expectError bool
	}{
		"success": {
			envMap: map[string]string{
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_password",
				"DB_ADDRESS":  "localhost:3306",
				"DB_NAME":     "test_db",
			},
		},
		"env is unset": {
			expectError: true,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			for key, val := range v.envMap {
				t.Setenv(key, val)
			}
			db, queries, err := maindb.NewQueryDB(os.LookupEnv)
			if v.expectError {
				require.Error(t, err)
				require.Nil(t, db)
				require.Nil(t, queries)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				require.NotNil(t, queries)
			}
		})
	}
}
