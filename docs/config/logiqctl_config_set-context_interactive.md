# `logiqctl config set-context interactive`

Run an interactive prompt that lets you select a namespace from a list.

## Synopsis


This command lets you set a default context interactively. Running 'logiqctl config set-context interactive' brings up an interactive list of namespaces from which you can select a namespace and set a context.
		

```
logiqctl config set-context interactive [flags]
```

## Options

```
  -h, --help   help for interactive
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

* [logiqctl config set-context](/config/logiqctl_config_set-context)	 - Sets the default context or namespace.

