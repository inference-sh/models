package models

import "encoding/json"

// UnmarshalJSON implements custom unmarshaling for FlowRunInput.
// Values can be plain scalars, arrays, objects, or connection descriptors.
func (f *FlowRunInput) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a connection
	var tryConn struct {
		Connection *FlowNodeConnection `json:"connection"`
		Value      json.RawMessage     `json:"value"`
	}
	if err := json.Unmarshal(data, &tryConn); err == nil && tryConn.Connection != nil {
		f.Connection = tryConn.Connection
		if len(tryConn.Value) > 0 {
			var nested FlowRunInput
			if err := json.Unmarshal(tryConn.Value, &nested); err == nil {
				f.Value = nested
				return nil
			}
			var generic any
			if err := flowUnmarshalRecursive(tryConn.Value, &generic); err == nil {
				f.Value = generic
				return nil
			}
		}
		return nil
	}

	var generic any
	if err := flowUnmarshalRecursive(data, &generic); err != nil {
		return err
	}
	f.Value = generic
	return nil
}

// MarshalJSON implements custom marshaling for FlowRunInput.
func (f FlowRunInput) MarshalJSON() ([]byte, error) {
	if f.Connection != nil {
		value, err := flowMarshalRecursive(f.Value)
		if err != nil {
			return nil, err
		}
		return json.Marshal(struct {
			Connection *FlowNodeConnection `json:"connection"`
			Value      json.RawMessage     `json:"value,omitempty"`
		}{f.Connection, value})
	}
	return flowMarshalRecursive(f.Value)
}

func flowUnmarshalRecursive(data []byte, out *any) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*out = flowProcessRaw(raw)
	return nil
}

func flowProcessRaw(raw any) any {
	switch v := raw.(type) {
	case []any:
		out := make([]any, len(v))
		for i, elem := range v {
			out[i] = flowProcessRaw(elem)
		}
		return out
	case map[string]any:
		if _, ok := v["connection"]; ok {
			b, _ := json.Marshal(v)
			var nested FlowRunInput
			if err := json.Unmarshal(b, &nested); err == nil {
				return nested
			}
		}
		out := make(map[string]any)
		for key, val := range v {
			out[key] = flowProcessRaw(val)
		}
		return out
	default:
		return v
	}
}

func flowMarshalRecursive(value any) ([]byte, error) {
	switch v := value.(type) {
	case FlowRunInput:
		return json.Marshal(v)
	case *FlowRunInput:
		return json.Marshal(v)
	case []any:
		arr := make([]any, len(v))
		for i, elem := range v {
			marshaled, err := flowMarshalRecursive(elem)
			if err != nil {
				return nil, err
			}
			var unmarshaled any
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				return nil, err
			}
			arr[i] = unmarshaled
		}
		return json.Marshal(arr)
	case map[string]any:
		mapped := make(map[string]any)
		for key, val := range v {
			marshaled, err := flowMarshalRecursive(val)
			if err != nil {
				return nil, err
			}
			var unmarshaled any
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				return nil, err
			}
			mapped[key] = unmarshaled
		}
		return json.Marshal(mapped)
	default:
		return json.Marshal(v)
	}
}
