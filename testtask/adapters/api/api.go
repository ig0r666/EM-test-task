package api

import (
	"EMtask/testtask/config"
	"EMtask/testtask/core"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	log    *slog.Logger
	config config.Config
	client http.Client
}

func NewClient(config config.Config, log *slog.Logger) *Client {
	return &Client{
		log:    log,
		config: config,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetAge(name string) (core.APIAgeResponse, error) {
	url := fmt.Sprintf("https://%s/?name=%s", c.config.APIConfig.AgeAPIUrl, name)
	resp, err := c.client.Get(url)
	if err != nil {
		c.log.Error("GetAge API request failed", "error", err)
		return core.APIAgeResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.log.Error("GetAge API non-200 response", "status", resp.StatusCode)
		return core.APIAgeResponse{}, errors.New("api returned non-200 status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("GetAge failed to read response body", "error", err)
		return core.APIAgeResponse{}, err
	}

	var ageResp core.APIAgeResponse
	if err := json.Unmarshal(body, &ageResp); err != nil {
		c.log.Error("GetAge failed to unmarshal JSON", "error", err)
		return core.APIAgeResponse{}, err
	}

	return ageResp, nil
}

func (c *Client) GetGender(name string) (core.APIGenderResponse, error) {
	url := fmt.Sprintf("https://%s/?name=%s", c.config.APIConfig.GenderAPIUrl, name)
	resp, err := c.client.Get(url)
	if err != nil {
		c.log.Error("GetGender API request failed", "error", err)
		return core.APIGenderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.APIGenderResponse{}, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var genderResp core.APIGenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderResp); err != nil {
		c.log.Error("GetGender failed to decode JSON", "error", err)
		return core.APIGenderResponse{}, err
	}

	return genderResp, nil
}

func (c *Client) GetNation(name string) (core.APINationResponse, error) {
	url := fmt.Sprintf("https://%s/?name=%s", c.config.APIConfig.NationAPIUrl, name)
	resp, err := c.client.Get(url)
	if err != nil {
		c.log.Error("GetNation API request failed", "error", err)
		return core.APINationResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.APINationResponse{}, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var nationResp core.APINationResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationResp); err != nil {
		c.log.Error("GetNation failed to decode JSON", "error", err)
		return core.APINationResponse{}, err
	}

	return nationResp, nil
}
