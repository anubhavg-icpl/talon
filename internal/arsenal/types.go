package arsenal

// paramSpec describes one tool parameter for both the MCP input schema (what
// the LLM sees) and the outbound request (what the Arsenal Engine expects).
type paramSpec struct {
	Name     string
	Kind     string // "string", "integer", "number", "boolean", "object", "array"
	Required bool
	Default  any
}

// bodyField maps a function param to the JSON key the Arsenal Engine
// expects it under -- usually identical to ParamName, occasionally renamed
// (e.g. generate_payload takes `payload_type` but sends it as JSON key
// "type").
type bodyField struct {
	ParamName string
	JSONKey   string
}

// toolSpec fully describes one Arsenal tool: enough to build both its MCP
// schema and its HTTP call.
type toolSpec struct {
	Name         string
	Description  string
	Params       []paramSpec
	HTTPMethod   string // "GET" or "POST"
	Endpoint     string // may contain fmt %v placeholders for path params
	EndpointVars []string
	BodyFields   []bodyField
	StaticFields map[string]any
}
