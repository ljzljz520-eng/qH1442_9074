package aftercare_test

import (
	"context"
	"testing"

	"aftercare/internal/aftercare"
)

func TestSubmissionStates(t *testing.T) {
	service := aftercare.NewService(aftercare.NewMemoryRepository())
	var observed []aftercare.View

	view := service.Submit(context.Background(), "  屏幕   出现闪烁  ", func(state aftercare.View) {
		observed = append(observed, state)
	})

	if len(observed) != 1 {
		t.Fatalf("observed %d intermediate states, want 1", len(observed))
	}
	if observed[0].State != aftercare.StateLoading || !observed[0].Loading {
		t.Fatalf("intermediate state = %#v, want loading", observed[0])
	}
	if view.State != aftercare.StateResult || view.Loading || view.Error != nil || view.Result == nil {
		t.Fatalf("final state = %#v, want result", view)
	}
	if view.Result.Description != "屏幕 出现闪烁" {
		t.Fatalf("description = %q, want %q", view.Result.Description, "屏幕 出现闪烁")
	}
	if view.Result.CharacterCount != 7 {
		t.Fatalf("character count = %d, want 7", view.Result.CharacterCount)
	}
}

func TestSubmissionValidation(t *testing.T) {
	service := aftercare.NewService(aftercare.NewMemoryRepository())
	view := service.Submit(context.Background(), " \n\t ", nil)

	if view.State != aftercare.StateError || view.Loading || view.Error == nil || view.Result != nil {
		t.Fatalf("state = %#v, want error", view)
	}
	if view.Error.Code != "description_required" {
		t.Fatalf("error code = %q, want description_required", view.Error.Code)
	}
}
