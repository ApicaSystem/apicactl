## logiqctl logs

View logs for the given namespace and application

### Synopsis


logs command is used to view historical logs. This expects a namespace and an application to be available to return results. Set the default namespace using 'logiqctl set-context' command or pass as '-n=NAMESPACE' flag. Application name needs to be passed as an argument to the command or use the 'interactive' command to choose from the list of available applications and processes.   

**Note:**
- The global flag '--time-format' is not applicable for this command.
- The global flag '--output' only supports JSON format for this command.

```
logiqctl logs [flags]
```

### Examples

```

Print logs for the LOGIQ ingest server
- logiqctl logs -a <application_name>

Print logs in JSON format
- logiqctl -o=json logs -a <application_name>

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for logiq-flash ingest server filtered by process logiq-flash-2
The --process (-p) flag lets you view logs for the individual pod
- logiqctl logs -p=<proc_id> -a <application_name>

Runs an interactive prompt to let user choose filters
- logiqctl logs interactive|i

Search logs for specific keywords or terms
- logiqctl logs -a <application_name> search <searchterm>
- logiqctl logs -a <application_name> -p <proc_id> search <searchterm>

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

- stream logs contains log pattern-signature (PS).
- Example:  % logiqctl config set-context <namespace>
            % logiqctl logs -a <proc_id> -s 10s -f 
            % logiqctl logs -a <application_name> -p -s 10s -f
            % logiqctl logs -a <application_name> -s 10s -w outputfile.txt
  (You might want to pipe above dump into file for later cross-reference)
- after done logs streaming, two files will be created.
  notice that these files are reset for every logs query session.
  * ps_stat.out: compute byte and log counts and percentage for each pattern signature 


```

### Options

```
  -a, --application string     Filter logs by application
  -f, --follow                 Specify if the logs should be streamed.
  -h, --help                   help for logs
  -m, --max-file-size int      Max output file size (default 10)
      --page-size uint32       Number of log entries to return in one page (default 30)
  -p, --process string         Filter logs by  proc id
  -s, --since string           Only return logs newer than a relative duration. This is in relative to the last
                               seen log time for a specified application or processes within the namespace.
                               A duration string is a possibly signed sequence of decimal numbers, each with optional
                               fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h" (default "1h")
  -w, --write-to-file string   Path to file
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
* [logiqctl logs interactive](logiqctl_logs_interactive.md)	 - Runs an interactive prompt to let the user select application and filters
* [logiqctl logs search](logiqctl_logs_search.md)	 - Search given text in logs

