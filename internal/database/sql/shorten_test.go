package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: all tests in this file using the same list of Shorten instance
// Test steps:
// 1. Create 3 data then assert them
// 2. Delete 1st data so there are should be 2 data left then assert them
// 3. Get 1 data then assert them
// 4. List 2 data then assert them
// 5. Update 1st data then assert them

var listOfShortenInstances []*Shorten

func TestQueries_CreateShorten(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		sample CreateShortenParams
	}{
		{
			name: "Should pass and returned Shorten instance has expected value as the one inputted as sample",
			sample: CreateShortenParams{
				Url:       "gl",
				Target:    "https://google.com",
				Permanent: true,
			},
		},
		{
			name: "Should pass and returned instance has value as expected and has permanent value of `false`",
			sample: CreateShortenParams{
				Url:    "yt",
				Target: "https://youtube.com",
			},
		},
		{
			name: "Should pass and returned instance has value as expected and has permanent value of `true`",
			sample: CreateShortenParams{
				Url:       "go",
				Target:    "https://go.dev",
				Permanent: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sh, err := testQueries.CreateShorten(ctx, tc.sample)
			require.NoError(t, err)

			assert.NotEmpty(t, sh.ID)
			assert.Equal(t, tc.sample.Url, sh.Url)
			assert.Equal(t, tc.sample.Target, sh.Target)
			assert.Equal(t, tc.sample.Permanent, sh.Permanent)
			listOfShortenInstances = append(listOfShortenInstances, &sh)
		})
	}
}

func TestQueries_DeleteShorten(t *testing.T) {
	err := testQueries.DeleteShorten(context.Background(), listOfShortenInstances[0].ID)
	require.NoError(t, err)

	// remove first element from list
	listOfShortenInstances = listOfShortenInstances[1:]
}

func TestQueries_GetShorten(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name     string
		sampleID int64
		wantErr  bool
	}{
		{
			name:     "Should pass and has value as expected",
			sampleID: listOfShortenInstances[0].ID,
		},
		{
			name:     "Should fail when searching using non exist ID",
			sampleID: 1001,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sh, err := testQueries.GetShorten(ctx, tc.sampleID)
			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, *listOfShortenInstances[0], sh)
			case true:
				require.Error(t, err)
				assert.Empty(t, sh, "Not found data should be empty")
			}
		})
	}
}

func TestQueries_GetShortenByUrl(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name      string
		sampleUrl string
		wantErr   bool
	}{
		{
			name:      "Should pass and has value as expected",
			sampleUrl: listOfShortenInstances[0].Url,
		},
		{
			name:      "Should fail when searching using non exist url",
			sampleUrl: "not-exist-url",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sh, err := testQueries.GetShortenByUrl(ctx, tc.sampleUrl)
			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, *listOfShortenInstances[0], sh)
			case true:
				require.Error(t, err)
				assert.Empty(t, sh, "Not found data should be empty")
			}
		})
	}
}

func TestQueries_ListShorten(t *testing.T) {
	shortens, err := testQueries.ListShorten(context.Background(), "updated_at", ASC)
	require.NoError(t, err)

	assert.Less(t, 1, len(shortens), "Should has more than 1 data after deleting only one in previous test")
	for i := range shortens {
		for _, sh := range listOfShortenInstances {
			// if the ID same then make sure they are equal
			if sh.ID == shortens[i].ID {
				newSh := *sh
				assert.Equal(t, newSh, shortens[i])
			}
		}
	}
}

func TestQueries_UpdateShorten(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		sample UpdateShortenParams
	}{
		{
			name: "Should pass updating 1st data and has `target` field with new value as expected",
			sample: UpdateShortenParams{
				ID:        listOfShortenInstances[0].ID,
				Url:       listOfShortenInstances[0].Url,
				Target:    "https://new.example.com",
				Permanent: listOfShortenInstances[0].Permanent,
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			},
		},
		{
			name: "Should pass updating 2nd data and has `url` field with new value as expected",
			sample: UpdateShortenParams{
				ID:        listOfShortenInstances[1].ID,
				Url:       "new",
				Target:    listOfShortenInstances[1].Target,
				Permanent: listOfShortenInstances[1].Permanent,
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sh, err := testQueries.UpdateShorten(ctx, tc.sample)
			require.NoError(t, err)

			assert.Equal(t, tc.sample.ID, sh.ID)
			assert.Equal(t, tc.sample.Url, sh.Url)
			assert.Equal(t, tc.sample.Target, sh.Target)
			assert.Equal(t, tc.sample.Permanent, sh.Permanent)
			assert.NotEqual(t, tc.sample.UpdatedAt, sh.UpdatedAt)
		})
	}
}
