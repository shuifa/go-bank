package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/shuifa/go-bank/db/mock"
	db "github.com/shuifa/go-bank/db/sqlc"
	"github.com/shuifa/go-bank/token"
	"github.com/shuifa/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name        string
		id          int64
		setupAuth   func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		codeStub    func(store *mockdb.MockStore)
		checkResult func(t *testing.T, record *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorizationHeader(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			codeStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResult: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, record.Code)
				checkResponse(t, record.Body, account)
			},
		},
		{
			name: "NotFount",
			id:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorizationHeader(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			codeStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResult: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, record.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			id:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorizationHeader(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			codeStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResult: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, record.Code)
			},
		},
		{
			name: "StatusBadRequest",
			id:   0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				AddAuthorizationHeader(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			codeStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResult: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, record.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)

			server := NewTestServer(t, store)

			testCase.codeStub(store)

			record := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", testCase.id)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			if testCase.setupAuth != nil {
				testCase.setupAuth(t, request, server.tokenMaker)
			}
			server.route.ServeHTTP(record, request)
			testCase.checkResult(t, record)
		})
	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func checkResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	byteData, err := io.ReadAll(body)
	require.NoError(t, err)

	var account1 db.Account

	err = json.Unmarshal(byteData, &account1)
	require.NoError(t, err)

	require.Equal(t, account, account1)
}
