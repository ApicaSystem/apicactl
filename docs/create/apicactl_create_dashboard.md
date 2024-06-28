## apicactl create dashboard

Create a dashboard

### Synopsis


The crowd-sourced dashboards available in https://github.com/logiqai/logiqhub can be downloaded and applied to any clusters. 
One can also export dashboards created using "apicactl get dashboard" command and apply on different clusters.


```
apicactl create dashboard [flags]
```

### Examples

```
apicactl create dashboard|d -f <path to dashboard spec>
```

### Options

```
  -f, --file string   Path to file
  -h, --help          help for dashboard
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `apicactl set-cluster' command
  -n, --namespace string     Override the default context set by `apicactl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [apicactl create](apicactl_create)	 - Create a LOGIQ resource

