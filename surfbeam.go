package surfbeam

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// The default URI used to access the modem.
const DefaultModemURI string = "http://192.168.100.1"
const modemStatusPath string = "/index.cgi?page=modemStatusData"
const triaStatusPath string = "/index.cgi?page=triaStatusData"

// A Client is used to query the modem for status information.
type Client struct {
	modemURI string
	client   *http.Client
}

// New creates a new Client with the default http client.
func New(modemURI string) *Client {
	return &Client{modemURI, http.DefaultClient}
}

// NewWithClient creates a new Client with the given http client.
func NewWithClient(modemURI string, client *http.Client) *Client {
	return &Client{modemURI, client}
}

// ModemStatus retrieves information describing the state of the modem.
func (c *Client) ModemStatus() (*ModemStatus, error) {
	resp, err := c.client.Get(c.modemStatusURI())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response was not ok: %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseModemStatus(string(contents))
}

func (c *Client) modemStatusURI() string {
	return c.modemURI + modemStatusPath
}

func (c *Client) triaStatusURI() string {
	return c.modemURI + triaStatusPath
}

func parseModemStatus(str string) (*ModemStatus, error){
  parts := strings.Split(str, "##")
  var status ModemStatus
  status.IPAddress = parts[0]
  status.MACAddress = parts[1]
  status.SoftwareVersion = parts[2]
  status.HardwareVersion = parts[3]
  status.Status = parts[4]
  if err := parseInt64(parts[5], &status.TransmittedPackets); err != nil { return nil, err }
  if err := parseInt64(parts[6], &status.TransmittedBytes); err != nil { return nil, err }
  if err := parseInt64(parts[7], &status.ReceivedPackets); err != nil { return nil, err }
  if err := parseInt64(parts[8], &status.ReceivedBytes); err != nil { return nil, err }
  status.OnlineTime = parts[9]
  if err := parseInt64(parts[10], &status.LossOfSyncCount); err != nil{ return nil, err}
  if err := parseFloat64(parts[11], &status.RxSNR); err != nil { return nil, err}
  if err := parsePercentage(parts[12], &status.RxSNRPercentage); err != nil { return nil, err}
  status.SerialNumber = parts[13]
  if err := parseFloat64(parts[14], &status.RxPower); err != nil { return nil, err}
  if err := parsePercentage(parts[15], &status.RxPowerPercentage); err != nil { return nil, err}
  if err := parseFloat64(parts[16], &status.CableResistance); err != nil { return nil, err}
  if err := parsePercentage(parts[17], &status.CableResistancePercentage); err != nil { return nil, err}
  status.ODUTelemetryStatus = parts[18]
  if err := parseFloat64(parts[19], &status.CableAttenuation); err != nil { return nil, err}
  if err := parsePercentage(parts[20], &status.CableAttenuationPercentage); err != nil { return nil, err}
  status.IFLType = parts[21]
  status.PartNumber = parts[22]
  status.StatusImageURI = parts[23]
  status.SatelliteStatusURI = parts[24]
  status.Unknown25 = parts[25]
  status.StatusHTML = parts[26]
  status.HealthHTML = parts[27]
  status.Unknown28 = parts[28]
  status.Unknown29 = parts[29]
  status.LastPageLoadDuration = parts[30]
  status.Unknown31 = parts[31]
  status.Unknown32 = parts[32]
  return &status, nil
}

// Parses an integer from the given string. Can handle leading/trailing spaces
// as well as embedded commas.
func parseInt64(str string, dest *int64) error {
  str = strings.TrimSpace(str)
	if len(str) == 0 {
		*dest = 0
    return nil
	}

	str = strings.Replace(str, ",", "", -1) // Remove commas
	if val, err := strconv.ParseInt(str, 10, 64); err != nil{
    return err
  } else {
    *dest = val
    return nil
  }
}

// Parses a float from the given string. Can handle leading/trailing spaces as
// well as embedded commas.
func parseFloat64(str string, dest *float64) error{
  str = strings.TrimSpace(str)
  if len(str) == 0 {
    *dest = 0
    return nil
  }

  str = strings.Replace(str, ",", "", -1) // Remove commas
  if val, err := strconv.ParseFloat(str, 64); err != nil{
    return err
  } else {
    *dest = val
    return nil
  }
}

// Parses a float from the given string. This method works the same as
// parseFloat64 except that it also trims percentage signs (%).
func parsePercentage(str string, dest *float64) error{
  str = strings.Replace(str, "%", "", -1)
  return parseFloat64(str, dest)
}

// ModemStatus provides information describing the state of the modem.
type ModemStatus struct {
	IPAddress                  string
	MACAddress                 string
	SoftwareVersion            string
	HardwareVersion            string
	Status                     string
	TransmittedPackets         int64
	TransmittedBytes           int64
	ReceivedPackets            int64
	ReceivedBytes              int64
	OnlineTime                 string // TODO: use time.Duration
	LossOfSyncCount            int64
	RxSNR                      float64
	RxSNRPercentage            float64
	SerialNumber               string
	RxPower                    float64
	RxPowerPercentage          float64
	CableResistance            float64
	CableResistancePercentage  float64
	ODUTelemetryStatus         string
	CableAttenuation           float64
	CableAttenuationPercentage float64
	IFLType                    string
	PartNumber                 string
	StatusImageURI             string
	SatelliteStatusURI         string
	Unknown25                  string
	StatusHTML                 string
	HealthHTML                 string
	Unknown28                  string
	Unknown29                  string
	LastPageLoadDuration       string
	Unknown31                  string
	Unknown32                  string
}
