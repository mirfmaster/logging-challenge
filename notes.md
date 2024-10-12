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
