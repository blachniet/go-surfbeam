# go-surfbeam [![Build Status](https://travis-ci.org/blachniet/go-surfbeam.svg?branch=master)](https://travis-ci.org/blachniet/go-surfbeam) [![GoDoc](https://godoc.org/github.com/blachniet/go-surfbeam?status.svg)](https://godoc.org/github.com/blachniet/go-surfbeam) [![Rager Releases](http://rager.io/badge.svg)](http://rager.io/projects/search?badge=1&query=https://github.com/blachniet/go-surfbeam)

This package provides easy access modem status data provided by ViaSat's
SurfBeam® 2 satellite modems' CGI api.

## Usage

```go
surfbeamClient := surfbeam.New(surfbeam.DefaultModemURI)
if s, err := surfbeamClient.ModemStatus(); err == nil{
  fmt.Printf("Status: %v", s.Status)
}
```

## Compatibility

Check out the [Compatibility](https://github.com/blachniet/go-surfbeam/wiki/Compatibility) page to see which versions of the SurfBeam® 2 hardware/software that this package has been tested against.

## License

Copyright (c) 2015 Brian Lachniet. See the [LICENSE](LICENSE) file for license rights and limitations (MIT).
