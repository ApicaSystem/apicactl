# `logiqctl logs`

View logs for the given namespace and application

## Synopsis


The 'logs' command is used to view historical logs. This command expects a namespace and an application to be available to return results. You can set the default namespace using the 'logiqctl set-context' command or pass the namespace as '-n=NAMESPACE' flag. The application name also needs to be passed as an argument to the command. You can also use the 'interactive' command to choose from the list of available applications and processes.   

!!! tip "Note"
  - The global flag '--time-format' is not applicable for this command.
  - The global flag '--output' only supports JSON format for this command.

```
logiqctl logs [flags]
```

## Examples

```

Print logs for the LOGIQ ingest server
  % logiqctl logs -a <application_name>

Print logs in JSON format:
  % logiqctl -o=json logs -a <application_name>

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for logiq-flash ingest server filtered by process logiq-flash-2
The --process (-p) flag lets you view logs for the individual pod
  % logiqctl logs -p=<proc_id> -a <application_name>

Runs an interactive prompt that lets you choose filters
  % logiqctl logs interactive|i

Search logs for specific keywords or terms see help:
  % logiqctl logs search --help

More examples:  
  % logiqctl logs -a <application_name> search <searchterm>
  % logiqctl logs -a <application_name> -p <proc_id> search <searchterm>
  % logiqctl logs -a <application_name> -p <proc_id> search <searchterm> -g

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

One can automatically generate pattern-signature (PS) for logs using flag --psmod (-g).
Add-on executable "psmod" from logiqhub is required to run side-by-side with logiqctl. 
Enable PS generation will generate stat file ps_stat.out that computes byte and log counts and 
percentage for each pattern signature 

More examples:  
  % logiqctl config set-context <namespace>
  % logiqctl logs -a <application_name> 
  % logiqctl logs -a <application_name> -p <proc_id_name> 
  % logiqctl logs -a <application_name> -p <proc_id_name> -g


```

## Options

```
  -a, --application string     Filter logs by application
  -b, --begtime string         Search begin time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -e, --endtime string         Search end time range format "yyyy-MM-dd hh:mm:ss +0000". 
                               "+0000" suffix is required for search using UTC time.  
                               Localtime time search is assumed WITHOUT specifying "+0000."
  -f, --follow                 Specify if the logs should be streamed.
  -h, --help                   help for logs
  -m, --max-file-size int      Max output file size (default 10)
      --page-size uint32       Number of log entries to return in one page (default 30)
  -p, --process string         Filter logs by  proc id
  -g, --psmod                  Enable pattern signature generation module
  -s, --since string           Only return logs newer than a relative duration. This is in relative to the last
                               seen log time for a specified application or processes within the namespace.
                               A duration string is a possibly signed sequence of decimal numbers, each with optional
                               fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h"
  -w, --write-to-file string   Path to file
  -x, --xutc                   Force UTC date-time
```

## Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl](/)	 - Logiqctl - CLI for Logiq Observability stack
* [logiqctl logs interactive](/logs/logiqctl_logs_interactive)	 - Runs an interactive prompt to display logs.
* [logiqctl logs search](/logs/logiqctl_logs_search)	 - Search logs for specific keywords or terms.

