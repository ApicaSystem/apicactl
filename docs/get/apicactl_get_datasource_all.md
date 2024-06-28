# `apicactl get datasource all`

List all the available datasources

## Synopsis

List all the available datasources

```
apicactl get datasource all [flags]
```

## Examples

```
apicactl get datasource all
```

## Options

```
  -h, --help   help for all
```

## Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `apicactl set-cluster' command
  -n, --namespace string     Override the default context set by `apicactl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [apicactl get datasource](/get/apicactl_get_datasource)	 - Get a datasource

