# MixPanel

[![GoDoc](https://godoc.org/vizzlo.com/mixpanel?status.png)](https://godoc.org/vizzlo.com/mixpanel)
[![Build Status](https://travis-ci.org/vizzlo/mixpanel.svg?branch=master)](https://travis-ci.org/vizzlo/mixpanel)
[![Coverage Status](https://codecov.io/gh/vizzlo/mixpanel/branch/master/graph/badge.svg)](https://codecov.io/gh/vizzlo/mixpanel)
[![Go Report Card](https://goreportcard.com/badge/vizzlo.com/mixpanel)](https://goreportcard.com/report/vizzlo.com/mixpanel)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

- Latest version: 1.1.0
- Supported Go versions: 1.5 … 1.14

This is an unofficial client for the Mixpanel [event tracking](https://developer.mixpanel.com/docs/http) and
[data export API](https://developer.mixpanel.com/docs/data-export-api)s.

Next to the Go standard library, this library has no dependencies. Versions are correctly tagged---it supports Go Modules.

# Quickstart

```golang
import "vizzlo.com/mixpanel"

…

// API Token is used to access the event tracking and user engagement API
mp := mixpanel.New(token)

err := mp.Track("abc123…", "My Event", map[string]interface{}{
   "property1": "value1",
   "property2": 2,
   "property3": true,
})

// API Secret is used to access the data export API
client := mixpanel.NewExport(apiSecret)

// Downloads all profiles that have been seen during the last hour
profiles, err := exp.ListProfiles(&mixpanel.ProfileQuery{
    LastSeenAfter:   time.Now().Add(-time.Hour),
})
```

For more info, see the [API reference](https://godoc.org/vizzlo.com/mixpanel) or check out the examples folder.
