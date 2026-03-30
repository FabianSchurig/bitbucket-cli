// Package handlers implements the Bitbucket API dispatch layer for each CLI command.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// Dispatch performs a generic Bitbucket API request.
//
// It substitutes {param} placeholders in urlTemplate with pathParams,
// sets query parameters, sends body for POST/PUT/PATCH, and handles
// Bitbucket's cursor-based pagination when all is true.
func Dispatch(
	ctx context.Context,
	c *client.BBClient,
	method, urlTemplate string,
	pathParams, queryParams map[string]string,
	body string,
	all bool,
) error {
	// Build URL from template and path params.
	url := urlTemplate
	for k, v := range pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", v)
	}

	var allValues []any
	fetchURL := url

	for {
		req := c.R().SetContext(ctx)

		// Set query params (skip empty values) only on the first request;
		// subsequent pagination URLs are absolute and already contain params.
		if fetchURL == url {
			for k, v := range queryParams {
				if v != "" && v != "0" && v != "false" {
					req = req.SetQueryParam(k, v)
				}
			}
		}

		// Set body for methods that accept one.
		if body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
			req = req.SetHeader("Content-Type", "application/json").SetBody(body)
		}

		resp, err := req.Execute(method, fetchURL)
		if err != nil {
			return fmt.Errorf("%s %s: %w", method, fetchURL, err)
		}
		if resp.IsError() {
			return client.ParseError(resp)
		}

		// If the response is empty (e.g. 204 No Content), we're done.
		respBody := resp.Body()
		if len(respBody) == 0 {
			fmt.Println("OK")
			return nil
		}

		// If the response is not JSON (e.g. raw diff), print it as-is.
		ct := resp.Header().Get("Content-Type")
		if !strings.Contains(ct, "json") {
			fmt.Print(string(respBody))
			return nil
		}

		// Parse the response as generic JSON.
		var result any
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Errorf("parsing response: %w", err)
		}

		// Check for paginated response pattern: {"values": [...], "next": "..."}
		m, isMap := result.(map[string]any)
		if isMap {
			if values, ok := m["values"]; ok {
				if arr, ok := values.([]any); ok {
					allValues = append(allValues, arr...)

					if all {
						if next, ok := m["next"].(string); ok && next != "" {
							fetchURL = next
							continue
						}
					}

					// Paginated response — render collected values.
					return output.Render(allValues)
				}
			}
		}

		// Non-paginated response — render the whole result.
		return output.Render(result)
	}
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
