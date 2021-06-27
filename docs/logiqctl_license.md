## logiqctl license

View or update LOGIQ license

### Synopsis


The LOGIQ Observability platform comes preconfigured with a 30-day trial license. You can obtain a valid license by contacting LOGIQ at license@logiq.ai.
This command lets you view your existing LOGIQ license or apply a new one. 


### Examples

```

Upload your LOGIQ platform license
- logiqctl license set -f license.jws

View your LOGIQ license information
 - logiqctl license get 

```

### Options

```
  -h, --help   help for license
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
* [logiqctl license get](logiqctl_license_get.md)	 - View license information
* [logiqctl license set](logiqctl_license_set.md)	 - Configure license for LOGIQ

