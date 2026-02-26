// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestServer creates a mock HTTP server and an API client pointed at it.
func newTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)
	client := NewClient(server.URL, "test-token", 1)
	return server, client
}

// --- Client Tests ---

// TestNewClient verifies the client is created with correct defaults.
func TestNewClient(t *testing.T) {
	client := NewClient("https://example.com", "my-token", 5)

	if client.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "https://example.com")
	}

	if client.accessToken != "my-token" {
		t.Errorf("accessToken = %q, want %q", client.accessToken, "my-token")
	}

	if client.accountID != 5 {
		t.Errorf("accountID = %d, want %d", client.accountID, 5)
	}
}

// TestSetBaseURL verifies the base URL can be overridden.
func TestSetBaseURL(t *testing.T) {
	client := NewClient("https://old.com", "token", 1)
	client.SetBaseURL("https://new.com")

	if client.baseURL != "https://new.com" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "https://new.com")
	}
}

// TestSetAccountID verifies the account ID can be changed.
func TestSetAccountID(t *testing.T) {
	client := NewClient("https://example.com", "token", 1)
	client.SetAccountID(42)

	if client.accountID != 42 {
		t.Errorf("accountID = %d, want %d", client.accountID, 42)
	}
}

// TestGetSendsAuthHeader verifies that GET requests include the Bearer token.
func TestGetSendsAuthHeader(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, err := client.get("/test", nil)
	if err != nil {
		t.Fatalf("get() error = %v", err)
	}
}

// TestGetWithParams verifies that query parameters are appended to the URL.
func TestGetWithParams(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "25" {
			t.Errorf("limit param = %q, want %q", r.URL.Query().Get("limit"), "25")
		}
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("page param = %q, want %q", r.URL.Query().Get("page"), "2")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, err := client.get("/test", map[string]string{"limit": "25", "page": "2"})
	if err != nil {
		t.Fatalf("get() error = %v", err)
	}
}

// TestGetSkipsEmptyParams verifies that empty parameter values are not sent.
func TestGetSkipsEmptyParams(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("empty") != "" {
			t.Error("empty param should not be sent")
		}
		if r.URL.Query().Get("filled") != "yes" {
			t.Errorf("filled param = %q, want %q", r.URL.Query().Get("filled"), "yes")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, err := client.get("/test", map[string]string{"empty": "", "filled": "yes"})
	if err != nil {
		t.Fatalf("get() error = %v", err)
	}
}

// TestPostSendsJSON verifies that POST requests send JSON and include auth.
func TestPostSendsJSON(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}

		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"name"`) {
			t.Errorf("body = %q, expected to contain 'name'", string(body))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1}`))
	})
	defer server.Close()

	_, err := client.post("/test", map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("post() error = %v", err)
	}
}

// TestPostNoAuthSkipsBearer verifies that postNoAuth does not send an Authorization header.
func TestPostNoAuthSkipsBearer(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("Authorization should be empty, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, err := client.postNoAuth("/test", map[string]string{})
	if err != nil {
		t.Fatalf("postNoAuth() error = %v", err)
	}
}

// TestPutSendsJSON verifies that PUT requests send JSON with auth.
func TestPutSendsJSON(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}

		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, err := client.put("/test", map[string]string{"name": "updated"})
	if err != nil {
		t.Fatalf("put() error = %v", err)
	}
}

// TestDeleteSendsAuth verifies that DELETE requests include the Bearer token.
func TestDeleteSendsAuth(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}

		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	_, err := client.delete("/test")
	if err != nil {
		t.Fatalf("delete() error = %v", err)
	}
}

