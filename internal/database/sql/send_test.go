package database

import (
	"context"
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

var listOfSendInstances []*Send

func TestQueries_CreateSend(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		sample CreateSendParams
	}{
		{
			name: "Should pass and returned Send instance has values as expected",
			sample: CreateSendParams{
				Url:       "text",
				File:      "note.txt",
				Size:      "1B",
				Permanent: true,
			},
		},
		{
			name: "Should pass and returned Send instance has values as expected and has permanent value of `false`",
			sample: CreateSendParams{
				Url:  "doc",
				File: "note.docx",
				Size: "2Kb",
			},
		},
		{
			name: "Should pass and returned Send instance has values as expected and has permanent value of `true`",
			sample: CreateSendParams{
				Url:       "zip",
				File:      "files.zip",
				Size:      "6Mb",
				Permanent: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sn, err := testQueries.CreateSend(ctx, tc.sample)
			require.NoError(t, err)

			assert.NotEmpty(t, sn.ID)
			assert.Equal(t, tc.sample.Url, sn.Url)
			assert.Equal(t, tc.sample.File, sn.File)
			assert.Equal(t, tc.sample.Size, sn.Size)
			assert.Equal(t, tc.sample.Permanent, sn.Permanent)
			listOfSendInstances = append(listOfSendInstances, &sn)
		})
	}
}

func TestQueries_DeleteSend(t *testing.T) {
	err := testQueries.DeleteSend(context.Background(), listOfSendInstances[0].ID)
	require.NoError(t, err)

	// remove 1st element from list
	listOfSendInstances = listOfSendInstances[1:]
}

func TestQueries_GetSend(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name     string
		sampleID int64
		wantErr  bool
	}{
		{
			name:     "Should pass using ID that exist and has values as expected",
			sampleID: listOfSendInstances[0].ID,
		},
		{
			name:     "Should fail using non exist ID and return empty instance",
			sampleID: 1001,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sn, err := testQueries.GetSend(ctx, tc.sampleID)
			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, *listOfSendInstances[0], sn)
			case true:
				require.Error(t, err)
				assert.Empty(t, sn, "Not found data should return empty instance")
			}
		})
	}
}

func TestQueries_GetSendByUrl(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name      string
		sampleUrl string
		wantErr   bool
	}{
		{
			name:      "Should pass using Url that exist and has values as expected",
			sampleUrl: listOfSendInstances[0].Url,
		},
		{
			name:      "Should fail using non exist Url and return empty instance",
			sampleUrl: "not-exist-url",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sn, err := testQueries.GetSendByUrl(ctx, tc.sampleUrl)
			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, *listOfSendInstances[0], sn)
			case true:
				require.Error(t, err)
				assert.Empty(t, sn, "Not found data should return empty instance")
			}
		})
	}
}

func TestQueries_ListSend(t *testing.T) {
	sends, err := testQueries.ListSend(context.Background())
	require.NoError(t, err)

	assert.Less(t, 1, len(sends), "Should has more than 1 data after deleting only one in previous test")
	for i := range sends {
		for _, sn := range listOfSendInstances {
			// if the ID same then make sure they are equal
			if sn.ID == sends[i].ID {
				newSn := *sn
				assert.Equal(t, newSn, sends[i])
			}
		}
	}
}

func TestQueries_UpdateSend(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		sample UpdateSendParams
	}{
		{
			name: "Should pass updating 1st data and has `url` field with new value",
			sample: UpdateSendParams{
				ID:        listOfSendInstances[0].ID,
				Url:       "newFile",
				File:      listOfSendInstances[0].File,
				Size:      listOfSendInstances[0].Size,
				Permanent: listOfSendInstances[0].Permanent,
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "Should pass updating 2nd data and has `file` & `size` fields with new value",
			sample: UpdateSendParams{
				ID:        listOfSendInstances[1].ID,
				Url:       listOfSendInstances[1].Url,
				File:      "newFile.rar",
				Size:      "11Mb",
				Permanent: listOfSendInstances[1].Permanent,
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sn, err := testQueries.UpdateSend(ctx, tc.sample)
			require.NoError(t, err)

			assert.Equal(t, tc.sample.ID, sn.ID)
			assert.Equal(t, tc.sample.Url, sn.Url)
			assert.Equal(t, tc.sample.File, sn.File)
			assert.Equal(t, tc.sample.Size, sn.Size)
			assert.Equal(t, tc.sample.Permanent, sn.Permanent)
			assert.NotEqual(t, tc.sample.UpdatedAt, sn.UpdatedAt)
		})
	}
}
