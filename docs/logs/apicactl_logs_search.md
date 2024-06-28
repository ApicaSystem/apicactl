# `apicactl logs search`

Search logs for specific keywords or terms.

## Synopsis

Search for specific keywords or terms in logs within a namespace, app, proc

```
apicactl logs search [SearchString] [flags]
```

## Examples

```

apicactl logs search supports many time range options
  - RFC3339 and epoch timestamp formats support automatically
  - Time format in format "yyyy-MM-dd hh:mm:ss.sssss +zzzz"
  - Suffix "+zzzz" will default to UTC-to-Localtime offset
    for example, 0700 is PDT and 0000 is UTC 
    One can use option --xutc (-x) to force UTC without specifying "+zzzz"
  - Different time search range options
    * --begtime (-b) --endtime (-e) => begtime, endtime
    * --begtime (-b) and --since (-s) => begtime, begtime + duration 
    * --endtime (-e) and --since (-s) => endtime - duration, endtime
    * Single duration --since (-s) => now() - duration, now()
    * Durations --since (-s) examples are 1m, 1d, 1s, etc., default=1h

Examples:
  % apicactl -a <application_name> -p <proc_id> logs search <search_string>
  %	apicactl -a <application_name> logs search <search_string> -b "2021-07-04 23:30:00.1234 0000" -s 5m
  %	apicactl -a <application_name> logs search <search_string> -b "2021-07-04 23:30:00.1234" -e "2021-07-04 23:35:00.1234"


```

## Options

```
  -h, --help   help for search
```

## Options inherited from parent commands

```
  -a, --application string     Filter logs by application
  -b, --begtime string         Search begin time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -c, --cluster string         Override the default cluster set by `apicactl set-cluster' command
  -e, --endtime string         Search end time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -f, --follow                 Specify if the logs should be streamed.
  -m, --max-file-size int      Max output file size (default 10)
  -n, --namespace string       Override the default context set by `apicactl set-context' command
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

* [apicactl logs](/logs/apicactl_logs)	 - View logs for the given namespace and application

