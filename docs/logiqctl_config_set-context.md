## logiqctl config set-context

Sets the default context or namespace.

### Synopsis


All logiqctl operations require a context or namespace. To override the default namespace set for an individual command, use the flag '-n'. 
		

```
logiqctl config set-context [flags]
```

### Examples

```
set-context <namespace name>
```

### Options

```
  -h, --help   help for set-context
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl config](logiqctl_config.md)	 - Modify your logiqctl configuration.
* [logiqctl config set-context interactive](logiqctl_config_set-context_interactive.md)	 - Run an interactive prompt that lets you select a namespace from a list.

