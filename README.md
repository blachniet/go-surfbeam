# SurfBeam

This package provides easy access modem status data provided by ViaSat's
SurfBeam® 2 satellite modems' CGI api.

## Usage

```go
surfbeamClient := surfbeam.New(surfbeam.DefaultModemURI)
if s, err := surfbeamClient.ModemStatus(); err == nil{
  fmt.Printf("Status: %v", s.Status)
}
```
