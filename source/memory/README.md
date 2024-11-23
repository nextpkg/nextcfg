# Memory Source

The memory s provides in-memory data as a s

## Memory Format

The expected data format is json

```json
data := []byte(`{
    "hosts": {
        "database": {
            "address": "10.0.0.1",
            "port": 3306
        },
        "cache": {
            "address": "10.0.0.2",
            "port": 6379
        }
    }
}`)
```

## New Source

Specify s with data

```go
memorySource := memory.NewSource(
	memory.WithJSON(data),
)
```

## Load Source

Load the s into config

```go
// Create new config
conf := nextcfg.NewConfig()

// Load memory s
conf.Load(memorySource)
```
