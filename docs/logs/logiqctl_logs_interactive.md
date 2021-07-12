# `logiqctl logs interactive`

Runs an interactive prompt to display logs.

## Synopsis

Runs an interactive prompt to display logs.

```
logiqctl logs interactive [flags]
```

## Options

```
  -h, --help   help for interactive
```

## Options inherited from parent commands

```
  -a, --application string     Filter logs by application
  -b, --begtime string         Search begin time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -c, --cluster string         Override the default cluster set by `logiqctl set-cluster' command
  -e, --endtime string         Search end time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -f, --follow                 Specify if the logs should be streamed.
  -m, --max-file-size int      Max output file size (default 10)
  -n, --namespace string       Override the default context set by `logiqctl set-context' command
  -o, --output string          Output format. One of: table|json|yaml. 
                               JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
      --page-size uint32       Number of log entries to return in one page (default 30)
  -p, --process string         Filter logs by  proc id
  -g, --psmod                  Enable pattern signature generation module
  -s, --since string           Only return logs newer than a relative duration. This is in relative to the last
                               seen log time for a specified application or processes within the namespace.
                               A duration string is a possibly signed sequence of decimal numbers, each with optional
                               fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h"
  -t, --time-format string     Time formatting options. One of: relative|epoch|RFC3339. 
                               This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
  -w, --write-to-file string   Path to file
  -x, --xutc                   Force UTC date-time
```

## SEE ALSO

* [logiqctl logs](/logs/logiqctl_logs)	 - View logs for the given namespace and application

