## logiqctl config set-credential

Set your LOGIQ user credentials

### Synopsis


This command lets you set your LOGIQ user credentials. You'll need valid user credentials in order to access all operations.
		

```
logiqctl config set-credential [flags]
```

### Examples

```
logiqctl set-credential login password
```

### Options

```
  -h, --help   help for set-credential
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

