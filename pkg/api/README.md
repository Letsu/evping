# API Requests
Möglichkeiten:

  GET - host:port/api/allhosts
  
  GET - host:port/api/dataofhost

  POST - host:port/api/host

  DELETE - host:port/api/host

## GET - allhosts
Respond Body (JSON):
```
[
    {
        "host": "1.1.1.1",
        "ping_frequency": 1
    },
    {
        "host": "8.8.8.8",
        "ping_frequency": 1
    }
]
```

## GET - dataofhost
Request Body (JSON):
```
{
	"ip_addresses": ["1.1.1.1","8.8.8.8"],
	"start_time": "2024-01-05T08:37:59.159143+01:00",
	"end_time": "2024-01-05T08:38:02.159143+01:00"
}
```
Respond Body (JSON):
```
[
    {
        "ip_address": "1.1.1.1",
        "rows": [
            {
                "timestamp": "2024-01-05T08:37:59.159143+01:00",
                "ip_address": "1.1.1.1",
                "dns_name": "1.1.1.1",
                "rtt": "91.049ms"
            },
            {
                "timestamp": "2024-01-05T08:38:00.087686+01:00",
                "ip_address": "1.1.1.1",
                "dns_name": "1.1.1.1",
                "rtt": "20.181ms"
            },
            {
                "timestamp": "2024-01-05T08:38:01.145903+01:00",
                "ip_address": "1.1.1.1",
                "dns_name": "1.1.1.1",
                "rtt": "77.802ms"
            }
        ]
    },
    {
        "ip_address": "8.8.8.8",
        "rows": [
            {
                "timestamp": "2024-01-05T08:37:59.159143+01:00",
                "ip_address": "8.8.8.8",
                "dns_name": "8.8.8.8",
                "rtt": "91.049ms"
            },
            {
                "timestamp": "2024-01-05T08:38:00.087686+01:00",
                "ip_address": "8.8.8.8",
                "dns_name": "8.8.8.8",
                "rtt": "20.181ms"
            },
            {
                "timestamp": "2024-01-05T08:38:01.145903+01:00",
                "ip_address": "8.8.8.8",
                "dns_name": "8.8.8.8",
                "rtt": "77.802ms"
            }
        ]
    }
]
```
## POST - host
Request Body (JSON):
```
{
	"ip_address": "8.8.8.8",
	"ping_frequency": 1
}
```

## DELETE - host
Request Body (JSON):
```
{
	"ip_address": "169.254.13.18"
}
```
