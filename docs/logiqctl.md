## logiqctl

Logiqctl - CLI for Logiq Observability stack

### Synopsis


LOGIQ comes with an inbuilt command-line toolkit that lets you interact with the LOGIQ Observability platform without logging into the UI. Using logiqctl, you can:
- Stream logs in real-time
- Query historical application logs
- Search within logs across namespaces
- Query and view events across your LOGIQ stack
- View and create event rules
- Create and manage dashboards
- Query and view all your resources on LOGIQ such as applications, dashboards, namespaces, processes, and queries
- Manage LOGIQ licenses
- Log pattern signature extraction and reporting (max 50,000 log-lines)

Find more information, please contact support@LOGIQ.ai.


### Options

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -h, --help                 help for logiqctl
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl config](logiqctl_config.md)	 - Modify your logiqctl configuration.
* [logiqctl create](logiqctl_create.md)	 - Create a LOGIQ resource
* [logiqctl get](logiqctl_get.md)	 - Display one or more of your LOGIQ resources
* [logiqctl license](logiqctl_license.md)	 - View or update LOGIQ license
* [logiqctl logs](logiqctl_logs.md)	 - View logs for the given namespace and application
* [logiqctl tail](logiqctl_tail.md)	 - Stream logs sent to your LOGIQ Observability platform in real-time.

