# elasticsearch-pipeline-tester
A simple CLI tool to use the _simulate API of elasticsearch to quickly test pipelines

```shell
usage: pipelinetester [<flags>] <pipeline> [<logs>]

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -e, --elasticsearch="http://localhost:9200"  
                                 URL of Elasticsearch. Scheme and port are mandatory (but HTTPS is not tested xD)
  -p, --pipeline-name="testing"  Name of the pipeline to create on elasticsearch
  -u, --ugly                     Deactivate pretty printing
  -b, --bulk                     Use bulk if you want to make a single request to ES and get a single response with many documents
  -s, --skip=0                   Skip the N first documents
  -n, --total=0                  Process a total of N documents. You can use it with skip
  -i, --stdin                    Instead of using an input log file, read from stdin

Args:
  <pipeline>  File path of the pipeline to use
  [<logs>]    The log filepath to parse. Use stdin flag if you want to work just on the cli

```

## Example usage

```shell
$ pwd                                                                                                       [11:16:24]
/home/sayden/go/src/github.com/elastic/beats/filebeat/module/kibana
$ head log/ingest/pipeline.yml                                                                              [11:16:25]
description: Pipeline for parsing Kibana logs
on_failure:
- set:
    field: error.message
    value: '{{ _ingest.on_failure_message }}'
processors:
- set:
    field: event.ingested
    value: '{{_ingest.timestamp}}'
- rename:
$
$ head log/test/test.log                                                                                    [11:16:37]
{"type":"response","@timestamp":"2018-05-09T10:57:55Z","tags":[],"pid":69410,"method":"get","statusCode":304,"req":{"url":"/ui/fonts/open_sans/open_sans_v15_latin_600.woff2","method":"get","headers":{"host":"localhost:5601","connection":"keep-alive","origin":"http://localhost:5601","user-agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36","accept":"*/*","referer":"http://localhost:5601/app/kibana","accept-encoding":"gzip, deflate, br","accept-language":"en-US,en;q=0.9,de;q=0.8","if-none-match":"\"24234c1c81b3948758c1a0be8e5a65386ca94c52\"","if-modified-since":"Thu, 03 May 2018 09:45:28 GMT"},"remoteAddress":"127.0.0.1","userAgent":"127.0.0.1","referer":"http://localhost:5601/app/kibana"},"res":{"statusCode":304,"responseTime":26,"contentLength":9},"message":"GET /ui/fonts/open_sans/open_sans_v15_latin_600.woff2 304 26ms - 9.0B"}
{"type":"log","@timestamp":"2018-05-09T10:59:12Z","tags":["debug","monitoring-ui","kibana-monitoring"],"pid":69776,"message":"Fetching data from kibana_stats collector"}
{"type":"log","@timestamp":"2018-05-09T10:59:12Z","tags":["reporting","debug","exportTypes"],"pid":69776,"message":"Found exportType at /Users/ruflin/Downloads/6.3/kibana-6.3.0-darwin-x86_64/node_modules/x-pack/plugins/reporting/export_types/csv/server/index.js"}
{"type":"response","@timestamp":"2018-05-09T10:57:55Z","tags":[],"pid":69410,"method":"get","statusCode":304,"req":{"url":"/ui/fonts/open_sans/open_sans_v15_latin_600.woff2","method":"get","headers":{"host":"localhost:5601","connection":"keep-alive","origin":"http://localhost:5601","user-agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36","accept":"*/*","referer":"http://localhost:5601/app/kibana","accept-encoding":"gzip, deflate, br","accept-language":"en-US,en;q=0.9,de;q=0.8","if-none-match":"\"24234c1c81b3948758c1a0be8e5a65386ca94c52\"","if-modified-since":"Thu, 03 May 2018 09:45:28 GMT"},"remoteAddress":"127.0.0.1","userAgent":"127.0.0.1","referer":"http://localhost:5601/app/kibana"},"res":{"statusCode":304,"responseTime":3000,"contentLength":9},"message":"GET /ui/fonts/open_sans/open_sans_v15_latin_600.woff2 304 26ms - 9.0B"}
$
$ pipelinetester -n 1 ~/go/src/github.com/elastic/beats/filebeat/module/kibana/log/ingest/pipeline.yml ~/go/src/github.com/elastic/beats/filebeat/module/kibana/log/test/test.log
  INFO[0000](2021-10-13T11:25:05.162856172+02:00) github.com/sayden/pipelinetester/main.go:217 Pipeline insertion output:
{
  "acknowledged": true
}

  INFO[0000](2021-10-13T11:25:05.164341632+02:00) github.com/sayden/pipelinetester/main.go:168 Response:
{sayden➜github.com/sayden/pipelinetester(main✗)» ./pipelinetester -n 1 ~/go/src/github.com/elastic/beats/filebeat/module/kibana/log/ingest/pipeline.yml ~/go/src/github.com/elastic/beats/filebeat/module/kibana/log/test/test.log                     [11:20:17]
Pipeline output:
{
  "acknowledged": true
}
  INFO[0000](2021-10-13T11:22:31.629405406+02:00) github.com/sayden/pipelinetester/main.go:168 Result:
{
  "docs": [
    {
      "doc": {
        "_id": "",
        "_index": "",
        "_ingest": {
          "timestamp": "2021-10-13T09:22:31.628472071Z"
        },
        "_source": {
          "error": {
            "message": "field [@timestamp] doesn't exist"
          },
          "event": {
            "ingested": "2021-10-13T09:22:31.628472071Z"
          },
          "message": {
            "@timestamp": "2018-05-09T10:57:55Z",
            "message": "GET /ui/fonts/open_sans/open_sans_v15_latin_600.woff2 304 26ms - 9.0B",
            "method": "get",
            "pid": 69410,
            "req": {
              "headers": {
                "accept": "*/*",
                "accept-encoding": "gzip, deflate, br",
                "accept-language": "en-US,en;q=0.9,de;q=0.8",
                "connection": "keep-alive",
                "host": "localhost:5601",
                "if-modified-since": "Thu, 03 May 2018 09:45:28 GMT",
                "if-none-match": "\"24234c1c81b3948758c1a0be8e5a65386ca94c52\"",
                "origin": "http://localhost:5601",
                "referer": "http://localhost:5601/app/kibana",
                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
              },
              "method": "get",
              "referer": "http://localhost:5601/app/kibana",
              "remoteAddress": "127.0.0.1",
              "url": "/ui/fonts/open_sans/open_sans_v15_latin_600.woff2",
              "userAgent": "127.0.0.1"
            },
            "res": {
              "contentLength": 9,
              "responseTime": 26,
              "statusCode": 304
            },
            "statusCode": 304,
            "tags": [],
            "type": "response"
          }
        }
      }
    }
  ]
  "docs": [
    {
      "doc": {
        "_id": "",
        "_index": "",
        "_ingest": {
          "timestamp": "2021-10-13T09:25:05.163496474Z"
        },
        "_source": {
          "error": {
            "message": "field [@timestamp] doesn't exist"
          },
          "event": {
            "ingested": "2021-10-13T09:25:05.163496474Z"
          },
          "message": {
            "@timestamp": "2018-05-09T10:57:55Z",
            "message": "GET /ui/fonts/open_sans/open_sans_v15_latin_600.woff2 304 26ms - 9.0B",
            "method": "get",
            "pid": 69410,
            "req": {
              "headers": {
                "accept": "*/*",
                "accept-encoding": "gzip, deflate, br",
                "accept-language": "en-US,en;q=0.9,de;q=0.8",
                "connection": "keep-alive",
                "host": "localhost:5601",
                "if-modified-since": "Thu, 03 May 2018 09:45:28 GMT",
                "if-none-match": "\"24234c1c81b3948758c1a0be8e5a65386ca94c52\"",
                "origin": "http://localhost:5601",
                "referer": "http://localhost:5601/app/kibana",
                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
              },
              "method": "get",
              "referer": "http://localhost:5601/app/kibana",
              "remoteAddress": "127.0.0.1",
              "url": "/ui/fonts/open_sans/open_sans_v15_latin_600.woff2",
              "userAgent": "127.0.0.1"
            },
            "res": {
              "contentLength": 9,
              "responseTime": 26,
              "statusCode": 304
            },
            "statusCode": 304,
            "tags": [],
            "type": "response"
          }
        }
      }
    }
  ]
}
```