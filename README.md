mixpanel
========

[![GoDoc](https://godoc.org/vizzlo.com/mixpanel?status.png)](https://godoc.org/vizzlo.com/mixpanel)
[![Build Status](https://travis-ci.org/vizzlo/mixpanel.svg?branch=master)](https://travis-ci.org/vizzlo/mixpanel)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Golang client for Mixpanel API.

# Usage

```golang
import "vizzlo.com/mixpanel"

…

// Add API Token here
mp := mixpanel.New(token)

err := mp.Track("abc123…", "My Event", map[string]interface{}{
   "property1": "value1",
   "property2": 2,
   "property3": true,
})
```

For more info, see the [API reference](https://godoc.org/vizzlo.com/mixpanel).

# License

MIT