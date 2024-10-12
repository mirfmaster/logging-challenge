Dashboard Query
```
histogram_quantile(0.95, sum(rate(grafana_http_request_duration_seconds_bucket{}[5m])) by (le,handler))
```

## Fluentbit
- Fluentbit provides data pipeline feature that allows to process the logs before sending the to centralized locations
### Features
  - Parse the logs
  - filter
  - add additional information
  - even route it

```
# Explanation
[SERVICE]
  Flush     1
  Log_Level info

[INPUT] # The input will use `tail` to read app.log and tag the result as `http-service`
  Name  tail
  Path  /app/logs/app.log
  Tag   http-service

[INPUT] # 
  Name  forward
  Listen 0.0.0.0
  port 24224

[OUTPUT]
  Name  stdout
  Match *

[OUTPUT]
  Name    loki # name of the output
  Match   http-service # to only sent log with label `http-service` that provided from the input above
  host    loki # use host `loki` because the loki use docker
  port    3100
  labels  app=http-service # the log forwarded will labeled as app with http-service
  # by default loki will get 2 fiels, app(label above) & log. since we only need the log forwarded we drop the `log` key
  drop_single_key true
  line_format key_value # to transform from json to key value informations
```


## Grafana Loki
- Loki is Centralized Logging System
- Focuses on collecting, indexing and searching logs.

### Features
- Log Aggregation: Loki able to  collect log data from various sources
- Label-based indexing: Loki uses label-based indexing, similar to Prometheus.
- Integration with Grafana: Able to visualize logs

### Queries

[Full reference](https://grafana.com/docs/loki/latest/query/log_queries/)
- Basic operation
`{app="http-service"} | json`
This will extract keys and values from a json formatted log line as labels. The extracted labels can be used in label filter expressions and used as values for a range aggregation via the unwrap operation.

- Filter operation with regex of "info"
`{app="http-service"} | json |~ "info"`

- Filter operation with labels
`{app="http-service"} | json | level = "info"`

- Aggregation by counting occurrences of a specfiic log entry
`count_over_time({app="http-service"} |~ "debug" [1m])`

- Calculate the rate of log entries per second (Log per second):
`rate({app="http-service"}[1m])`

- Log per second grouped by level
`sum by (level) (rate({app="http-service"} | json [1m]))`

- Get highest log by n per second
`topk(1, sum by (level) (rate({app="http-service"} | json [1m])))`

```
avg by (level) 
(avg_over_time({app="http-service"}
| json 
| message = "request processed" 
| unwrap elapsed_ms [1m]))
```

filter by pattern

`<_>` is for ignoring
`{app="nginx"} | logfmt | pattern '<_> - - <_> "<method> <url> <_>" <status> <_> <_> "<user_agent>" <_>`

- aggregate sum grouped by user_agent
```
sum by (user_agent) (rate({app="nginx"} | logfmt | pattern `<_> - - <_> "<method> <url> <_>" <status> <_> <_> "<user_agent>" <_>`[1m]))
```



## ETC
- for each 1 sec send 2 request
`watch -n 0.5 curl "localhost"`
