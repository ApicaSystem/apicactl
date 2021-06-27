## logiqctl config

Modify your logiqctl configuration.

### Synopsis


The 'logiqctl config' command lets you configure your LOGIQ CLI. If this is your first time configuring logiqctl, you'll need an API token in order to use this command. To know how to generate an API token, read https://docs.logiq.ai/vewing-logs/logiqctl/obtaining-api-key.

Note: The values you provide during configuration will be written to the configuration file located at (~/.logiqctl)


### Examples

```

View current context
	logiqctl config view

Runs an interactive prompt that lets you configure logiqctl
	logiqctl config init

Set default cluster
	logiqctl config set-cluster END-POINT

Set default context
	logiqctl config set-context namespace

Runs an interactive prompt and lets you select a namespace from a list of namespaces
	logiqctl config set-context i

Set API token
	logiqctl config set-token api_token

```

### Options

```
  -h, --help   help for config
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

* [logiqctl](logiqctl.md)	 - Logiqctl - CLI for Logiq Observability stack
* [logiqctl config init](logiqctl_config_init.md)	 - Configure logiqctl interactively
* [logiqctl config set-cluster](logiqctl_config_set-cluster.md)	 - Set your LOGIQ platform endpoint
* [logiqctl config set-context](logiqctl_config_set-context.md)	 - Sets the default context or namespace.
* [logiqctl config set-credential](logiqctl_config_set-credential.md)	 - Set your LOGIQ user credentials
* [logiqctl config set-token](logiqctl_config_set-token.md)	 - Set your LOGIQ API token
* [logiqctl config view](logiqctl_config_view.md)	 - View your current logiqctl configuration.

