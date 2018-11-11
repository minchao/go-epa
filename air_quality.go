package epa

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ResourceIDAirQualityForecast = "355000000I-000001"
)

type AirQualityForecastResponse struct {
	Success bool                     `json:"success"`
	Result  AirQualityForecastResult `json:"result"`
}

type AirQualityForecastResult struct {
	ResourceId string               `json:"resource_id"`
	Fields     []Fields             `json:"fields"`
	Records    []AirQualityForecast `json:"records"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	Total      int                  `json:"total"`
}

type AirQualityForecast struct {
	Content           string `json:"Content"`
	Area              string `json:"Area"`
	MajorPollutant    string `json:"MajorPollutant"`
	AQI               string `json:"AQI"`
	ForecastDate      string `json:"ForecastDate"`
	MinorPollutant    string `json:"MinorPollutant"`
	MinorPollutantAQI string `json:"MinorPollutantAQI"`
	PublishTime       string `json:"PublishTime"`
}

func (c *Client) GetAirQualityForecast(ctx context.Context, options url.Values) (*AirQualityForecastResponse, *http.Response, error) {
	u, _ := url.Parse(fmt.Sprintf("webapi/api/rest/datastore/%v", ResourceIDAirQualityForecast))
	u.RawQuery = options.Encode()

	forecast := new(AirQualityForecastResponse)
	req, err := c.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Do(ctx, req, forecast)
	if err != nil {
		return nil, resp, err
	}

	return forecast, resp, nil
}
