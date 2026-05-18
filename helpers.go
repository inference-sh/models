package models

import (
	"encoding/json"
	"strings"
)

// ParseRef parses a reference string in format "namespace/name@shortVersionId:function"
// Returns (namespace, name, shortVersionId, function). If no @ present, shortID is empty.
// If no : present, function is empty. If no / present, namespace is empty.
// "@latest" is treated as unversioned (shortID = "").
func ParseRef(ref string) (namespace, name, shortID, function string) {
	fullRef := ref
	if idx := strings.LastIndex(ref, ":"); idx != -1 {
		atIdx := strings.LastIndex(ref, "@")
		if atIdx == -1 || idx > atIdx {
			fullRef = ref[:idx]
			function = ref[idx+1:]
		}
	}

	fullName := fullRef
	if idx := strings.LastIndex(fullRef, "@"); idx != -1 {
		fullName = fullRef[:idx]
		shortID = fullRef[idx+1:]
	}

	if shortID == "latest" {
		shortID = ""
	}

	if idx := strings.Index(fullName, "/"); idx != -1 {
		namespace = fullName[:idx]
		name = fullName[idx+1:]
	} else {
		name = fullName
	}
	return namespace, name, shortID, function
}

// FullName returns "namespace/name" from the parsed components.
func FullName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return namespace + "/" + name
}

// BuildRef builds a reference string in format "namespace/name@shortVersionId:function".
// Empty components are omitted.
func BuildRef(namespace, name, shortID, function string) string {
	ref := FullName(namespace, name)
	if shortID != "" {
		ref += "@" + shortID
	}
	if function != "" {
		ref += ":" + function
	}
	return ref
}

// MicrocentsToDollars converts microcents to dollars.
func MicrocentsToDollars(microcents int64) float64 {
	return float64(microcents) / 100_000_000
}

// taskStatusToString maps TaskStatus int values to their string representation.
var taskStatusToString = map[TaskStatus]string{
	TaskStatusUnknown:    "unknown",
	TaskStatusReceived:   "received",
	TaskStatusQueued:     "queued",
	TaskStatusDispatched: "dispatched",
	TaskStatusPreparing:  "preparing",
	TaskStatusServing:    "serving",
	TaskStatusSettingUp:  "setting_up",
	TaskStatusRunning:    "running",
	TaskStatusCancelling: "cancelling",
	TaskStatusUploading:  "uploading",
	TaskStatusCompleted:  "completed",
	TaskStatusFailed:     "failed",
	TaskStatusCancelled:  "cancelled",
}

// String returns the human-readable name of a TaskStatus.
func (ts TaskStatus) String() string {
	if s, ok := taskStatusToString[ts]; ok {
		return s
	}
	return "unknown"
}

// IsTerminal returns true if the task is in a final state.
func (ts TaskStatus) IsTerminal() bool {
	return ts == TaskStatusCompleted || ts == TaskStatusFailed || ts == TaskStatusCancelled
}

// AffinityKey represents the fingerprint for task-to-worker matching.
type AffinityKey struct {
	AppID     string  `json:"app_id"`
	VersionID string  `json:"version_id"`
	Setup     string  `json:"setup,omitempty"`
	SessionID *string `json:"session_id,omitempty"`
}

// AffinityKeyFromTaskDTO creates an AffinityKey from a TaskDTO.
func AffinityKeyFromTaskDTO(task *TaskDTO) AffinityKey {
	key := AffinityKey{
		AppID:     task.AppID,
		VersionID: task.AppVersionID,
		SessionID: task.SessionID,
	}
	if task.Setup != nil {
		key.Setup = string(*task.Setup)
	}
	return key
}

// HasSession returns true if this key has a session binding.
func (k AffinityKey) HasSession() bool {
	return k.SessionID != nil && *k.SessionID != ""
}

// SameAs checks if two keys represent the same affinity.
func (a AffinityKey) SameAs(b AffinityKey) bool {
	return a.AppID == b.AppID &&
		a.VersionID == b.VersionID &&
		a.Setup == b.Setup &&
		ptrEqual(a.SessionID, b.SessionID)
}

// MatchesApp checks if keys are for the same app.
func (a AffinityKey) MatchesApp(b AffinityKey) bool {
	return a.AppID != "" && a.AppID == b.AppID
}

// MatchesVersion checks if keys are for the same app+version.
func (a AffinityKey) MatchesVersion(b AffinityKey) bool {
	return a.MatchesApp(b) && a.VersionID != "" && a.VersionID == b.VersionID
}

// MatchesSetup checks if keys are for the same app+version+setup (hot container).
func (a AffinityKey) MatchesSetup(b AffinityKey) bool {
	return a.MatchesVersion(b) && a.Setup == b.Setup
}

// CanReuseContainer checks if a task can reuse a container with this key.
func (a AffinityKey) CanReuseContainer(b AffinityKey) bool {
	if !a.MatchesSetup(b) {
		return false
	}
	if a.SessionID != nil && *a.SessionID == "new" {
		return false
	}
	return ptrEqual(a.SessionID, b.SessionID)
}

// String returns a compact string representation of the key.
func (k AffinityKey) String() string {
	s := k.AppID
	if k.VersionID != "" {
		s += ":" + k.VersionID
	}
	if k.Setup != "" {
		setup := k.Setup
		if len(setup) > 20 {
			setup = setup[:20] + "..."
		}
		s += "@" + setup
	}
	if k.SessionID != nil && *k.SessionID != "" {
		s += "#" + *k.SessionID
	}
	return s
}

// ToWarmApp converts the key to a warm app string.
func (k AffinityKey) ToWarmApp() string {
	s := k.AppID
	if k.VersionID != "" {
		s += ":" + k.VersionID
	}
	if k.Setup != "" {
		s += "@" + k.Setup
	}
	return s
}

// WorkerStatuses provides named constants for worker status values.
var WorkerStatuses = struct {
	Reserved WorkerStatus
	Busy     WorkerStatus
	Idle     WorkerStatus
	Inactive WorkerStatus
}{
	Reserved: "reserved",
	Busy:     "busy",
	Idle:     "idle",
	Inactive: "inactive",
}

// DeepCopy returns a deep copy of the TaskDTO via JSON round-trip.
func (t *TaskDTO) DeepCopy() *TaskDTO {
	if t == nil {
		return nil
	}
	data, err := json.Marshal(t)
	if err != nil {
		return nil
	}
	var copy TaskDTO
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil
	}
	return &copy
}

func ptrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
