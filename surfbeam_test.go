package surfbeam

import (
	"bytes"
	"net"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("http://192.168.100.1")
	if c == nil {
		t.Error("client was null")
		t.FailNow()
	}

	expected := "http://192.168.100.1"
	if c.modemURI != expected {
		t.Errorf("Expected modemUri %v but was %v", expected, c.modemURI)
	}

	if c.client != http.DefaultClient {
		t.Errorf("Expected http.DefaultClient but was something else")
	}
}

var uint32Tests = []struct {
	in  string
	out uint32
}{
	{"", 0},
	{" ", 0},
	{"12", 12},
	{"1,234,554", 1234554},
	{"1,23,44", 12344},
}

var uint64Tests = []struct {
	in  string
	out uint64
}{
	{"", 0},
	{" ", 0},
	{"12", 12},
	{"1,234,554", 1234554},
	{"1,23,44", 12344},
}

var floatTests = []struct {
	in  string
	out float64
}{
	{"", 0.},
	{"  ", 0.},
	{"12", 12.},
	{"12.345", 12.345},
	{"1,234,554.678", 1234554.678},
	{"1,22,3.677", 1223.677},
}

var percentageTests = []struct {
	in  string
	out float64
}{
	{"", 0.},
	{"  ", 0.},
	{"12", 12.},
	{"12.345", 12.345},
	{"1,234,554.678", 1234554.678},
	{"1,22,3.677", 1223.677},
	{"12%", 12.},
	{"12.345 %", 12.345},
	{"1,234,554.678  %", 1234554.678},
}

func TestParseUint32(t *testing.T) {
	var status ModemStatus
	for _, tt := range uint32Tests {
		if err := parseUint32(tt.in, &status.LossOfSyncCount); err != nil {
			t.Errorf("parseUint32 %q: Error parsing: %v", tt.in, err)
		} else {
			if status.LossOfSyncCount != tt.out {
				t.Errorf("parseUint32 %q: Expected %d but was %d", tt.in, tt.out, status.TransmittedPackets)
			}
		}
	}
}

func TestParseUint64(t *testing.T) {
	var status ModemStatus
	for _, tt := range uint64Tests {
		if err := parseUint64(tt.in, &status.TransmittedPackets); err != nil {
			t.Errorf("parseUint64 %q: Error parsing: %v", tt.in, err)
		} else {
			if status.TransmittedPackets != tt.out {
				t.Errorf("parseUint64 %q: Expected %d but was %d", tt.in, tt.out, status.TransmittedPackets)
			}
		}
	}
}

func TestParseFloat64(t *testing.T) {
	var status ModemStatus
	for _, tt := range floatTests {
		if err := parseFloat64(tt.in, &status.RxPower); err != nil {
			t.Errorf("parseFloat64 %q: Error parsing: %v", tt.in, err)
		} else {
			if status.RxPower != tt.out {
				t.Errorf("parseFloat64 %q: Expected %f but was %f", tt.in, tt.out, status.RxPower)
			}
		}
	}
}

func TestParsePercentage(t *testing.T) {
	var status ModemStatus
	for _, tt := range percentageTests {
		if err := parsePercentage(tt.in, &status.RxPowerPercentage); err != nil {
			t.Errorf("parsePercentage %q: Error parsing: %v", tt.in, err)
		} else {
			if status.RxPowerPercentage != tt.out {
				t.Errorf("parsePercentage %q: Expected %f but was %f", tt.in, tt.out, status.RxPowerPercentage)
			}
		}
	}
}

var modemStatusTests = []struct {
	in                 string
	IPAddress          net.IP
	MACAddress         net.HardwareAddr
	SoftwareVersion    string
	HardwareVersion    string
	Status             string
	TransmittedPackets uint64
	TransmittedBytes   uint64
	ReceivedPackets    uint64
	ReceivedBytes      uint64
}{
	{
		`55.55.55.55##55:55:55:55:55:55##UT_1.5.2.2.3##UT_7 P3_V1##Online##7,088,473##1,693,689,905##9,378,757##9,532,929,796##000:01:37:54##2##6.2##32%##012345678901##-44.3##47%##1.5##6%##Active##12.4##82%##Single##0123456789##images/Modem_Status_005_Online.png##/images/Satellite_Status_Purple.png##0##<p style=""color:green"">Connected</p>##<p style=""color:green"">Good</p>##0.00%##0.00%##6.67s##195##10000000##`,
		net.IPv4(0x37, 0x37, 0x37, 0x37),
		[]byte{0x37, 0x37, 0x37, 0x37, 0x37, 0x37},
		"UT_1.5.2.2.3",
		"UT_7 P3_V1",
		"Online",
		7088473,
		1693689905,
		9378757,
		9532929796,
	},
}

func TestParseModemStatus(t *testing.T) {
	for _, tt := range modemStatusTests {
		if status, err := parseModemStatus(tt.in); err != nil {
			t.Errorf("ParseModemStatus error: %v (in: %q)", err, tt.in)
		} else {
			assertIPEqual(t, tt.IPAddress, status.IPAddress, "IPAddress")
			assertMACEqual(t, tt.MACAddress, status.MACAddress, "MACAddress")
			assertEqual(t, tt.SoftwareVersion, status.SoftwareVersion, "SoftwareVersion")
			assertEqual(t, tt.HardwareVersion, status.HardwareVersion, "HardwareVersion")
			assertEqual(t, tt.Status, status.Status, "Status")
			assertEqual(t, tt.TransmittedPackets, status.TransmittedPackets, "TransmittedPackets")
			assertEqual(t, tt.TransmittedBytes, status.TransmittedBytes, "TransmittedBytes")
			assertEqual(t, tt.ReceivedPackets, status.ReceivedPackets, "ReceivedPackets")
			assertEqual(t, tt.ReceivedBytes, status.ReceivedBytes, "ReceivedBytes")
		}
	}
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}, name string) {
	if expected != actual {
		t.Errorf("Expected %v %v but was %v", name, expected, actual)
	}
}

func assertIPEqual(t *testing.T, expected net.IP, actual net.IP, name string) {
	if !expected.Equal(actual) {
		t.Errorf("Expected %v to be %v but was %v", name, expected, actual)
	}
}

func assertMACEqual(t *testing.T, expected, actual net.HardwareAddr, name string) {
	if bytes.Equal(expected, actual) {
		t.Errorf("Expected %v to be %v but was %v", name, expected, actual)
	}
}
