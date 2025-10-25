package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListLocations -
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := c.cache.Get(url); ok {
		locationsResp := RespShallowLocations{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return RespShallowLocations{}, err
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil) //builds the request
	if err != nil {
		return RespShallowLocations{}, err
	}

	resp, err := c.httpClient.Do(req) //performs the request
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RespShallowLocations{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	dat, err := io.ReadAll(resp.Body) //converts the response into a slice of bytes - needed to be done the long way for caching
	if err != nil {
		return RespShallowLocations{}, err
	}

	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(dat, &locationsResp) //unmarshals the slice of bytes via encoding/json into the initialized struct
	if err != nil {
		return RespShallowLocations{}, err
	}

	c.cache.Add(url, dat)

	return locationsResp, nil
}
