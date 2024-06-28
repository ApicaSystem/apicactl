# `apicactl get events`

List all available events for the namespace

## Synopsis

List all available events for the namespace

```
apicactl get events [flags]
```

## Examples

```

List last 30 events
- apicactl get events|e

List events by application 
- apicactl get events -a=sshd


```

## Options

```
  -a, --application string   Filter events by application name
  -h, --help                 help for events
  -p, --process string       Filter events by application and process name
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

