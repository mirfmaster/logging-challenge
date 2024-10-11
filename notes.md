Dashboard Query
```
histogram_quantile(0.95, sum(rate(grafana_http_request_duration_seconds_bucket{}[5m])) by (le,handler))
```

