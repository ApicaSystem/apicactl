## logiqctl config set-cluster

Set your LOGIQ platform endpoint

### Synopsis


This command lets you set your LOGIQ cluster endpoint. You'll need a valid LOGIQ cluster endpoint in order to complete all operations. 
		

```
logiqctl config set-cluster [flags]
```

### Examples

```
logiqctl set-cluster END-POINT
```

### Options

```
  -h, --help   help for set-cluster
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

