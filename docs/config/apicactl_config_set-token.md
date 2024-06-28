# `apicactl config set-token`

Set your LOGIQ API token

## Synopsis


This command lets you set your LOGIQ API token. You'll need a valid LOGIQ cluster endpoint in order to complete all operations. 
		

```
apicactl config set-token [flags]
``R

## Examples

```
apicactl set-token api_token
```

## Options

```
  -h, --help   help for set-token
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

