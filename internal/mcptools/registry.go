package mcptools

// AllToolGroups is the central registry of all generated tool groups.
// Each generated .gen.go file appends its ToolGroup via an init() function.
// The server entry point iterates this slice to filter and register tools
// at startup based on the runtime configuration.
var AllToolGroups []ToolGroup
