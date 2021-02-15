package model

const GrafanaDashboardJSON = `{
      "__inputs": [
        {
          "name": "DS_PROMETHEUS",
          "label": "Prometheus",
          "description": "",
          "type": "datasource",
          "pluginId": "prometheus",
          "pluginName": "Prometheus"
        }
      ],
      "__requires": [
        {
          "type": "grafana",
          "id": "grafana",
          "name": "Grafana",
          "version": "6.2.1"
        },
        {
          "type": "panel",
          "id": "graph",
          "name": "Graph",
          "version": ""
        },
        {
          "type": "datasource",
          "id": "prometheus",
          "name": "Prometheus",
          "version": "1.0.0"
        },
        {
          "type": "panel",
          "id": "singlestat",
          "name": "Singlestat",
          "version": ""
        }
      ],
      "annotations": {
        "list": [
          {
            "builtIn": 1,
            "datasource": "-- Grafana --",
            "enable": true,
            "hide": true,
            "iconColor": "rgba(0, 211, 255, 1)",
            "name": "Annotations & Alerts",
            "type": "dashboard"
          }
        ]
      },
      "description": "Monitoring Keycloak metrics",
	  "editable": true,
	  "gnetId": null,
	  "graphTooltip": 1,
	  "id": 2,
	  "iteration": 1613387898703,
	  "links": [],
	  "panels": [
		{
		  "collapsed": false,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 0
		  },
		  "id": 33,
		  "panels": [],
		  "title": "System Health",
		  "type": "row"
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 8,
			"w": 5,
			"x": 0,
			"y": 1
		  },
		  "id": 25,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null as zero",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(up{namespace=\"$namespace\", job=\"keycloak\"}) by (pod)",
			  "format": "time_series",
			  "hide": false,
			  "intervalFactor": 1,
			  "legendFormat": "Pod {{pod}}",
			  "refId": "A"
			},
			{
			  "expr": "sum((wildfly_datasources_pool_available_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"} + wildfly_datasources_pool_active_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}) >bool 0) by (pod)",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "interval": "",
			  "intervalFactor": 1,
			  "legendFormat": "Database",
			  "refId": "B"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Readiness Probes",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "decimals": 0,
			  "format": "short",
			  "label": "Ready",
			  "logBase": 1,
			  "max": "1",
			  "min": "0",
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": false
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 8,
			"w": 10,
			"x": 5,
			"y": 1
		  },
		  "id": 27,
		  "legend": {
			"alignAsTable": false,
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"rightSide": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [
			{}
		  ],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(increase(keycloak_response_errors{namespace=\"$namespace\", job=\"keycloak\"}[30m])) by (route, code)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "",
			  "refId": "A"
			},
			{
			  "expr": "",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "refId": "B"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "HTTP Errors [30m]",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": "# of errors",
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": false
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "cacheTimeout": null,
		  "dashLength": 10,
		  "dashes": false,
		  "description": "",
		  "fill": 1,
		  "gridPos": {
			"h": 8,
			"w": 9,
			"x": 15,
			"y": 1
		  },
		  "id": 51,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pluginVersion": "6.2.4",
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/auth\"}[15m])) / count(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/auth\"}[15m]))",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "/auth",
			  "refId": "A"
			},
			{
			  "expr": "sum(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/token\"}[15m])) / count(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/token\"}[15m]))",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "/token",
			  "refId": "B"
			},
			{
			  "expr": "sum(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/introspect\"}[15m])) / count(rate(keycloak_request_duration_sum{namespace=\"$namespace\", route=~\".+/introspect\"}[15m]))",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "/introspect",
			  "refId": "C"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Endpoint latency [15m]",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "s",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "collapsed": false,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 9
		  },
		  "id": 35,
		  "panels": [],
		  "title": "Detailed metrics",
		  "type": "row"
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "datasource": "Prometheus",
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 0,
			"y": 10
		  },
		  "hideTimeOverride": false,
		  "id": 12,
		  "legend": {
			"alignAsTable": false,
			"avg": true,
			"current": true,
			"hideEmpty": true,
			"hideZero": true,
			"max": true,
			"min": true,
			"rightSide": false,
			"show": true,
			"sideWidth": 70,
			"total": false,
			"values": true
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "connected",
		  "options": {},
		  "percentage": false,
		  "pointradius": 5,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "kube_pod_container_resource_requests{namespace=\"$namespace\", resource=\"memory\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "Requests",
			  "refId": "A"
			},
			{
			  "expr": "kube_pod_container_resource_limits{namespace=\"$namespace\", resource=\"memory\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "hide": false,
			  "intervalFactor": 1,
			  "legendFormat": "Limits",
			  "refId": "C"
			},
			{
			  "expr": "container_memory_rss{namespace=\"$namespace\", container=\"keycloak\"}  ",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "RSS for {{pod}}",
			  "refId": "B"
			},
			{
			  "expr": "sum(jvm_memory_bytes_committed{namespace=\"$namespace\",job=\"keycloak\"}) by (pod)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "JVM Committed for {{pod}}",
			  "refId": "D"
			},
			{
			  "expr": "sum(jvm_memory_bytes_max{namespace=\"$namespace\",job=\"keycloak\"}) by (pod)",
			  "format": "time_series",
			  "hide": false,
			  "intervalFactor": 1,
			  "legendFormat": "JVM Max for {{pod}}",
			  "refId": "E"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Pod Memory Usage",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "decimals": null,
			  "format": "bytes",
			  "label": "",
			  "logBase": 1,
			  "max": null,
			  "min": "0",
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 2,
			  "max": null,
			  "min": null,
			  "show": false
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 5,
			"x": 7,
			"y": 10
		  },
		  "id": 31,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "kube_pod_container_resource_limits{namespace=\"$namespace\", resource=\"cpu\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "CPU Limit for {{pod}}",
			  "refId": "A"
			},
			{
			  "expr": "kube_pod_container_resource_requests{namespace=\"$namespace\", resource=\"cpu\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "CPU Request for {{pod}}",
			  "refId": "C"
			},
			{
			  "expr": "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=\"$namespace\", container=\"keycloak\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "CPU Load for {{pod}}",
			  "refId": "B"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "CPU Load",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "cacheTimeout": null,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 5,
			"x": 12,
			"y": 10
		  },
		  "id": 37,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null as zero",
		  "options": {},
		  "percentage": false,
		  "pluginVersion": "6.2.4",
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(increase(jvm_gc_collection_seconds_count{namespace=\"$namespace\", job=\"keycloak\"}[30m])) by (pod)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "GC Collections for {{pod}}",
			  "refId": "A"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Number of GC Collections [30m]",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 17,
			"y": 10
		  },
		  "id": 45,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(vendor_cache_manager_default_cache_sessions_statistics_number_of_entries{namespace=\"datagrid\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Sessions",
			  "refId": "A"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_offline_sessions_statistics_number_of_entries{namespace=\"datagrid\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Offline Sessions",
			  "refId": "B"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_client_sessions_statistics_number_of_entries{namespace=\"datagrid\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Client Sessions",
			  "refId": "C"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_offline_client_sessions_statistics_number_of_entries{namespace=\"$datagrid_namespace\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Offline Client Sessions",
			  "refId": "D"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_login_failures_statistics_number_of_entries{namespace=\"$datagrid_namespace\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Offline Client Sessions",
			  "refId": "E"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_action_tokens_statistics_number_of_entries{namespace=\"$datagrid_namespace\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Offline Client Sessions",
			  "refId": "F"
			},
			{
			  "expr": "sum(vendor_cache_manager_default_cache_work_statistics_number_of_entries{namespace=\"$datagrid_namespace\"}) by (job)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Cache Offline Client Sessions",
			  "refId": "G"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Number of entries in caches",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "datasource": "Prometheus",
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 0,
			"y": 16
		  },
		  "hideTimeOverride": false,
		  "id": 39,
		  "legend": {
			"alignAsTable": false,
			"avg": true,
			"current": true,
			"hideEmpty": true,
			"hideZero": true,
			"max": true,
			"min": true,
			"rightSide": false,
			"show": true,
			"sideWidth": 70,
			"total": false,
			"values": true
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "connected",
		  "options": {},
		  "percentage": false,
		  "pointradius": 5,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "kube_pod_container_resource_requests{namespace=\"$namespace\", resource=\"memory\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "hide": false,
			  "instant": false,
			  "intervalFactor": 1,
			  "legendFormat": "Requests",
			  "refId": "A"
			},
			{
			  "expr": "kube_pod_container_resource_limits{namespace=\"$namespace\", resource=\"memory\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Limits",
			  "refId": "C"
			},
			{
			  "expr": "sum(base_memory_usedHeap_bytes{namespace=\"$namespace\"})",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "JVM Heap for {{pod}}",
			  "refId": "E"
			},
			{
			  "expr": "sum(base_memory_usedNonHeap_bytes{namespace=\"$namespace\"})",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "JVM Non Heap for {{pod}}",
			  "refId": "B"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "JVM Memory Usage",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "decimals": null,
			  "format": "bytes",
			  "label": "",
			  "logBase": 1,
			  "max": null,
			  "min": "0",
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 2,
			  "max": null,
			  "min": null,
			  "show": false
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "cacheTimeout": null,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 5,
			"x": 12,
			"y": 16
		  },
		  "id": 38,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null as zero",
		  "options": {},
		  "percentage": false,
		  "pluginVersion": "6.2.4",
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum(increase(jvm_gc_collection_seconds_sum{namespace=\"$namespace\", job=\"keycloak\"}[30m])) by (pod)",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Total GC Time for {{pod}}",
			  "refId": "A"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Total GC Time [30m]",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "s",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 17,
			"y": 16
		  },
		  "id": 48,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "wildfly_datasources_pool_active_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Active for {{pod}}",
			  "refId": "A"
			},
			{
			  "expr": "wildfly_datasources_pool_available_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Available for {{pod}}",
			  "refId": "B"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Number of database connections in pool",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 17,
			"y": 22
		  },
		  "id": 49,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "increase(wildfly_datasources_pool_xacommit_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}[30m])",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Number of commits for {{pod}}",
			  "refId": "A"
			},
			{
			  "expr": "increase(wildfly_datasources_pool_xarollback_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}[30m])",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Number of rollbacks for {{pod}}",
			  "refId": "B"
			},
			{
			  "expr": "increase(wildfly_datasources_pool_blocking_failure_count{namespace=\"$namespace\", xa_data_source=\"keycloak_postgresql-DB\"}[30m])",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "Number of failures for {{pod}}",
			  "refId": "C"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Number of transactions [30m]",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "cacheTimeout": null,
		  "dashLength": 10,
		  "dashes": false,
		  "fill": 1,
		  "gridPos": {
			"h": 6,
			"w": 7,
			"x": 17,
			"y": 28
		  },
		  "id": 50,
		  "legend": {
			"avg": false,
			"current": false,
			"max": false,
			"min": false,
			"show": true,
			"total": false,
			"values": false
		  },
		  "lines": true,
		  "linewidth": 1,
		  "links": [],
		  "nullPointMode": "null",
		  "options": {},
		  "percentage": false,
		  "pluginVersion": "6.2.4",
		  "pointradius": 2,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "wildfly_io_io_thread_count{namespace=\"$namespace\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "IO Threads for {{pod}}",
			  "refId": "A"
			},
			{
			  "expr": "wildfly_io_queue_size{namespace=\"$namespace\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "I/O Queue size for {{pod}}",
			  "refId": "B"
			},
			{
			  "expr": "wildfly_io_busy_task_thread_count{namespace=\"$namespace\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "I/O Busy tasks for {{pod}}",
			  "refId": "C"
			},
			{
			  "expr": "wildfly_io_max_pool_size{namespace=\"$namespace\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "IO Max pool size for {{pod}}",
			  "refId": "E"
			},
			{
			  "expr": "kube_pod_container_resource_limits{namespace=\"$namespace\", resource=\"cpu\", pod=\"keycloak-0\"}",
			  "format": "time_series",
			  "intervalFactor": 1,
			  "legendFormat": "CPU limit",
			  "refId": "D"
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "IO Threads",
		  "tooltip": {
			"shared": true,
			"sort": 0,
			"value_type": "individual"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "short",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "collapsed": true,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 34
		  },
		  "id": 41,
		  "panels": [
			{
			  "aliasColors": {},
			  "bars": false,
			  "dashLength": 10,
			  "dashes": false,
			  "fill": 1,
			  "gridPos": {
				"h": 5,
				"w": 5,
				"x": 0,
				"y": 35
			  },
			  "id": 46,
			  "legend": {
				"avg": false,
				"current": false,
				"max": false,
				"min": false,
				"show": true,
				"total": false,
				"values": false
			  },
			  "lines": true,
			  "linewidth": 1,
			  "links": [],
			  "nullPointMode": "null as zero",
			  "options": {},
			  "percentage": false,
			  "pointradius": 2,
			  "points": false,
			  "renderer": "flot",
			  "seriesOverrides": [],
			  "spaceLength": 10,
			  "stack": false,
			  "steppedLine": false,
			  "targets": [
				{
				  "expr": "up{namespace=\"$namespace\", job=\"keycloak-operator-metrics\", endpoint=\"http-metrics\"}",
				  "format": "time_series",
				  "hide": false,
				  "intervalFactor": 1,
				  "legendFormat": "Pod {{pod}}",
				  "refId": "A"
				}
			  ],
			  "thresholds": [],
			  "timeFrom": null,
			  "timeRegions": [],
			  "timeShift": null,
			  "title": "Operator Readiness Probes",
			  "tooltip": {
				"shared": true,
				"sort": 0,
				"value_type": "individual"
			  },
			  "type": "graph",
			  "xaxis": {
				"buckets": null,
				"mode": "time",
				"name": null,
				"show": true,
				"values": []
			  },
			  "yaxes": [
				{
				  "decimals": 0,
				  "format": "short",
				  "label": "Ready",
				  "logBase": 1,
				  "max": "1",
				  "min": "0",
				  "show": true
				},
				{
				  "format": "short",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": false
				}
			  ],
			  "yaxis": {
				"align": false,
				"alignLevel": null
			  }
			}
		  ],
		  "title": "Operator metrics",
		  "type": "row"
		},
		{
		  "collapsed": true,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 35
		  },
		  "id": 43,
		  "panels": [
			{
			  "aliasColors": {},
			  "bars": false,
			  "dashLength": 10,
			  "dashes": false,
			  "datasource": "Prometheus",
			  "fill": 1,
			  "gridPos": {
				"h": 5,
				"w": 6,
				"x": 0,
				"y": 36
			  },
			  "hideTimeOverride": false,
			  "id": 1,
			  "legend": {
				"alignAsTable": true,
				"avg": false,
				"current": false,
				"max": false,
				"min": false,
				"rightSide": true,
				"show": true,
				"sideWidth": 100,
				"total": false,
				"values": false
			  },
			  "lines": true,
			  "linewidth": 1,
			  "links": [],
			  "nullPointMode": "connected",
			  "options": {},
			  "percentage": false,
			  "pointradius": 5,
			  "points": false,
			  "renderer": "flot",
			  "seriesOverrides": [],
			  "spaceLength": 10,
			  "stack": false,
			  "steppedLine": false,
			  "targets": [
				{
				  "expr": "sum by (realm)(increase(keycloak_logins{namespace=\"$namespace\",job=\"keycloak\"}[1h]))",
				  "format": "time_series",
				  "hide": false,
				  "interval": "",
				  "intervalFactor": 1,
				  "legendFormat": "{{realm}}",
				  "refId": "A"
				}
			  ],
			  "thresholds": [],
			  "timeFrom": null,
			  "timeRegions": [],
			  "timeShift": null,
			  "title": "Logins per REALM",
			  "tooltip": {
				"shared": true,
				"sort": 0,
				"value_type": "individual"
			  },
			  "type": "graph",
			  "xaxis": {
				"buckets": null,
				"mode": "time",
				"name": null,
				"show": true,
				"values": []
			  },
			  "yaxes": [
				{
				  "decimals": 0,
				  "format": "none",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": "0",
				  "show": true
				},
				{
				  "format": "short",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": false
				}
			  ],
			  "yaxis": {
				"align": false,
				"alignLevel": null
			  }
			},
			{
			  "aliasColors": {},
			  "bars": false,
			  "dashLength": 10,
			  "dashes": false,
			  "datasource": "Prometheus",
			  "fill": 1,
			  "gridPos": {
				"h": 5,
				"w": 6,
				"x": 6,
				"y": 36
			  },
			  "hideTimeOverride": false,
			  "id": 7,
			  "legend": {
				"alignAsTable": true,
				"avg": false,
				"current": true,
				"hideEmpty": false,
				"hideZero": true,
				"max": false,
				"min": false,
				"rightSide": true,
				"show": true,
				"sideWidth": null,
				"sort": "current",
				"sortDesc": false,
				"total": false,
				"values": true
			  },
			  "lines": true,
			  "linewidth": 1,
			  "links": [],
			  "nullPointMode": "connected",
			  "options": {},
			  "percentage": false,
			  "pointradius": 5,
			  "points": false,
			  "renderer": "flot",
			  "seriesOverrides": [],
			  "spaceLength": 10,
			  "stack": false,
			  "steppedLine": false,
			  "targets": [
				{
				  "expr": "increase(keycloak_failed_login_attempts{namespace=\"$namespace\",job=\"keycloak\",provider=\"keycloak\",realm=\"$realm\"}[1h])",
				  "format": "time_series",
				  "hide": false,
				  "instant": false,
				  "interval": "",
				  "intervalFactor": 1,
				  "legendFormat": "{{error}}",
				  "refId": "A"
				},
				{
				  "expr": "increase(keycloak_failed_login_attempts{namespace=\"$namespace\",job=\"keycloak\",provider=\"keycloak\"}[1h])",
				  "format": "time_series",
				  "intervalFactor": 1,
				  "refId": "B"
				}
			  ],
			  "thresholds": [],
			  "timeFrom": null,
			  "timeRegions": [],
			  "timeShift": null,
			  "title": "Login Errors on realm $realm",
			  "tooltip": {
				"shared": true,
				"sort": 0,
				"value_type": "individual"
			  },
			  "type": "graph",
			  "xaxis": {
				"buckets": null,
				"mode": "time",
				"name": null,
				"show": true,
				"values": []
			  },
			  "yaxes": [
				{
				  "decimals": null,
				  "format": "none",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": "0",
				  "show": true
				},
				{
				  "format": "short",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": false
				}
			  ],
			  "yaxis": {
				"align": false,
				"alignLevel": null
			  }
			},
			{
			  "aliasColors": {},
			  "bars": false,
			  "dashLength": 10,
			  "dashes": false,
			  "datasource": "Prometheus",
			  "fill": 1,
			  "gridPos": {
				"h": 5,
				"w": 6,
				"x": 0,
				"y": 41
			  },
			  "hideTimeOverride": false,
			  "id": 18,
			  "legend": {
				"alignAsTable": true,
				"avg": false,
				"current": false,
				"max": false,
				"min": false,
				"rightSide": true,
				"show": true,
				"sideWidth": 100,
				"total": false,
				"values": false
			  },
			  "lines": true,
			  "linewidth": 1,
			  "links": [],
			  "nullPointMode": "connected",
			  "options": {},
			  "percentage": false,
			  "pointradius": 5,
			  "points": false,
			  "renderer": "flot",
			  "seriesOverrides": [],
			  "spaceLength": 10,
			  "stack": false,
			  "steppedLine": false,
			  "targets": [
				{
				  "expr": "increase(keycloak_logins{namespace=\"$namespace\",job=\"keycloak\",realm=\"$realm\",provider=\"keycloak\"}[1h])",
				  "format": "time_series",
				  "hide": false,
				  "interval": "",
				  "intervalFactor": 2,
				  "legendFormat": "{{client_id}}",
				  "refId": "A"
				}
			  ],
			  "thresholds": [],
			  "timeFrom": null,
			  "timeRegions": [],
			  "timeShift": null,
			  "title": "Logins per CLIENT on realm $realm",
			  "tooltip": {
				"shared": true,
				"sort": 0,
				"value_type": "individual"
			  },
			  "type": "graph",
			  "xaxis": {
				"buckets": null,
				"mode": "time",
				"name": null,
				"show": true,
				"values": []
			  },
			  "yaxes": [
				{
				  "decimals": 0,
				  "format": "none",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": "0",
				  "show": true
				},
				{
				  "format": "short",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": false
				}
			  ],
			  "yaxis": {
				"align": false,
				"alignLevel": null
			  }
			},
			{
			  "cacheTimeout": null,
			  "cards": {
				"cardPadding": null,
				"cardRound": null
			  },
			  "color": {
				"cardColor": "#b4ff00",
				"colorScale": "sqrt",
				"colorScheme": "interpolateBlues",
				"exponent": 0.5,
				"mode": "opacity"
			  },
			  "dataFormat": "timeseries",
			  "description": "",
			  "gridPos": {
				"h": 8,
				"w": 9,
				"x": 0,
				"y": 46
			  },
			  "heatmap": {},
			  "hideZeroBuckets": false,
			  "highlightCards": true,
			  "id": 29,
			  "legend": {
				"show": false
			  },
			  "links": [],
			  "options": {},
			  "pluginVersion": "6.2.4",
			  "reverseYBuckets": false,
			  "targets": [
				{
				  "expr": "sum(increase(keycloak_request_duration_bucket[5m])) by (le)",
				  "format": "time_series",
				  "instant": false,
				  "intervalFactor": 1,
				  "legendFormat": "Request taking less or equal than {{le}}",
				  "refId": "A"
				}
			  ],
			  "timeFrom": null,
			  "timeShift": null,
			  "title": "Endpoint latency [5m]",
			  "tooltip": {
				"show": true,
				"showHistogram": true
			  },
			  "type": "heatmap",
			  "xAxis": {
				"show": true
			  },
			  "xBucketNumber": null,
			  "xBucketSize": null,
			  "yAxis": {
				"decimals": null,
				"format": "ms",
				"logBase": 1,
				"max": null,
				"min": "0",
				"show": true,
				"splitFactor": null
			  },
			  "yBucketBound": "auto",
			  "yBucketNumber": null,
			  "yBucketSize": null
			},
			{
			  "aliasColors": {},
			  "bars": true,
			  "cacheTimeout": null,
			  "dashLength": 10,
			  "dashes": false,
			  "description": "",
			  "fill": 1,
			  "gridPos": {
				"h": 8,
				"w": 24,
				"x": 0,
				"y": 54
			  },
			  "id": 47,
			  "legend": {
				"avg": false,
				"current": false,
				"max": false,
				"min": false,
				"show": false,
				"total": false,
				"values": false
			  },
			  "lines": false,
			  "linewidth": 1,
			  "links": [],
			  "nullPointMode": "null",
			  "options": {},
			  "percentage": false,
			  "pluginVersion": "6.2.4",
			  "pointradius": 2,
			  "points": false,
			  "renderer": "flot",
			  "seriesOverrides": [],
			  "spaceLength": 10,
			  "stack": false,
			  "steppedLine": false,
			  "targets": [
				{
				  "expr": "sum(increase(keycloak_request_duration_sum{namespace=\"$namespace\"}[5m])) by (route) / sum(increase(keycloak_request_duration_count{namespace=\"$namespace\"}[5m])) by (route)",
				  "format": "time_series",
				  "instant": false,
				  "intervalFactor": 1,
				  "legendFormat": "{{route}}",
				  "refId": "A"
				}
			  ],
			  "thresholds": [],
			  "timeFrom": null,
			  "timeRegions": [],
			  "timeShift": null,
			  "title": "Endpoint latency [5m]",
			  "tooltip": {
				"shared": false,
				"sort": 0,
				"value_type": "individual"
			  },
			  "type": "graph",
			  "xaxis": {
				"buckets": null,
				"mode": "series",
				"name": null,
				"show": true,
				"values": [
				  "total"
				]
			  },
			  "yaxes": [
				{
				  "format": "ms",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": true
				},
				{
				  "format": "short",
				  "label": null,
				  "logBase": 1,
				  "max": null,
				  "min": null,
				  "show": true
				}
			  ],
			  "yaxis": {
				"align": false,
				"alignLevel": null
			  }
			}
		  ],
		  "title": "Experimental",
		  "type": "row"
		}
	  ],
	  "refresh": false,
	  "schemaVersion": 18,
	  "style": "dark",
	  "tags": [],
	  "templating": {
		"list": [
		  {
			"allValue": null,
			"current": {
			  "tags": [],
			  "text": "All",
			  "value": "$__all"
			},
			"datasource": "Prometheus",
			"definition": "label_values(keycloak_request_duration_count,namespace)",
			"hide": 0,
			"includeAll": true,
			"label": "Namespace",
			"multi": false,
			"name": "namespace",
			"options": [],
			"query": "label_values(keycloak_request_duration_count,namespace)",
			"refresh": 1,
			"regex": "",
			"skipUrlSync": false,
			"sort": 0,
			"tagValuesQuery": "",
			"tags": [],
			"tagsQuery": "",
			"type": "query",
			"useTags": false
		  },
		  {
			"allValue": null,
			"current": {
			  "tags": [],
			  "text": "All",
			  "value": "$__all"
			},
			"datasource": "Prometheus",
			"definition": "",
			"hide": 0,
			"includeAll": true,
			"label": "Realm",
			"multi": false,
			"name": "realm",
			"options": [],
			"query": "label_values(keycloak_logins{namespace=\"$namespace\",job=\"keycloak\",provider=\"keycloak\"},realm)",
			"refresh": 1,
			"regex": "",
			"skipUrlSync": false,
			"sort": 0,
			"tagValuesQuery": "",
			"tags": [],
			"tagsQuery": "",
			"type": "query",
			"useTags": false
		  },
		  {
			"allValue": null,
			"current": {
			  "tags": [],
			  "text": "All",
			  "value": "$__all"
			},
			"datasource": "Prometheus",
			"definition": "label_values(keycloak_logins{namespace=\"$namespace\",job=\"keycloak\",provider=\"keycloak\",realm=\"$realm\"},client_id)",
			"hide": 0,
			"includeAll": true,
			"label": "ClientId",
			"multi": false,
			"name": "ClientId",
			"options": [],
			"query": "label_values(keycloak_logins{namespace=\"$namespace\",job=\"keycloak\",provider=\"keycloak\",realm=\"$realm\"},client_id)",
			"refresh": 1,
			"regex": "",
			"skipUrlSync": false,
			"sort": 0,
			"tagValuesQuery": "",
			"tags": [],
			"tagsQuery": "",
			"type": "query",
			"useTags": false
		  }
		]
	  },
	  "time": {
		"from": "now-15m",
		"to": "now"
	  },
	  "timepicker": {
		"hidden": false,
		"refresh_intervals": [
		  "5s",
		  "10s",
		  "30s",
		  "1m",
		  "5m",
		  "15m",
		  "30m",
		  "1h",
		  "2h",
		  "1d"
		],
		"time_options": [
		  "5m",
		  "15m",
		  "1h",
		  "6h",
		  "12h",
		  "24h",
		  "2d",
		  "7d",
		  "30d"
		]
	  },
	  "title": "Keycloak Metrics",
	  "version": 3
    }`
