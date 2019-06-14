# Mock - A configurable REST API mock written in go

Mock is a simple configurable REST API mock. Using a file with regex patterns to
match URL paths and query parameters.

```
Usage:
  mock [flags]

Flags:
  -c, --config string    config file (default is ./mock.json) (default "./mock.json")
  -h, --help             help for mock
  -L, --logfile string   logfile (defaults to stdout)
  -l, --loglevel int     log level (defaults to 4 (Info)) (default 5)
  -P, --port int         Network port (defaults to 8000) (default 8000)
  -r, --rules string     rule file (default is ./rules.json) (default "./rules.json")

```
As mock is using viper/cobra, you can provide the commandline parameters also via
the config file.

You need to provide a rules file to specify the patterns for which the mock service
should return some responses. A sample file rules.json is provided.

By default, mock is logging to the console using the loglevel Info (4). Loglevels
range from 1 (Panic) to 5 (Debug). You can specify a log file to write to.

mock is creating a http server on the specifief port. If you provide a file named
"server.crt" accompanied by by "server.key", a https server is created instead.

