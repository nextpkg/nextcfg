# URL Source

The URL source reads config from a url.

It uses the url suffix as the format e.g `my.yaml` becomes `yaml`.
The content itself is not touched. If we can't find a format we'll use the encoder format.

## New Source

Specify url source with url. Defaults to `http://localhost:8080/config.yaml`.

```go
urlSource := url.NewSource(
url.WithURL("http://api.example.com/config.yaml"),
)
```

## Load Source

Load the source into config

```go
// Create new config
conf := nextcfg.NewConfig()

// Load url source
conf.Load(urlSource)
```

## 轮询时间

默认30S
