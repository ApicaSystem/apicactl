# `apicactl get applications`

List all available applications within the default namespace

## Synopsis

List all available applications within the default namespace

```
apicactl get applications [flags]
```

## Examples

```
apicactl get applications|apps|a
```

## Options

```
  -h, --help   help for applications
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

* [apicactl get](/get/apicactl_get)	 - Display one or more of your LOGIQ resources
* [apicactl get applications all](/get/apicactl_get_applications_all)	 - List all available applications across namespaces

