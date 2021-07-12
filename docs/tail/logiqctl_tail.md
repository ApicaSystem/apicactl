# `logiqctl tail`

Stream logs sent to your LOGIQ Observability platform in real-time.

## Synopsis


The 'logiqctl tail' command is similar to the 'tail -f' command. It allows you to stream the log data that is being sent to your LOGIQ Observability platform in real-time. You can see logs from the cluster at multiple levels. Running the command 'tail' without any options brings up an interactive prompt that lets you choose an application and process in the current context. 


```
logiqctl tail [flags]
```

## Examples

```

Tail logs 
  % logiqctl tail
  % logiqctl tail -g

```

## Options

```
  -h, --help                   help for tail
  -m, --max-file-size int      Max output file size (default 10)
  -g, --psmod                  Enable pattern signature generation module
  -w, --write-to-file string   Path to file
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

