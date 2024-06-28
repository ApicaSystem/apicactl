# `apicactl get applications all`

List all available applications across namespaces

## Synopsis

List all available applications across namespaces

```
apicactl get applications all [flags]
```

## Examples

```
apicactl get applications all
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

* [apicactl get applications](/get/apicactl_get_applications)	 - List all available applications within the default namespace

