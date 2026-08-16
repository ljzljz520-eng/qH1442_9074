package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aftercare/internal/aftercare"
	"aftercare/internal/api"
)

func TestTicketBusinessWorkflow(t *testing.T) {
	handler := newHandler()
	created := submit(t, handler, map[string]any{
		"description": "  商品   外包装\n破损  ",
	}, http.StatusCreated)

	if created.State != aftercare.StateResult || created.Loading || created.Error != nil || created.Result == nil {
		t.Fatalf("create state = %#v, want result", created)
	}
	if created.Result.ID != "T-000001" {
		t.Fatalf("ticket id = %q, want T-000001", created.Result.ID)
	}
	if created.Result.Description != "商品 外包装 破损" {
		t.Fatalf("description = %q, want %q", created.Result.Description, "商品 外包装 破损")
	}
	if created.Result.CharacterCount != 9 {
		t.Fatalf("character count = %d, want 9", created.Result.CharacterCount)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/tickets/"+created.Result.ID, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("find status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var found aftercare.View
	decode(t, recorder, &found)
	if found.State != aftercare.StateResult || found.Result == nil || *found.Result != *created.Result {
		t.Fatalf("found ticket = %#v, want %#v", found.Result, created.Result)
	}
}

func TestTicketRequestValidation(t *testing.T) {
	handler := newHandler()

	tests := []struct {
		name string
		body map[string]any
		code string
	}{
		{name: "missing field", body: map[string]any{}, code: "description_required"},
		{name: "blank field", body: map[string]any{"description": " \t\n "}, code: "description_required"},
		{name: "unknown field", body: map[string]any{"description": "无法开机", "priority": "high"}, code: "invalid_request"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view := submit(t, handler, test.body, expectedValidationStatus(test.code))
			if view.State != aftercare.StateError || view.Loading || view.Error == nil || view.Result != nil {
				t.Fatalf("state = %#v, want error", view)
			}
			if view.Error.Code != test.code {
				t.Fatalf("error code = %q, want %q", view.Error.Code, test.code)
			}
		})
	}
}

func TestMultilingualTicketSubmission(t *testing.T) {
	handler := newHandler()
	description := strings.Repeat("售", 100)
	view := submit(t, handler, map[string]any{"description": description}, http.StatusCreated)

	if view.Result == nil {
		t.Fatal("result is nil")
	}
	if view.Result.Description != description {
		t.Fatalf("description contains %d characters, want %d", runeCount(view.Result.Description), runeCount(description))
	}
	if view.Result.CharacterCount != 100 {
		t.Fatalf("character count = %d, want 100", view.Result.CharacterCount)
	}
}

func newHandler() http.Handler {
	repository := aftercare.NewMemoryRepository()
	service := aftercare.NewService(repository)
	return api.NewHandler(service)
}

func submit(t *testing.T, handler http.Handler, body map[string]any, wantStatus int) aftercare.View {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, wantStatus, recorder.Body.String())
	}
	var view aftercare.View
	decode(t, recorder, &view)
	return view
}

func decode(t *testing.T, recorder *httptest.ResponseRecorder, value any) {
	t.Helper()
	if err := json.NewDecoder(recorder.Body).Decode(value); err != nil {
		t.Fatal(err)
	}
}

func expectedValidationStatus(code string) int {
	if code == "invalid_request" {
		return http.StatusBadRequest
	}
	return http.StatusUnprocessableEntity
}

func runeCount(value string) int {
	return len([]rune(value))
}
