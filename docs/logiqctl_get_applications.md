## logiqctl get applications

List all available applications within the default namespace

### Synopsis

List all available applications within the default namespace

```
logiqctl get applications [flags]
```

### Examples

```
logiqctl get applications|apps|a
```

### Options

```
  -h, --help   help for applications
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

* [logiqctl get](logiqctl_get.md)	 - Display one or more of your LOGIQ resources
* [logiqctl get applications all](logiqctl_get_applications_all.md)	 - List all available applications across namespaces

