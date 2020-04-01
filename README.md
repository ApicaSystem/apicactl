# Logiq-box 
#### CLI for flash - Tail logs realtime

##### Build and Run
```bash

> ./generate_grpc.sh        #setup
> go build logiqbox.go      #build


> ./logiqbox                #run
NAME:
   Logiq-box - Logiq CLI Tool

USAGE:
   logiqbox [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Logiq Inc

COMMANDS:
     configure, c  Configure Logiq-box
     list, ls      List of applications that you can tail
     tail, t       tail app1 app2
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version


```