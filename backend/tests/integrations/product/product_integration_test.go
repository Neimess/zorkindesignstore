package product_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	testsuite "github.com/Neimess/zorkin-store-project/pkg/database/test_suite"
	http_utils "github.com/Neimess/zorkin-store-project/pkg/http_utils"
	migrator "github.com/Neimess/zorkin-store-project/pkg/migrator"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func createCategory(t *testing.T, db *sqlx.DB, name string) int64 {
	var id int64
	err := db.QueryRow(`INSERT INTO categories(name) VALUES($1) RETURNING category_id`, name).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestProductEndpoints_TableDriven(t *testing.T) {
	srv := testsuite.RunTestServer(t)
	require.NotNil(t, srv)

	db := srv.App.DB()
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	})
	require.NoError(t, migrator.Run(srv.Cfg.Storage.DSN(), migrator.Options{Mode: migrator.Up}))

	catID := createCategory(t, db, "TestCategory")
	client := srv.Server.Client()
	baseURL := srv.Server.URL
	token := srv.GenerateTestJWT(t, 1)
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	type testCase struct {
		name       string
		method     string
		route      string
		body       any
		expectCode int
		check      func(t *testing.T, data []byte)
	}

	// Pre-create a product for GET/PUT/DELETE tests
	prodID := func() int64 {
		// create sample product
		req := map[string]any{
			"name":        "Init Product",
			"price":       10.0,
			"description": "Desc",
			"category_id": catID,
			"image_url":   "http://example.com/image.jpg",
			"attributes":  []map[string]any{{"name": "Attr", "unit": "", "value": "Val"}},
		}
		data, code := http_utils.DoRequest(t, client, http.MethodPost, baseURL+"/api/admin/product", req, headers)
		require.Equal(t, http.StatusCreated, code, "Expected status code 201 while create init product, got %d", code)
		var created map[string]any
		require.NoError(t, json.Unmarshal(data, &created))
		return int64(created["product_id"].(float64))
	}()

	cases := []testCase{
		{
			name:   "Create_Success",
			method: http.MethodPost,
			route:  baseURL + "/api/admin/product",
			body: map[string]any{
				"name":        "NewProd",
				"price":       5.5,
				"description": "Description",
				"category_id": catID,
				"image_url":   "http://example.com/image.jpg",
				"attributes":  []map[string]any{{"name": "Attr", "unit": "кг", "value": "V"}},
			},
			expectCode: http.StatusCreated,
			check: func(t *testing.T, data []byte) {
				var obj map[string]any
				require.NoError(t, json.Unmarshal(data, &obj))
				require.NotZero(t, obj["product_id"])
			},
		},
		{
			name:       "Create_ValidationErrName",
			method:     http.MethodPost,
			route:      baseURL + "/api/admin/product",
			body:       map[string]any{"name": "", "price": 0, "category_id": catID},
			expectCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Get_Success",
			method:     http.MethodGet,
			route:      fmt.Sprintf("%s/api/product/%d", baseURL, prodID),
			body:       nil,
			expectCode: http.StatusOK,
			check: func(t *testing.T, data []byte) {
				var obj map[string]any
				require.NoError(t, json.Unmarshal(data, &obj))
				require.Equal(t, prodID, int64(obj["product_id"].(float64)))
			},
		},
		{
			name:       "Get_NotFound",
			method:     http.MethodGet,
			route:      fmt.Sprintf("%s/api/product/9999", baseURL),
			body:       nil,
			expectCode: http.StatusNotFound,
		},
		{
			name:       "List",
			method:     http.MethodGet,
			route:      baseURL + "/api/product/category/1",
			body:       nil,
			expectCode: http.StatusOK,
			check: func(t *testing.T, data []byte) {
				var list []map[string]any
				require.NoError(t, json.Unmarshal(data, &list))
				require.GreaterOrEqual(t, len(list), 1)
			},
		},
		{
			name:   "Update_Success",
			method: http.MethodPut,
			route:  fmt.Sprintf("%s/api/admin/product/%d", baseURL, prodID),
			body: map[string]any{
				"name":        "Upd",
				"price":       7.7,
				"description": "Xext",
				"category_id": catID,
				"image_url":   "http://example.com/new.jpg",
				"attributes":  []map[string]any{{"name": "A2", "unit": "", "value": "V2"}},
			},
			expectCode: http.StatusOK,
			check: func(t *testing.T, data []byte) {
				var obj map[string]any
				require.NoError(t, json.Unmarshal(data, &obj))
				require.Equal(t, "Upd", obj["name"])
			},
		},
		{
			name:       "Update_NotFound",
			method:     http.MethodPut,
			route:      fmt.Sprintf("%s/api/admin/product/9999", baseURL),
			body:       map[string]any{"name": "Name", "price": 1, "category_id": catID, "image_url": "http://example.com/new.jpg", "attributes": []map[string]any{{"name": "Attr", "unit": "", "value": "V"}}},
			expectCode: http.StatusNotFound,
		},
		{
			name:       "Delete_Success",
			method:     http.MethodDelete,
			route:      fmt.Sprintf("%s/api/admin/product/%d", baseURL, prodID),
			body:       nil,
			expectCode: http.StatusNoContent,
		},
		{
			name:       "Delete_NotFound_Idempotent",
			method:     http.MethodDelete,
			route:      fmt.Sprintf("%s/api/admin/product/9999", baseURL),
			body:       nil,
			expectCode: http.StatusNoContent,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test case: %+v", tc)
			data, code := http_utils.DoRequest(t, client, tc.method, tc.route, tc.body, headers)
			require.Equalf(
				t,
				tc.expectCode,
				code,
				"expected status code %d, got %d",
				tc.expectCode,
				code,
			)

			if tc.check != nil {
				tc.check(t, data)
			}
		})
	}
}
