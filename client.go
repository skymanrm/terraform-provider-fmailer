// client.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the FMailer API client
type Client struct {
	endpoint   string
	token      string
	httpClient *http.Client
}

// NewClient creates a new FMailer API client
func NewClient(endpoint string, token string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// DomainTemplate represents a domain template
type DomainTemplate struct {
	ID        int                  `json:"id,omitempty"`
	UUID      string               `json:"uuid,omitempty"`
	CreatedAt time.Time            `json:"created_at,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty"`
	Name      string               `json:"name"`
	Slug      string               `json:"slug"`
	AllowCopy bool                 `json:"allow_copy,omitempty"`
	Editable  bool                 `json:"editable,omitempty"`
	Domain    int                  `json:"domain"`
	Langs     []DomainTemplateLang `json:"langs,omitempty"`
}

// DomainTemplateLang represents a domain template language
type DomainTemplateLang struct {
	ID        int       `json:"id,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Lang      string    `json:"lang"`
	Default   bool      `json:"default,omitempty"`
	Template  int       `json:"template,omitempty"`
}

// PaginatedDomainTemplateList represents a paginated list of domain templates
type PaginatedDomainTemplateList struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []DomainTemplate `json:"results"`
}

// CreateDomainTemplate creates a new domain template
func (c *Client) CreateDomainTemplate(domainTemplate *DomainTemplate) (*DomainTemplate, error) {
	body, err := json.Marshal(domainTemplate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/domains/templates/", c.endpoint), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result DomainTemplate
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDomainTemplate gets a domain template by UUID
func (c *Client) GetDomainTemplate(uuid string) (*DomainTemplate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/domains/templates/%s/", c.endpoint, uuid), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result DomainTemplate
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDomainTemplate updates a domain template
func (c *Client) UpdateDomainTemplate(uuid string, domainTemplate *DomainTemplate) (*DomainTemplate, error) {
	body, err := json.Marshal(domainTemplate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/domains/templates/%s/", c.endpoint, uuid), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result DomainTemplate
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDomainTemplate deletes a domain template
func (c *Client) DeleteDomainTemplate(uuid string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/domains/templates/%s/", c.endpoint, uuid), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(responseBody))
	}

	return nil
}

// ListDomainTemplates lists domain templates with optional filters
func (c *Client) ListDomainTemplates(domain *int, search *string, page *int, ordering *string) (*PaginatedDomainTemplateList, error) {
	url := fmt.Sprintf("%s/api/domains/templates/", c.endpoint)

	// Add query parameters
	first := true
	if domain != nil {
		url += fmt.Sprintf("?domain=%d", *domain)
		first = false
	}

	if search != nil {
		if first {
			url += "?"
			first = false
		} else {
			url += "&"
		}
		url += fmt.Sprintf("search=%s", *search)
	}

	if page != nil {
		if first {
			url += "?"
			first = false
		} else {
			url += "&"
		}
		url += fmt.Sprintf("page=%d", *page)
	}

	if ordering != nil {
		if first {
			url += "?"
		} else {
			url += "&"
		}
		url += fmt.Sprintf("ordering=%s", *ordering)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result PaginatedDomainTemplateList
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DuplicateDomainTemplate duplicates a domain template
func (c *Client) DuplicateDomainTemplate(uuid string, name string, slug string) error {
	payload := struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}{
		Name: name,
		Slug: slug,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/domains/templates/%s/duplicate/", c.endpoint, uuid), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
