package models

import "encoding/json"

// Flow graph action type constants.
const (
	ActionNodeAdd        = "node.add"
	ActionNodeRemove     = "node.remove"
	ActionNodeMove       = "node.move"
	ActionNodeMoveMany   = "node.move_many"
	ActionNodeDuplicate  = "node.duplicate"
	ActionNodeRename     = "node.rename"
	ActionNodeSetApp     = "node.set_app"
	ActionNodeUpdate     = "node.update"
	ActionNodeSetInput   = "node.set_input"
	ActionNodeClearInput = "node.clear_input"

	ActionEdgeAdd    = "edge.add"
	ActionEdgeRemove = "edge.remove"

	ActionFlowSetInputSchema      = "flow.set_input_schema"
	ActionFlowSetOutputSchema     = "flow.set_output_schema"
	ActionFlowSetOutputMapping    = "flow.set_output_mapping"
	ActionFlowRemoveOutputMapping = "flow.remove_output_mapping"
	ActionFlowRenameOutputField   = "flow.rename_output_field"
)

// FlowAction represents a single graph mutation sent to POST /flows/{id}/actions.
type FlowAction struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// FlowActionsRequest is the request body for POST /flows/{id}/actions.
type FlowActionsRequest struct {
	Actions []FlowAction `json:"actions"`
}

// FlowActionsResponse is the response from the actions endpoint.
type FlowActionsResponse struct {
	Version int          `json:"version"`
	Actions []FlowAction `json:"actions"`
	Errors  []FlowActionError `json:"errors,omitempty"`
}

// FlowActionError is an error returned from an action.
type FlowActionError struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
}

// --- Action payloads ---

// AddNodePayload is the payload for node.add.
type AddNodePayload struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Position FlowNodePosition `json:"position"`
}

// RemoveNodePayload is the payload for node.remove.
type RemoveNodePayload struct {
	ID string `json:"id"`
}

// MoveNodePayload is the payload for node.move.
type MoveNodePayload struct {
	ID       string           `json:"id"`
	Position FlowNodePosition `json:"position"`
}

// MoveNodesPayload is the payload for node.move_many.
type MoveNodesPayload struct {
	Positions map[string]FlowNodePosition `json:"positions"`
}

// DuplicateNodePayload is the payload for node.duplicate.
type DuplicateNodePayload struct {
	SourceID string           `json:"source_id"`
	NewID    string           `json:"new_id"`
	Offset   FlowNodePosition `json:"offset"`
}

// RenameNodePayload is the payload for node.rename.
type RenameNodePayload struct {
	OldID string `json:"old_id"`
	NewID string `json:"new_id"`
}

// SetNodeAppPayload is the payload for node.set_app.
type SetNodeAppPayload struct {
	NodeID       string `json:"node_id"`
	AppID        string `json:"app_id"`
	AppVersionID string `json:"app_version_id"`
	Function     string `json:"function"`
}

// UpdateNodeDataPayload is the payload for node.update.
type UpdateNodeDataPayload struct {
	NodeID string         `json:"node_id"`
	Patch  map[string]any `json:"patch"`
}

// SetInputPayload is the payload for node.set_input.
type SetInputPayload struct {
	NodeID   string       `json:"node_id"`
	InputKey string       `json:"input_key"`
	Input    FlowRunInput `json:"input"`
}

// ClearInputPayload is the payload for node.clear_input.
type ClearInputPayload struct {
	NodeID   string `json:"node_id"`
	InputKey string `json:"input_key"`
}

// AddEdgePayload is the payload for edge.add.
type AddEdgePayload struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	SourceHandle *string `json:"source_handle,omitempty"`
	TargetHandle *string `json:"target_handle,omitempty"`
}

// RemoveEdgePayload is the payload for edge.remove.
type RemoveEdgePayload struct {
	ID string `json:"id"`
}

// SetSchemaPayload is the payload for flow.set_input_schema and flow.set_output_schema.
type SetSchemaPayload struct {
	Schema json.RawMessage `json:"schema"`
}

// SetOutputMappingPayload is the payload for flow.set_output_mapping.
type SetOutputMappingPayload struct {
	Field   string             `json:"field"`
	Mapping OutputFieldMapping `json:"mapping"`
}

// RemoveOutputMappingPayload is the payload for flow.remove_output_mapping.
type RemoveOutputMappingPayload struct {
	Field string `json:"field"`
}

// RenameOutputFieldPayload is the payload for flow.rename_output_field.
type RenameOutputFieldPayload struct {
	OldField string `json:"old_field"`
	NewField string `json:"new_field"`
}

// NewAction creates a FlowAction from a type and typed payload.
func NewAction(actionType string, payload any) FlowAction {
	data, _ := json.Marshal(payload)
	return FlowAction{Type: actionType, Payload: data}
}
