# `apicactl logs`

View logs for the given namespace and application

## Synopsis


The 'logs' command is used to view historical logs. This command expects a namespace and an application to be available to return results. You can set the default namespace using the 'apicactl set-context' command or pass the namespace as '-n=NAMESPACE' flag. The application name also needs to be passed as an argument to the command. You can also use the 'interactive' command to choose from the list of available applications and processes.   

!!! tip "Note"
  - The global flag '--time-format' is not applicable for this command.
  - The global flag '--output' only supports JSON format for this command.

```
apicactl logs [flags]
```

## Examples

```

Print logs for the LOGIQ ingest server
  % apicactl logs -a <application_name>

Print logs in JSON format:
  % apicactl -o=json logs -a <application_name>

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for apica-flash ingest server filtered by process apica-flash-2
The --process (-p) flag lets you view logs for the individual pod
  % apicactl logs -p=<proc_id> -a <application_name>

Runs an interactive prompt that lets you choose filters
  % apicactl logs interactive|i

Search logs for specific keywords or terms see help:
  % apicactl logs search --help

More examples:  
  % apicactl logs -a <application_name> search <searchterm>
  % apicactl logs -a <application_name> -p <proc_id> search <searchterm>
  % apicactl logs -a <application_name> -p <proc_id> search <searchterm> -g

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

One can automatically generate pattern-signature (PS) for logs using flag --psmod (-g).
Add-on executable "psmod" from apicahub is required to run side-by-side with apicactl. 
Enable PS generation will generate stat file ps_stat.out that computes byte and log counts and 
percentage for each pattern signature 

More examples:  
  % apicactl config set-context <namespace>
  % apicactl logs -a <application_name> 
  % apicactl logs -a <application_name> -p <proc_id_name> 
  % apicactl logs -a <application_name> -p <proc_id_name> -g


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
  -c, --cluster string       Override the default cluster set by `apicactl set-cluster' command
  -n, --namespace string     Override the default context set by `apicactl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [apicactl](/)	 - Logiqctl - CLI for Logiq Observability stack
* [apicactl logs interactive](/logs/apicactl_logs_interactive)	 - Runs an interactive prompt to display logs.
* [apicactl logs search](/logs/apicactl_logs_search)	 - Search logs for specific keywords or terms.

