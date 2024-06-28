# `apicactl license`

View or update Apica Ascent license

## Synopsis


The Apica Ascent platform comes preconfigured with a 30-day trial license. You can obtain a valid license by contacting Apica Ascent at support@apica.io.
This command lets you view your existing Apica Ascent license or apply a new one. 


## Examples

```

Upload your Apica Ascent platform license
  % apicactl license set -f license.jws

View your Apica Ascent license information 
  % apicactl license get 
 
You can obtain a valid license by contacting Apica Ascent at license@apica.ai.
This command lets you view your existing Apica Ascent license or apply a new one. 

```

## Options

```
  -h, --help   help for license
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

* [apicactl](/)	 - Logiqctl - CLI for Apica Ascent stack
* [apicactl license get](/license/apicactl_license_get)	 - View license information
* [apicactl license set](/license/apicactl_license_set)	 - Configure license for Apica Ascent

