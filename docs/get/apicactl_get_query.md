# `apicactl get query`

Get a query

## Synopsis

Get a query

```
apicactl get query [flags]
```

## Examples

```
apicactl get query|q <query-id>
```

## Options

```
  -h, --help   help for query
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
* [apicactl get query all](/get/apicactl_get_query_all)	 - List all the available queries
* [apicactl get query result](/get/apicactl_get_query_result)	 - Get a query result

