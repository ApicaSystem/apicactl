# apicactl

`apicactl` is Apica Ascent's inbuilt command-line toolkit that lets you interact with the LOGIQ Observability platform without logging into the UI. Using `apicactl`, you can:
- Stream logs in real-time
- Query historical application logs
- Search within logs across namespaces
- Query and view events across your LOGIQ stack
- View and create event rules
- Create and manage dashboards
- Query and view all your resources on LOGIQ such as applications, dashboards, namespaces, processes, and queries
- Manage LOGIQ licenses

# Quickstart

The quickest way to start using `apicactl` is to download a pre-built binary from our [release page on GitHub](https://github.com/ApicaSyatem/apicactl/releases). 

## Configuring `apicactl`

Once you've downloaded the binary, you can configure `apicactl` to interact with your Apica Ascent instance by doing the following:
1. Set your cluster URL:
    ```
    apicactl config set-cluster CLUSTER_URL
    ```
1. Set your LOGIQ credentials:
    ```
    apicactl config set-ui-credential flash-userid password
    ```
1. Set your default namespace:
    ```
    apicactl config set-context NAMESPACE
    ```
1. Verify your `apicactl` configuration:
    ```
    apicactl get namespaces
    ```
This completes the installation of `apicactl`. You can now use `apicactl` to interact with your LOGIQ instance right from your terminal. 

# Building `apicactl` from source

Another way of installing `apicactl` is by building it from the source code. Building `apicactl` from its source code involves two steps:
- Installing dependencies
- Downloading and building the `apicactl` binary

## Installing dependencies

`apicactl` has the following dependencies:
- Go: You can install Go by following the instructions listed on [https://golang.org/dl/]
- Protocol Buffers: Download the binary and set it up by running the following commands:

On macOS:

```bash
PROTOC_ZIP=protoc-3.15.6-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
```

On Linux OS:
   
```bash
PROTOC_ZIP=protoc-3.15.6-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

```

This completes the installation of all `apicactl` dependencies. 

## Building `apicactl`

Run the following commands to build `apicactl` from the source code:
1. Create a directory inside your workspace in which to keep source code:
    ```
    mkdir -p $GOPATH/src/github.com/apicaai
    ```
1. Accesss the source code directory:
    ```
    cd $GOPATH/src/github.com/apicaai
    ```
1. Clone the `apicactl` GitHub repository into this folder:
    ```
    git clone git@github.com:ApicaSystem/apicactl.git
    ```
1. Access the repository you just cloned:
    ```
    cd apicactl
    ```
1. Build `apicactl`:
    ```
    go build apicactl.go
    ```
1. Make the binary `apicactl` executable:
    ```
    chmod +x ./apicactl
    ```
1. Verify the build:
    ```
    apicactl -h
    ```

`apicactl` is now built and ready for configuration and use. To configure `apicactl`, refer to the configuration instructions listed under [Configuring `apicactl`](#configuring-apicactl). 

# Available `apicactl` commands

| Command | Operation |
|---|---|
| [`apicactl config`](apicactl_config.md) | Configure `apicactl` or modify existing `apicactl` configuration |
| [`apicactl tail`](apicactl_tail.md) | Stream logs from your LOGIQ instance in real-time |
| [`apicactl create`](apicactl_create.md) | Create LOGIQ resources such as dashboards and event rules |
| [`apicactl get`](apicactl_get.md) | Display one or more LOGIQ resources |
| [`apicactl license`](apicactl_license.md) | View and manage your LOGIQ license |
| [`apicactl logs`](apicactl_logs.md) | View logs for the given namespace and application |

# Options

```
  -c, --cluster string       Override the default cluster set by `apicactl set-cluster' command
  -h, --help                 help for apicactl
  -n, --namespace string     Override the default context set by `apicactl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

To know more about the LOGIQ Observability stack, see https://apica.io/ and https://docs.apica.io/. 

In case of issues or questions, do reach out to us at [support@apica.io]. You can also [log an issue](https://github.com/ApicaSystem/apicactl/issues/new) in our `apica` source code repository on GitHub. 