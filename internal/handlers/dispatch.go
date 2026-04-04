// Package handlers implements the Bitbucket API dispatch layer for each CLI command.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// Request holds all parameters for a Dispatch call.
type Request struct {
	Method      string
	URLTemplate string
	PathParams  map[string]string
	QueryParams map[string]string
	Body        string
	All         bool
}

// pageResult holds extracted pagination data from a response.
type pageResult struct {
	values  []any
	nextURL string
}

// Dispatch performs a generic Bitbucket API request.
//
// It substitutes {param} placeholders in the URL template with path params,
// sets query parameters, sends body for POST/PUT/PATCH, and handles
// Bitbucket's cursor-based pagination when Request.All is true.
func Dispatch(ctx context.Context, c *client.BBClient, r Request) error {
	result, err := DispatchRaw(ctx, c, r)
	if err != nil {
		return err
	}
	if result == nil {
		fmt.Println("OK")
		return nil
	}
	return output.Render(result)
}

// DispatchRaw performs a generic Bitbucket API request and returns the raw
// parsed response (either a single object or a collected []any for paginated
// responses). Returns nil for empty/non-JSON responses. This is used by both
// the CLI (via Dispatch) and the MCP server.
func DispatchRaw(ctx context.Context, c *client.BBClient, r Request) (any, error) {
	baseURL := buildURL(r.URLTemplate, r.PathParams)
	var allValues []any
	fetchURL := baseURL

	for {
		resp, err := executeRequest(ctx, c, r, fetchURL, baseURL)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", r.Method, fetchURL, err)
		}
		if resp.IsError() {
			return nil, client.ParseError(resp)
		}

		result, nonJSON, err := parseResponseRaw(resp)
		if err != nil {
			return nil, err
		}
		if nonJSON {
			return nil, nil
		}
		if result == nil {
			return nil, nil
		}

		if page := extractPage(result); page != nil {
			allValues = append(allValues, page.values...)
			if r.All && page.nextURL != "" {
				fetchURL = page.nextURL
				continue
			}
			return allValues, nil
		}

		return result, nil
	}
}

// buildURL substitutes {param} placeholders in a URL template.
func buildURL(template string, pathParams map[string]string) string {
	url := template
	for k, v := range pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", v)
	}
	return url
}

// executeRequest builds and executes a single HTTP request.
// Query params are only set on the first request; pagination URLs already contain them.
func executeRequest(ctx context.Context, c *client.BBClient, r Request, fetchURL, baseURL string) (*resty.Response, error) {
	req := c.R().SetContext(ctx)

	if fetchURL == baseURL {
		for k, v := range r.QueryParams {
			if v != "" && v != "0" && v != "false" {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	if r.Body != "" && (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
		req = req.SetHeader("Content-Type", "application/json").SetBody(r.Body)
	}

	return req.Execute(r.Method, fetchURL)
}

// parseResponseRaw handles empty and non-JSON responses.
// Returns the parsed result, or nonJSON=true if the response is not JSON.
func parseResponseRaw(resp *resty.Response) (result any, nonJSON bool, err error) {
	respBody := resp.Body()
	if len(respBody) == 0 {
		return nil, false, nil
	}

	ct := resp.Header().Get("Content-Type")
	if !strings.Contains(ct, "json") {
		return string(respBody), true, nil
	}

	var parsed any
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, false, fmt.Errorf("parsing response: %w", err)
	}
	return parsed, false, nil
}

// extractPage checks whether result is a paginated Bitbucket response
// with a {"values": [...], "next": "..."} shape.
func extractPage(result any) *pageResult {
	m, ok := result.(map[string]any)
	if !ok {
		return nil
	}
	values, ok := m["values"]
	if !ok {
		return nil
	}
	arr, ok := values.([]any)
	if !ok {
		return nil
	}
	page := &pageResult{values: arr}
	if next, ok := m["next"].(string); ok && next != "" {
		page.nextURL = next
	}
	return page
}

// SetNested sets a value in a nested map using a dot-separated path.
// E.g., SetNested(m, "content.raw", "hello") produces {"content": {"raw": "hello"}}.
func SetNested(m map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	current := m
	for i, p := range parts {
		if i == len(parts)-1 {
			current[p] = value
		} else {
			if sub, ok := current[p]; ok {
				current = sub.(map[string]any)
			} else {
				sub := map[string]any{}
				current[p] = sub
				current = sub
			}
		}
	}
}

// GetNested retrieves a value from a nested map using a dot-separated path.
// Returns the value and true if found, or nil and false otherwise.
func GetNested(m map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := any(m)
	for _, p := range parts {
		cm, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = cm[p]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
