# Sensu Go Control By Web Plugin 
TravisCI: [![TravisCI Build Status](https://travis-ci.com/jspaleta/sensu-controlbyweb-plugins.svg?branch=master)](https://travis-ci.com/jspaleta/sensu-controlbyweb-plugin)

Collection of Sensu Agent check commands to interact with Particle Cloud

## Installation

Download the latest version from [releases][1],
or create an executable script from this source.

### Build from source
From the local path of this repository:

#### WX110 Metric Check
```
go build -o /usr/local/bin/wx110_metric_check ./wx110_metric_check/main.go
```

## Configuration

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "check",
    "metadata": {
        "namespace": "default",
        "name": "wx110"
    },
    "spec": {
        "...": "..."
    }
}
```

## Usage Examples


## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/jspaleta/sensu-controlbyweb-plugins/releases
