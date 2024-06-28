# `apicactl config set-context`

Sets the default context or namespace.

## Synopsis


All apicactl operations require a context or namespace. To override the default namespace set for an individual command, use the flag '-n'. 
		

```
apicactl config set-context [flags]
```

## Examples

```
set-context <namespace name>
```

## Options

```
  -h, --help   help for set-context
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

* [apicactl config](/config/apicactl_config)	 - Modify your apicactl configuration.
* [apicactl config set-context interactive](/config/apicactl_config_set-context_interactive)	 - Run an interactive prompt that lets you select a namespace from a list.

