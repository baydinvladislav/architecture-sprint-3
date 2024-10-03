package suppliers

import (
	"bytes"
	"fmt"
	"net/http"
)

type DeviceServiceSupplier struct {
	BaseURL string
}

func NewDeviceServiceSupplier(baseURL string) *DeviceServiceSupplier {
	return &DeviceServiceSupplier{BaseURL: baseURL}
}

func (ds *DeviceServiceSupplier) TurnOffModule(houseID string, moduleID string) error {
	url := fmt.Sprintf("%s/device/modules/houses/%s/modules/%s/turn-off", ds.BaseURL, houseID, moduleID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to turn off module: received status code %d", resp.StatusCode)
	}

	return nil
}
