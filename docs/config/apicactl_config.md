# `apicactl config`

Modify your apicactl configuration.

## Synopsis


The 'apicactl config' command lets you configure your LOGIQ CLI. If this is your first time configuring apicactl, you'll need an API token in order to use this command. To know how to generate an API token, read [https://docs.apica.io/vewing-logs/apicactl/obtaining-api-key].

Note: The values you provide during configuration will be written to the configuration file located at (~/.apicactl)


## Examples

```

View current context
	apicactl config view

Runs an interactive prompt that lets you configure apicactl
	apicactl config init

Set default cluster
	apicactl config set-cluster END-POINT

Set default context
	apicactl config set-context namespace

Runs an interactive prompt and lets you select a namespace from a list of namespaces
	apicactl config set-context i

Set API token
	apicactl config set-token api_token

```

## Options

```
  -h, --help   help for config
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

* [apicactl](/)	 - Logiqctl - CLI for Logiq Observability stack
* [apicactl config init](/config/apicactl_config_init)	 - Configure apicactl interactively
* [apicactl config set-cluster](/config/apicactl_config_set-cluster)	 - Set your LOGIQ platform endpoint
* [apicactl config set-context](/config/apicactl_config_set-context)	 - Sets the default context or namespace.
* [apicactl config set-credential](/config/apicactl_config_set-credential)	 - Set your LOGIQ user credentials
* [apicactl config set-token](/config/apicactl_config_set-token)	 - Set your LOGIQ API token
* [apicactl config view](/config/apicactl_config_view)	 - View your current apicactl configuration.

