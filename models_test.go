package models

import "testing"

func TestTaskStatus_Parse(t *testing.T) {
	tests := []struct {
		input string
		want  TaskStatus
		ok    bool
	}{
		{"completed", TaskStatusCompleted, true},
		{"running", TaskStatusRunning, true},
		{"queued", TaskStatusQueued, true},
		{"failed", TaskStatusFailed, true},
		{"cancelled", TaskStatusCancelled, true},
		{"bogus", TaskStatusUnknown, false},
		{"", TaskStatusUnknown, false},
	}
	for _, tt := range tests {
		got, ok := TaskStatus(0).Parse(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("TaskStatus.Parse(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestTaskStatus_ParseRoundTrips(t *testing.T) {
	statuses := []TaskStatus{
		TaskStatusReceived, TaskStatusQueued, TaskStatusDispatched,
		TaskStatusPreparing, TaskStatusServing, TaskStatusSettingUp,
		TaskStatusRunning, TaskStatusCancelling, TaskStatusUploading,
		TaskStatusCompleted, TaskStatusFailed, TaskStatusCancelled,
	}
	for _, s := range statuses {
		parsed, ok := TaskStatus(0).Parse(s.String())
		if !ok {
			t.Errorf("TaskStatus.Parse(%q) returned ok=false", s.String())
		}
		if parsed != s {
			t.Errorf("TaskStatus.Parse(%q) = %d, want %d", s.String(), parsed, s)
		}
	}
}
