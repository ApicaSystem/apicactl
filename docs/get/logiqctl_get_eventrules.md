# `logiqctl get eventrules`

List event rules

## Synopsis

List event rules

```
logiqctl get eventrules [flags]
```

## Examples

```
logiqctl get eventrule
```

## Options

```
  -w, --file string     Path to file to be written to
  -g, --groups string   list of groups separated by comma
  -h, --help            help for eventrules
```

## Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl get](/get/logiqctl_get)	 - Display one or more of your LOGIQ resources
* [logiqctl get eventrules all](/get/logiqctl_get_eventrules_all)	 - List all event rules
* [logiqctl get eventrules groups](/get/logiqctl_get_eventrules_groups)	 - List all event rules

