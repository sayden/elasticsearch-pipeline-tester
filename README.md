# elasticsearch-pipeline-tester
A simple CLI tool to use the _simulate API of elasticsearch to quickly test pipelines

```shell
usage: pipelinetester [<flags>] <pipeline> [<logs>]

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -e, --elasticsearch="http://localhost:9200"  
                                 URL of Elasticsearch. Scheme and port are mandatory (but HTTPS is not tested xD)
  -n, --pipeline-name="testing"  Name of the pipeline to create on elasticsearch
  -u, --ugly                     Deactivate pretty printing
  -b, --bulk                     Use bulk if you want to make a single request to ES and get a single response with many documents
      --skip=0                   Skip the N first documents
      --total=0                  Process a total of N documents. You can use it with skip
  -i, --stdin                    Instead of using an input log file, read from stdin

Args:
  <pipeline>  File path of the pipeline to use
  [<logs>]    The log filepath to parse. Use stdin flag if you want to work just on the cli

```

## Example usage

```shell

```