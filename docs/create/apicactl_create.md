# `apicactl create`

Create a LOGIQ resource

## Synopsis

This command helps you create LOGIQ resources such as dashboards and event rules from a resource specification.

## Examples

```

Create a dashboard
apicactl create dashboard -f <path to dashboard_spec_file.json>

Create eventrules
apicactl create eventrules -f <path to eventrules_file.json>

```

## Options

```
  -h, --help   help for create
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
* [apicactl create dashboard](/config/apicactl_create_dashboard)	 - Create a dashboard
* [apicactl create eventrules](/config/apicactl_create_eventrules)	 - Create an event rule

