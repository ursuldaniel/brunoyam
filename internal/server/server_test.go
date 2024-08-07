package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
	"github.com/ursuldaniel/brunoyam/mocks"
)

func TestHandleListUsers(t *testing.T) {
	type want struct {
		code  int
		users string
	}

	type test struct {
		name    string
		request string
		method  string
		users   []*models.User
		want    want
	}

	var srv Server
	r := gin.Default()
	r.GET("/", srv.handleListUsers)
	httpSrv := httptest.NewServer(r)

	tests := []test{
		{
			name:    "Test 'handleListUsers' #1; Default call",
			request: "/",
			method:  http.MethodGet,
			users: []*models.User{
				{
					Id:       1,
					Name:     "Vitya",
					Email:    "ex1@ya.ru",
					Password: "pass1",
				},
				{
					Id:       2,
					Name:     "Danya",
					Email:    "ex2@ya.ru",
					Password: "pass2",
				},
			},
			want: want{
				code:  http.StatusOK,
				users: `[{"id":1,"name":"Vitya","email":"ex1@ya.ru","password":"pass1"},{"id":2,"name":"Danya","email":"ex2@ya.ru","password":"pass2"}]`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mocks.NewMockStorage(ctrl)
			defer ctrl.Finish()

			m.EXPECT().ListUsers().Return(tc.users, nil)
			srv.store = m

			getReq := resty.New().R()
			getReq.Method = tc.method
			getReq.URL = httpSrv.URL + tc.request

			resp, err := getReq.Send()

			assert.NoError(t, err)
			assert.Equal(t, tc.want.users, string(resp.Body()))
			assert.Equal(t, tc.want.code, resp.StatusCode())
		})
	}

	httpSrv.Close()
}
