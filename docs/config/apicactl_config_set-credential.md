# `apicactl config set-credential`

Set your LOGIQ user credentials

## Synopsis


This command lets you set your LOGIQ user credentials. You'll need valid user credentials in order to access all operations.
		

```
apicactl config set-credential [flags]
```

## Examples

```
apicactl set-credential login password
```

## Options

```
  -h, --help   help for set-credential
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

