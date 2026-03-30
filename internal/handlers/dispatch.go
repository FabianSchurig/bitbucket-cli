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
	baseURL := buildURL(r.URLTemplate, r.PathParams)
	var allValues []any
	fetchURL := baseURL

	for {
		resp, err := executeRequest(ctx, c, r, fetchURL, baseURL)
		if err != nil {
			return fmt.Errorf("%s %s: %w", r.Method, fetchURL, err)
		}
		if resp.IsError() {
			return client.ParseError(resp)
		}

		result, done, err := parseResponse(resp)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		if page := extractPage(result); page != nil {
			allValues = append(allValues, page.values...)
			if r.All && page.nextURL != "" {
				fetchURL = page.nextURL
				continue
			}
			return output.Render(allValues)
		}

		return output.Render(result)
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

// parseResponse handles empty and non-JSON responses.
// Returns the parsed result, or done=true if output was already written.
func parseResponse(resp *resty.Response) (result any, done bool, err error) {
	respBody := resp.Body()
	if len(respBody) == 0 {
		fmt.Println("OK")
		return nil, true, nil
	}

	ct := resp.Header().Get("Content-Type")
	if !strings.Contains(ct, "json") {
		fmt.Print(string(respBody))
		return nil, true, nil
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
