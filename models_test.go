package models

import (
	"testing"
)

func TestParseTaskStatus(t *testing.T) {
	tests := []struct {
		input string
		want  TaskStatus
		ok    bool
	}{
		{"completed", TaskStatusCompleted, true},
		{"Completed", TaskStatusCompleted, true},
		{"RUNNING", TaskStatusRunning, true},
		{"queued", TaskStatusQueued, true},
		{"failed", TaskStatusFailed, true},
		{"cancelled", TaskStatusCancelled, true},
		{"received", TaskStatusReceived, true},
		{"dispatched", TaskStatusDispatched, true},
		{"preparing", TaskStatusPreparing, true},
		{"serving", TaskStatusServing, true},
		{"setting_up", TaskStatusSettingUp, true},
		{"uploading", TaskStatusUploading, true},
		{"cancelling", TaskStatusCancelling, true},
		{"bogus", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseTaskStatus(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseTaskStatus(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseTaskStatus_RoundTrips(t *testing.T) {
	statuses := []TaskStatus{
		TaskStatusReceived, TaskStatusQueued, TaskStatusDispatched,
		TaskStatusPreparing, TaskStatusServing, TaskStatusSettingUp,
		TaskStatusRunning, TaskStatusCancelling, TaskStatusUploading,
		TaskStatusCompleted, TaskStatusFailed, TaskStatusCancelled,
	}
	for _, s := range statuses {
		parsed, ok := ParseTaskStatus(s.String())
		if !ok {
			t.Errorf("ParseTaskStatus(%q) returned ok=false", s.String())
		}
		if parsed != s {
			t.Errorf("ParseTaskStatus(%q) = %d, want %d", s.String(), parsed, s)
		}
	}
}
