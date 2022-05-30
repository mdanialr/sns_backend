package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBPostgres_GenerateConnectionString(t *testing.T) {
	testCases := []struct {
		name   string
		sample DBPostgres
		expect string
	}{
		{
			name:   "Default value of `User` field is 'postgres' if not provided",
			sample: DBPostgres{},
			expect: "postgresql://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable",
		},
		{
			name:   "Default value of `Pass` field is 'postgres' if not provided",
			sample: DBPostgres{User: "user"},
			expect: "postgresql://user:postgres@127.0.0.1:5432/postgres?sslmode=disable",
		},
		{
			name:   "Default value of `Host` field is '127.0.0.1' if not provided",
			sample: DBPostgres{User: "user", Pass: "secret"},
			expect: "postgresql://user:secret@127.0.0.1:5432/postgres?sslmode=disable",
		},
		{
			name:   "Default value of `Port` field is '5432' if not provided",
			sample: DBPostgres{User: "user", Pass: "secret", Host: "localhost"},
			expect: "postgresql://user:secret@localhost:5432/postgres?sslmode=disable",
		},
		{
			name:   "Default value of `Name` field is 'postgres' if not provided",
			sample: DBPostgres{User: "user", Pass: "secret", Host: "localhost", Port: 5430},
			expect: "postgresql://user:secret@localhost:5430/postgres?sslmode=disable",
		},
		{
			name:   "SSL mode should be 'disable' by default because default value of `IsSSL` field is false",
			sample: DBPostgres{User: "user", Pass: "secret", Host: "localhost", Port: 5430, Name: "sns"},
			expect: "postgresql://user:secret@localhost:5430/sns?sslmode=disable",
		},
		{
			name:   "SSL mode should be 'enable' when the provided `IsSSL` field is true",
			sample: DBPostgres{User: "user", Pass: "secret", Host: "localhost", Port: 5430, Name: "sns", IsSSL: true},
			expect: "postgresql://user:secret@localhost:5430/sns?sslmode=enable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.sample.GenerateConnectionString()
			assert.Equal(t, tc.expect, out)
		})
	}
}

func TestDBPostgres_GetDriver(t *testing.T) {
	t.Run("Should always return the same value which is 'postgres'", func(t *testing.T) {
		sample := DBPostgres{}
		assert.Equal(t, "postgres", sample.GetDriver())
	})
}
