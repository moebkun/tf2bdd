package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/leighmacdonald/steamid/v4/steamid"
	"github.com/stretchr/testify/require"
)

func newTestDB(ctx context.Context) (*sql.DB, error) {
	db, errDB := openDB(":memory:")
	if errDB != nil {
		return nil, errDB
	}

	return db, setupDB(ctx, db)
}

func TestHandleGetSteamIDS(t *testing.T) {
	ctx := context.Background()
	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/steamids", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	recorder := httptest.NewRecorder()
	database, errApp := newTestDB(ctx)
	require.NoError(t, errApp)

	localPlayers := []Player{
		{
			SteamID:    steamid.New(76561198237337976),
			Attributes: []string{"cheater"},
			LastSeen:   LastSeen{},
		},
		{
			SteamID:    steamid.New(76561198834913692),
			Attributes: []string{"cheater"},
			LastSeen:   LastSeen{},
		},
	}
	for _, p := range localPlayers {
		require.NoError(t, addPlayer(ctx, database, p, 0))
	}

	createRouter(database).ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	var players PlayerListRoot
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&players))
	require.Equal(t, len(localPlayers), len(players.Players))
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
