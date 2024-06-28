# `apicactl get`

Display one or more of your LOGIQ resources

## Synopsis

Prints a table that displays the most important information about the LOGIQ resource you specify.

## Examples

```

List all applications for the selected context
apicactl get applications

List all applications for all available contexts
apicactl get applications all

List all dashboards
apicactl get dashboards all

Get dashboard
apicactl get dashboard dashboard-slug

List events for the available namespace
apicactl get events

List events for all available namespaces
apicactl get events all

List or export eventrules
apicactl get eventrules

List all namespaces
apicactl get namespaces

List all processes
apicactl get processes

List all queries
apicactl get queries all

Get httpingestkey
apicactl get httpingestkey

Get query
apicactl get query query-slug

```

## Options

```
  -h, --help   help for get
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

* [apicactl](/)	 - Logiqctl - CLI for Logiq Observability stack
* [apicactl get applications](/get/apicactl_get_applications)	 - List all available applications within the default namespace
* [apicactl get dashboard](/get/apicactl_get_dashboard)	 - Get a dashboard
* [apicactl get datasource](/get/apicactl_get_datasource)	 - Get a datasource
* [apicactl get eventrules](/get/apicactl_get_eventrules)	 - List event rules
* [apicactl get events](/get/apicactl_get_events)	 - List all available events for the namespace
* [apicactl get httpingestkey](/get/apicactl_get_httpingestkey)	 - Get httpingestkey
* [apicactl get namespaces](/get/apicactl_get_namespaces)	 - List all available namespaces
* [apicactl get processes](/get/apicactl_get_processes)	 - List all available processes. This command runs an interactive prompt that lets you choose an application from a list of available applications.
* [apicactl get query](/get/apicactl_get_query)	 - Get a query