// TestHTTPErrorReturnsError verifies that non-2xx status codes produce an error.
func TestHTTPErrorReturnsError(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"something went wrong"}`))
	})
	defer server.Close()

	_, err := client.get("/test", nil)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error = %q, expected to contain '500'", err.Error())
	}
}

// TestAccountPath verifies the account path builder.
func TestAccountPath(t *testing.T) {
	client := NewClient("https://example.com", "token", 42)

	path := client.accountPath("/ledger")
	expected := "/api/v3/42/ledger"
	if path != expected {
		t.Errorf("accountPath = %q, want %q", path, expected)
	}
}

// --- Auth Tests ---

// TestLogin verifies the login flow with mock server.
func TestLogin(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("path = %q, want /oauth/token", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req LoginRequest
		json.Unmarshal(body, &req)

		if req.Username != "user@example.com" {
			t.Errorf("username = %q, want user@example.com", req.Username)
		}
		if req.GrantType != "password" {
			t.Errorf("grant_type = %q, want password", req.GrantType)
		}

		resp := LoginResponse{
			AccessToken: "abc123token",
			UserID:      42,
			TokenType:   "bearer",
		}
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	resp, err := client.Login("user@example.com", "password123", "test-client-id")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if resp.AccessToken != "abc123token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "abc123token")
	}
	if resp.UserID != 42 {
		t.Errorf("UserID = %d, want %d", resp.UserID, 42)
	}
}

// TestLoginError verifies login handles API errors correctly.
func TestLoginError(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
	})
	defer server.Close()

	_, err := client.Login("bad@example.com", "wrong", "bad-client-id")
	if err == nil {
		t.Fatal("Login() expected error, got nil")
	}
}

// TestLogout verifies the logout flow.
func TestLogout(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/logout" {
			t.Errorf("path = %q, want /oauth/logout", r.URL.Path)
		}

		if r.URL.Query().Get("access_token") != "test-token" {
			t.Errorf("access_token = %q, want test-token", r.URL.Query().Get("access_token"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	err := client.Logout()
	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
}

// TestGetAuthUser verifies fetching the authenticated user.
func TestGetAuthUser(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/me" {
			t.Errorf("path = %q, want /oauth/me", r.URL.Path)
		}

		user := User{
			ID:        42,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Accounts: []Account{
				{ID: 1, Name: "Personal"},
				{ID: 2, Name: "Business"},
			},
		}
		json.NewEncoder(w).Encode(user)
	})
	defer server.Close()

	user, err := client.GetAuthUser()
	if err != nil {
		t.Fatalf("GetAuthUser() error = %v", err)
	}

	if user.ID != 42 {
		t.Errorf("ID = %d, want %d", user.ID, 42)
	}
	if len(user.Accounts) != 2 {
		t.Errorf("Accounts count = %d, want %d", len(user.Accounts), 2)
	}
}

// --- Account Tests ---

// TestGetAccount verifies fetching account details.
func TestGetAccount(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/account" {
			t.Errorf("path = %q, want /api/v3/1/account", r.URL.Path)
		}

		account := Account{ID: 1, Name: "Test Account", Currency: "USD"}
		json.NewEncoder(w).Encode(account)
	})
	defer server.Close()

	account, err := client.GetAccount()
	if err != nil {
		t.Fatalf("GetAccount() error = %v", err)
	}

	if account.Name != "Test Account" {
		t.Errorf("Name = %q, want %q", account.Name, "Test Account")
	}
	if account.Currency != "USD" {
		t.Errorf("Currency = %q, want %q", account.Currency, "USD")
	}
}

// TestUpdateAccount verifies updating account details.
func TestUpdateAccount(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var acct Account
		json.Unmarshal(body, &acct)

		if acct.Name != "Updated Name" {
			t.Errorf("Name = %q, want %q", acct.Name, "Updated Name")
		}

		acct.ID = 1
		json.NewEncoder(w).Encode(acct)
	})
	defer server.Close()

	updated, err := client.UpdateAccount(&Account{Name: "Updated Name"})
	if err != nil {
		t.Fatalf("UpdateAccount() error = %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", updated.Name, "Updated Name")
	}
}

// TestGetBilling verifies fetching billing information.
func TestGetBilling(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/account/billing" {
			t.Errorf("path = %q, want /api/v3/1/account/billing", r.URL.Path)
		}

		billing := Billing{Status: "Active", Subscription: "Monthly"}
		json.NewEncoder(w).Encode(billing)
	})
	defer server.Close()

	billing, err := client.GetBilling()
	if err != nil {
		t.Fatalf("GetBilling() error = %v", err)
	}

	if billing.Status != "Active" {
		t.Errorf("Status = %q, want %q", billing.Status, "Active")
	}
}

// --- Ledger Tests ---

// TestGetLedgers verifies fetching a list of ledger entries.
func TestGetLedgers(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/ledger" {
			t.Errorf("path = %q, want /api/v3/1/ledger", r.URL.Path)
		}

		if r.URL.Query().Get("limit") != "25" {
			t.Errorf("limit = %q, want 25", r.URL.Query().Get("limit"))
		}

		ledgers := []Ledger{
			{ID: 1, Amount: -50.00, Date: "2026-01-15"},
			{ID: 2, Amount: 100.00, Date: "2026-01-20"},
		}
		json.NewEncoder(w).Encode(ledgers)
	})
	defer server.Close()

	ledgers, err := client.GetLedgers(map[string]string{"limit": "25"})
	if err != nil {
		t.Fatalf("GetLedgers() error = %v", err)
	}

	if len(ledgers) != 2 {
		t.Fatalf("count = %d, want %d", len(ledgers), 2)
	}
	if ledgers[0].Amount != -50.00 {
		t.Errorf("Amount = %f, want %f", ledgers[0].Amount, -50.00)
	}
}

// TestGetLedger verifies fetching a single ledger entry.
func TestGetLedger(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/ledger/42" {
			t.Errorf("path = %q, want /api/v3/1/ledger/42", r.URL.Path)
		}

		ledger := Ledger{ID: 42, Amount: -75.50, Date: "2026-02-01"}
		json.NewEncoder(w).Encode(ledger)
	})
	defer server.Close()

	ledger, err := client.GetLedger(42)
	if err != nil {
		t.Fatalf("GetLedger() error = %v", err)
	}

	if ledger.ID != 42 {
		t.Errorf("ID = %d, want %d", ledger.ID, 42)
	}
}

// TestCreateLedger verifies creating a new ledger entry.
func TestCreateLedger(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v3/1/ledger" {
			t.Errorf("path = %q, want /api/v3/1/ledger", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req LedgerCreateRequest
		json.Unmarshal(body, &req)

		if req.Amount != -25.50 {
			t.Errorf("Amount = %f, want %f", req.Amount, -25.50)
		}

		ledger := Ledger{ID: 99, Amount: req.Amount, Date: req.Date}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ledger)
	})
	defer server.Close()

	ledger, err := client.CreateLedger(&LedgerCreateRequest{
		Amount:     -25.50,
		Date:       "2026-02-25",
		ContactID:  1,
		CategoryID: 2,
	})
	if err != nil {
		t.Fatalf("CreateLedger() error = %v", err)
	}

	if ledger.ID != 99 {
		t.Errorf("ID = %d, want %d", ledger.ID, 99)
	}
}

// TestUpdateLedger verifies updating a ledger entry.
func TestUpdateLedger(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v3/1/ledger/42" {
			t.Errorf("path = %q, want /api/v3/1/ledger/42", r.URL.Path)
		}

		ledger := Ledger{ID: 42, Amount: -30.00}
		json.NewEncoder(w).Encode(ledger)
	})
	defer server.Close()

	ledger, err := client.UpdateLedger(42, &LedgerUpdateRequest{Amount: -30.00})
	if err != nil {
		t.Fatalf("UpdateLedger() error = %v", err)
	}

	if ledger.Amount != -30.00 {
		t.Errorf("Amount = %f, want %f", ledger.Amount, -30.00)
	}
}

// TestDeleteLedger verifies deleting a ledger entry.
func TestDeleteLedger(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/ledger/42" {
			t.Errorf("path = %q, want /api/v3/1/ledger/42", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteLedger(42)
	if err != nil {
		t.Fatalf("DeleteLedger() error = %v", err)
	}
}

// TestGetLedgerSummary verifies fetching a ledger summary.
func TestGetLedgerSummary(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/ledger/summary" {
			t.Errorf("path = %q, want /api/v3/1/ledger/summary", r.URL.Path)
		}

		summary := LedgerSummary{Income: 5000, Expense: 3000, Profit: 2000}
		json.NewEncoder(w).Encode(summary)
	})
	defer server.Close()

	summary, err := client.GetLedgerSummary(nil)
	if err != nil {
		t.Fatalf("GetLedgerSummary() error = %v", err)
	}

	if summary.Profit != 2000 {
		t.Errorf("Profit = %f, want %f", summary.Profit, 2000.0)
	}
}

// TestGetLedgerPL verifies fetching the ledger P&L.
func TestGetLedgerPL(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/ledger/pl" {
			t.Errorf("path = %q, want /api/v3/1/ledger/pl", r.URL.Path)
		}

		report := PnlReport{Income: 10000, Expense: 7000, Profit: 3000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetLedgerPL(nil)
	if err != nil {
		t.Fatalf("GetLedgerPL() error = %v", err)
	}

	if report.Profit != 3000 {
		t.Errorf("Profit = %f, want %f", report.Profit, 3000.0)
	}
}

// --- Category Tests ---

// TestGetCategories verifies fetching a list of categories.
func TestGetCategories(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/categories" {
			t.Errorf("path = %q, want /api/v3/1/categories", r.URL.Path)
		}

		categories := []Category{
			{ID: 1, Name: "Meals", Type: "1"},
			{ID: 2, Name: "Sales", Type: "2"},
		}
		json.NewEncoder(w).Encode(categories)
	})
	defer server.Close()

	categories, err := client.GetCategories(nil)
	if err != nil {
		t.Fatalf("GetCategories() error = %v", err)
	}

	if len(categories) != 2 {
		t.Fatalf("count = %d, want %d", len(categories), 2)
	}
	if categories[0].Name != "Meals" {
		t.Errorf("Name = %q, want %q", categories[0].Name, "Meals")
	}
}

// TestGetCategory verifies fetching a single category.
func TestGetCategory(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/categories/5" {
			t.Errorf("path = %q, want /api/v3/1/categories/5", r.URL.Path)
		}

		category := Category{ID: 5, Name: "Travel", Type: "1"}
		json.NewEncoder(w).Encode(category)
	})
	defer server.Close()

	category, err := client.GetCategory(5)
	if err != nil {
		t.Fatalf("GetCategory() error = %v", err)
	}

	if category.Name != "Travel" {
		t.Errorf("Name = %q, want %q", category.Name, "Travel")
	}
}

// TestCreateCategory verifies creating a new category.
func TestCreateCategory(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req CategoryCreateRequest
		json.Unmarshal(body, &req)

		if req.Name != "Office Supplies" {
			t.Errorf("Name = %q, want %q", req.Name, "Office Supplies")
		}

		category := Category{ID: 10, Name: req.Name, Type: req.Type}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(category)
	})
	defer server.Close()

	category, err := client.CreateCategory(&CategoryCreateRequest{
		Name: "Office Supplies",
		Type: "1",
	})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	if category.ID != 10 {
		t.Errorf("ID = %d, want %d", category.ID, 10)
	}
}

// TestUpdateCategory verifies updating a category.
func TestUpdateCategory(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v3/1/categories/5" {
			t.Errorf("path = %q, want /api/v3/1/categories/5", r.URL.Path)
		}

		category := Category{ID: 5, Name: "Updated Category"}
		json.NewEncoder(w).Encode(category)
	})
	defer server.Close()

	category, err := client.UpdateCategory(5, &CategoryUpdateRequest{Name: "Updated Category"})
	if err != nil {
		t.Fatalf("UpdateCategory() error = %v", err)
	}

	if category.Name != "Updated Category" {
		t.Errorf("Name = %q, want %q", category.Name, "Updated Category")
	}
}

// TestDeleteCategory verifies deleting a category.
func TestDeleteCategory(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/categories/5" {
			t.Errorf("path = %q, want /api/v3/1/categories/5", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteCategory(5)
	if err != nil {
		t.Fatalf("DeleteCategory() error = %v", err)
	}
}

// --- Label Tests ---

// TestGetLabels verifies fetching a list of labels.
func TestGetLabels(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/labels" {
			t.Errorf("path = %q, want /api/v3/1/labels", r.URL.Path)
		}

		labels := []Label{
			{ID: 1, Name: "Tax Deductible"},
			{ID: 2, Name: "Q1 2026"},
		}
		json.NewEncoder(w).Encode(labels)
	})
	defer server.Close()

	labels, err := client.GetLabels(nil)
	if err != nil {
		t.Fatalf("GetLabels() error = %v", err)
	}

	if len(labels) != 2 {
		t.Fatalf("count = %d, want %d", len(labels), 2)
	}
}

// TestGetLabel verifies fetching a single label.
func TestGetLabel(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/labels/3" {
			t.Errorf("path = %q, want /api/v3/1/labels/3", r.URL.Path)
		}

		label := Label{ID: 3, Name: "Recurring"}
		json.NewEncoder(w).Encode(label)
	})
	defer server.Close()

	label, err := client.GetLabel(3)
	if err != nil {
		t.Fatalf("GetLabel() error = %v", err)
	}

	if label.Name != "Recurring" {
		t.Errorf("Name = %q, want %q", label.Name, "Recurring")
	}
}

// TestCreateLabel verifies creating a new label.
func TestCreateLabel(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req LabelCreateRequest
		json.Unmarshal(body, &req)

		label := Label{ID: 7, Name: req.Name}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(label)
	})
	defer server.Close()

	label, err := client.CreateLabel(&LabelCreateRequest{Name: "New Label"})
	if err != nil {
		t.Fatalf("CreateLabel() error = %v", err)
	}

	if label.ID != 7 {
		t.Errorf("ID = %d, want %d", label.ID, 7)
	}
}

// TestUpdateLabel verifies updating a label.
func TestUpdateLabel(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v3/1/labels/3" {
			t.Errorf("path = %q, want /api/v3/1/labels/3", r.URL.Path)
		}

		label := Label{ID: 3, Name: "Updated Label"}
		json.NewEncoder(w).Encode(label)
	})
	defer server.Close()

	label, err := client.UpdateLabel(3, &LabelUpdateRequest{Name: "Updated Label"})
	if err != nil {
		t.Fatalf("UpdateLabel() error = %v", err)
	}

	if label.Name != "Updated Label" {
		t.Errorf("Name = %q, want %q", label.Name, "Updated Label")
	}
}

// TestDeleteLabel verifies deleting a label.
func TestDeleteLabel(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/labels/3" {
			t.Errorf("path = %q, want /api/v3/1/labels/3", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteLabel(3)
	if err != nil {
		t.Fatalf("DeleteLabel() error = %v", err)
	}
}

// --- Contact Tests ---

// TestGetContacts verifies fetching a list of contacts.
func TestGetContacts(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/contacts" {
			t.Errorf("path = %q, want /api/v3/1/contacts", r.URL.Path)
		}

		contacts := []Contact{
			{ID: 1, Name: "Acme Corp"},
			{ID: 2, Name: "Widget Inc"},
		}
		json.NewEncoder(w).Encode(contacts)
	})
	defer server.Close()

	contacts, err := client.GetContacts(nil)
	if err != nil {
		t.Fatalf("GetContacts() error = %v", err)
	}

	if len(contacts) != 2 {
		t.Fatalf("count = %d, want %d", len(contacts), 2)
	}
}

// TestGetContact verifies fetching a single contact.
func TestGetContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/contacts/10" {
			t.Errorf("path = %q, want /api/v3/1/contacts/10", r.URL.Path)
		}

		contact := Contact{ID: 10, Name: "Acme Corp", Email: "billing@acme.com"}
		json.NewEncoder(w).Encode(contact)
	})
	defer server.Close()

	contact, err := client.GetContact(10)
	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact.Email != "billing@acme.com" {
		t.Errorf("Email = %q, want %q", contact.Email, "billing@acme.com")
	}
}

// TestCreateContact verifies creating a new contact.
func TestCreateContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req ContactCreateRequest
		json.Unmarshal(body, &req)

		if req.Name != "New Vendor" {
			t.Errorf("Name = %q, want %q", req.Name, "New Vendor")
		}

		contact := Contact{ID: 20, Name: req.Name}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contact)
	})
	defer server.Close()

	contact, err := client.CreateContact(&ContactCreateRequest{Name: "New Vendor"})
	if err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	if contact.ID != 20 {
		t.Errorf("ID = %d, want %d", contact.ID, 20)
	}
}

// TestUpdateContact verifies updating a contact.
func TestUpdateContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v3/1/contacts/10" {
			t.Errorf("path = %q, want /api/v3/1/contacts/10", r.URL.Path)
		}

		contact := Contact{ID: 10, Name: "Updated Vendor"}
		json.NewEncoder(w).Encode(contact)
	})
	defer server.Close()

	contact, err := client.UpdateContact(10, &ContactUpdateRequest{Name: "Updated Vendor"})
	if err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}

	if contact.Name != "Updated Vendor" {
		t.Errorf("Name = %q, want %q", contact.Name, "Updated Vendor")
	}
}

// TestDeleteContact verifies deleting a contact.
func TestDeleteContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/contacts/10" {
			t.Errorf("path = %q, want /api/v3/1/contacts/10", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteContact(10)
	if err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}
}

// --- File Tests ---

// TestUploadFile verifies uploading a file with multipart form data.
func TestUploadFile(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v3/1/files" {
			t.Errorf("path = %q, want /api/v3/1/files", r.URL.Path)
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want multipart/form-data", contentType)
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}

		_, _, err = r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile('file') error = %v", err)
		}

		if r.FormValue("ledger_id") != "42" {
			t.Errorf("ledger_id = %q, want 42", r.FormValue("ledger_id"))
		}

		file := File{ID: 1, Name: "receipt.jpg", Type: "image/jpeg"}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(file)
	})
	defer server.Close()

	// Create a temp file to upload.
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "receipt.jpg")
	os.WriteFile(tmpFile, []byte("fake image data"), 0644)

	file, err := client.UploadFile(tmpFile, "42")
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	if file.Name != "receipt.jpg" {
		t.Errorf("Name = %q, want %q", file.Name, "receipt.jpg")
	}
}

// TestUploadFileNoLedger verifies uploading a file without a ledger association.
func TestUploadFileNoLedger(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}

		if r.FormValue("ledger_id") != "" {
			t.Errorf("ledger_id should be empty, got %q", r.FormValue("ledger_id"))
		}

		file := File{ID: 2, Name: "doc.pdf"}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(file)
	})
	defer server.Close()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "doc.pdf")
	os.WriteFile(tmpFile, []byte("fake pdf data"), 0644)

	file, err := client.UploadFile(tmpFile, "")
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	if file.ID != 2 {
		t.Errorf("ID = %d, want %d", file.ID, 2)
	}
}

// --- Activity Tests ---

// TestGetActivities verifies fetching account activities.
func TestGetActivities(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/activities" {
			t.Errorf("path = %q, want /api/v3/1/activities", r.URL.Path)
		}

		activities := []Activity{
			{ID: 1, Action: "create", Message: "Created ledger entry"},
			{ID: 2, Action: "update", Message: "Updated category"},
		}
		json.NewEncoder(w).Encode(activities)
	})
	defer server.Close()

	activities, err := client.GetActivities(nil)
	if err != nil {
		t.Fatalf("GetActivities() error = %v", err)
	}

	if len(activities) != 2 {
		t.Fatalf("count = %d, want %d", len(activities), 2)
	}
	if activities[0].Action != "create" {
		t.Errorf("Action = %q, want %q", activities[0].Action, "create")
	}
}

// --- Report Tests ---

// TestGetPnlReport verifies fetching a P&L report.
func TestGetPnlReport(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/pnl" {
			t.Errorf("path = %q, want /api/v3/1/reports/pnl", r.URL.Path)
		}

		report := PnlReport{Income: 50000, Expense: 30000, Profit: 20000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetPnlReport(nil)
	if err != nil {
		t.Fatalf("GetPnlReport() error = %v", err)
	}

	if report.Profit != 20000 {
		t.Errorf("Profit = %f, want %f", report.Profit, 20000.0)
	}
}

// TestGetPnlByLabel verifies fetching P&L by label.
func TestGetPnlByLabel(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/pnl/label" {
			t.Errorf("path = %q, want /api/v3/1/reports/pnl/label", r.URL.Path)
		}

		report := PnlReport{Income: 10000, Expense: 5000, Profit: 5000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetPnlByLabel(nil)
	if err != nil {
		t.Fatalf("GetPnlByLabel() error = %v", err)
	}

	if report.Profit != 5000 {
		t.Errorf("Profit = %f, want %f", report.Profit, 5000.0)
	}
}

// TestGetPnlByCategory verifies fetching P&L by category.
func TestGetPnlByCategory(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/pnl/category" {
			t.Errorf("path = %q, want /api/v3/1/reports/pnl/category", r.URL.Path)
		}

		report := PnlReport{Income: 8000, Expense: 4000, Profit: 4000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetPnlByCategory(nil)
	if err != nil {
		t.Fatalf("GetPnlByCategory() error = %v", err)
	}

	if report.Profit != 4000 {
		t.Errorf("Profit = %f, want %f", report.Profit, 4000.0)
	}
}

// TestGetPnlCurrent verifies fetching the current year P&L.
func TestGetPnlCurrent(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/pnl/current" {
			t.Errorf("path = %q, want /api/v3/1/reports/pnl/current", r.URL.Path)
		}

		report := PnlReport{Income: 12000, Expense: 9000, Profit: 3000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetPnlCurrent()
	if err != nil {
		t.Fatalf("GetPnlCurrent() error = %v", err)
	}

	if report.Profit != 3000 {
		t.Errorf("Profit = %f, want %f", report.Profit, 3000.0)
	}
}

// TestGetIncomeByContact verifies fetching income by contact.
func TestGetIncomeByContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/income/by-contact" {
			t.Errorf("path = %q, want /api/v3/1/reports/income/by-contact", r.URL.Path)
		}

		report := PnlReport{
			Income: 15000,
			Breakdown: []PnlBreakdown{
				{Name: "Acme Corp", Amount: 10000},
				{Name: "Widget Inc", Amount: 5000},
			},
		}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetIncomeByContact(nil)
	if err != nil {
		t.Fatalf("GetIncomeByContact() error = %v", err)
	}

	if len(report.Breakdown) != 2 {
		t.Fatalf("Breakdown count = %d, want %d", len(report.Breakdown), 2)
	}
}

// TestGetExpensesByContact verifies fetching expenses by contact.
func TestGetExpensesByContact(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/reports/expenses/by-contact" {
			t.Errorf("path = %q, want /api/v3/1/reports/expenses/by-contact", r.URL.Path)
		}

		report := PnlReport{Expense: 8000}
		json.NewEncoder(w).Encode(report)
	})
	defer server.Close()

	report, err := client.GetExpensesByContact(nil)
	if err != nil {
		t.Fatalf("GetExpensesByContact() error = %v", err)
	}

	if report.Expense != 8000 {
		t.Errorf("Expense = %f, want %f", report.Expense, 8000.0)
	}
}

// --- Me Tests ---

// TestGetMe verifies fetching the current user's profile.
func TestGetMe(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/me" {
			t.Errorf("path = %q, want /api/v3/1/me", r.URL.Path)
		}

		me := MeResponse{ID: 42, FirstName: "John", LastName: "Doe", Email: "john@example.com"}
		json.NewEncoder(w).Encode(me)
	})
	defer server.Close()

	me, err := client.GetMe()
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}

	if me.Email != "john@example.com" {
		t.Errorf("Email = %q, want %q", me.Email, "john@example.com")
	}
}

// TestUpdateMe verifies updating the user profile.
func TestUpdateMe(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %q, want PUT", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req MeUpdateRequest
		json.Unmarshal(body, &req)

		if req.FirstName != "Jane" {
			t.Errorf("FirstName = %q, want %q", req.FirstName, "Jane")
		}

		me := MeResponse{ID: 42, FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}
		json.NewEncoder(w).Encode(me)
	})
	defer server.Close()

	me, err := client.UpdateMe(&MeUpdateRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	})
	if err != nil {
		t.Fatalf("UpdateMe() error = %v", err)
	}

	if me.FirstName != "Jane" {
		t.Errorf("FirstName = %q, want %q", me.FirstName, "Jane")
	}
}

// TestChangePassword verifies changing the user password.
func TestChangePassword(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v3/1/me/change-password" {
			t.Errorf("path = %q, want /api/v3/1/me/change-password", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req ChangePasswordRequest
		json.Unmarshal(body, &req)

		if req.CurrentPassword != "oldpass" {
			t.Errorf("CurrentPassword = %q, want %q", req.CurrentPassword, "oldpass")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	err := client.ChangePassword(&ChangePasswordRequest{
		CurrentPassword: "oldpass",
		NewPassword:     "newpass",
	})
	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
}

// --- User Tests ---

// TestGetUsers verifies fetching users in an account.
func TestGetUsers(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/users" {
			t.Errorf("path = %q, want /api/v3/1/users", r.URL.Path)
		}

		users := []User{
			{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			{ID: 2, FirstName: "Jane", LastName: "Smith", Email: "jane@example.com"},
		}
		json.NewEncoder(w).Encode(users)
	})
	defer server.Close()

	users, err := client.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("count = %d, want %d", len(users), 2)
	}
}

// TestRemoveUser verifies removing a user from an account.
func TestRemoveUser(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/users/5" {
			t.Errorf("path = %q, want /api/v3/1/users/5", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.RemoveUser(5)
	if err != nil {
		t.Fatalf("RemoveUser() error = %v", err)
	}
}

// TestGetInvites verifies fetching pending invitations.
func TestGetInvites(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/1/users/invite" {
			t.Errorf("path = %q, want /api/v3/1/users/invite", r.URL.Path)
		}

		invites := []Invite{
			{ID: 1, Email: "new@example.com", FirstName: "New", LastName: "User"},
		}
		json.NewEncoder(w).Encode(invites)
	})
	defer server.Close()

	invites, err := client.GetInvites()
	if err != nil {
		t.Fatalf("GetInvites() error = %v", err)
	}

	if len(invites) != 1 {
		t.Fatalf("count = %d, want %d", len(invites), 1)
	}
}

// TestCreateInvite verifies sending an invitation.
func TestCreateInvite(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req InviteCreateRequest
		json.Unmarshal(body, &req)

		if req.Email != "invite@example.com" {
			t.Errorf("Email = %q, want %q", req.Email, "invite@example.com")
		}

		invite := Invite{ID: 5, Email: req.Email, FirstName: req.FirstName}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(invite)
	})
	defer server.Close()

	invite, err := client.CreateInvite(&InviteCreateRequest{
		Email:     "invite@example.com",
		FirstName: "Test",
		LastName:  "User",
	})
	if err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}

	if invite.ID != 5 {
		t.Errorf("ID = %d, want %d", invite.ID, 5)
	}
}

// TestCancelInvite verifies cancelling an invitation.
func TestCancelInvite(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v3/1/user-invite/5" {
			t.Errorf("path = %q, want /api/v3/1/user-invite/5", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.CancelInvite(5)
	if err != nil {
		t.Fatalf("CancelInvite() error = %v", err)
	}
}

// Suppress unused import warnings.
var _ = fmt.Sprintf
