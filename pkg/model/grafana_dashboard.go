package model

import (
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GrafanaDashboard(cr *v1alpha1.Keycloak) *integreatlyv1alpha1.GrafanaDashboard {
	return &integreatlyv1alpha1.GrafanaDashboard{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": MonitoringKey,
			},
		},
		Spec: integreatlyv1alpha1.GrafanaDashboardSpec{
			Json: `{
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
				"editable": true,
				"gnetId": null,
				"graphTooltip": 0,
				"id": 3,
				"iteration": 1554910334349,
				"links": [],
				"panels": [
				{
					"cacheTimeout": null,
					"colorBackground": false,
					"colorValue": false,
					"colors": [
						"#299c46",
					"rgba(237, 129, 40, 0.89)",
					"#d44a3a"
				],
					"datasource": "Prometheus",
					"description": "All registered users in the application.",
					"format": "none",
					"gauge": {
					"maxValue": 100,
					"minValue": 0,
					"show": false,
					"thresholdLabels": false,
					"thresholdMarkers": true
				},
					"gridPos": {
					"h": 5,
					"w": 5,
					"x": 0,
					"y": 0
				},
					"hideTimeOverride": true,
					"id": 4,
					"interval": null,
					"links": [],
					"mappingType": 1,
					"mappingTypes": [
				{
					"name": "value to text",
					"value": 1
				},
				{
					"name": "range to text",
					"value": 2
				}
				],
					"maxDataPoints": 100,
					"nullPointMode": "connected",
					"nullText": null,
					"postfix": "",
					"postfixFontSize": "50%",
					"prefix": "",
					"prefixFontSize": "50%",
					"rangeMaps": [
				{
					"from": "null",
					"text": "0",
					"to": "null"
				}
				],
					"sparkline": {
					"fillColor": "rgba(31, 118, 189, 0.18)",
					"full": false,
					"lineColor": "rgb(31, 120, 193)",
					"show": false
				},
					"tableColumn": "",
					"targets": [
				{
					"expr": "sum(keycloak_registrations{namespace=\"$namespace\"})",
					"format": "time_series",
					"hide": false,
					"intervalFactor": 2,
					"legendFormat": "",
					"refId": "A"
				}
				],
					"thresholds": "",
					"timeFrom": "5s",
					"title": "Total Registrations",
					"transparent": false,
					"type": "singlestat",
					"valueFontSize": "110%",
					"valueMaps": [
				{
					"op": "=",
					"text": "N/A",
					"value": "null"
				}
				],
					"valueName": "current"
				},
				{
					"cacheTimeout": null,
					"colorBackground": false,
					"colorValue": false,
					"colors": [
						"#299c46",
					"#f9934e",
					"#d44a3a"
				],
					"datasource": "Prometheus",
					"decimals": null,
					"description": "All occurred log in events. This does not show current active users.",
					"format": "none",
					"gauge": {
					"maxValue": 100,
					"minValue": 0,
					"show": false,
					"thresholdLabels": false,
					"thresholdMarkers": true
				},
					"gridPos": {
					"h": 5,
					"w": 5,
					"x": 5,
					"y": 0
				},
					"hideTimeOverride": true,
					"id": 3,
					"interval": null,
					"links": [],
					"mappingType": 1,
					"mappingTypes": [
				{
					"name": "value to text",
					"value": 1
				},
				{
					"name": "range to text",
					"value": 2
				}
				],
					"maxDataPoints": 100,
					"nullPointMode": "connected",
					"nullText": null,
					"postfix": "",
					"postfixFontSize": "50%",
					"prefix": "",
					"prefixFontSize": "50%",
					"rangeMaps": [
				{
					"from": "null",
					"text": "N/A",
					"to": "null"
				}
				],
					"sparkline": {
					"fillColor": "rgba(31, 118, 189, 0.18)",
					"full": false,
					"lineColor": "rgb(31, 120, 193)",
					"show": false
				},
					"tableColumn": "",
					"targets": [
				{
					"expr": "sum(keycloak_logins{namespace=\"$namespace\"})",
					"format": "time_series",
					"hide": false,
					"intervalFactor": 2,
					"refId": "A"
				}
				],
					"thresholds": "",
					"timeFrom": "5s",
					"title": "Total Logins",
					"type": "singlestat",
					"valueFontSize": "110%",
					"valueMaps": [
				{
					"op": "=",
					"text": "0",
					"value": "null"
				}
				],
					"valueName": "current"
				},
				{
					"cacheTimeout": null,
					"colorBackground": false,
					"colorValue": false,
					"colors": [
						"#299c46",
					"rgba(237, 129, 40, 0.89)",
					"#d44a3a"
				],
					"datasource": "Prometheus",
					"description": "All failed login attempts that have occurred.",
					"format": "none",
					"gauge": {
					"maxValue": 100,
					"minValue": 0,
					"show": false,
					"thresholdLabels": false,
					"thresholdMarkers": true
				},
					"gridPos": {
					"h": 5,
					"w": 5,
					"x": 10,
					"y": 0
				},
					"hideTimeOverride": true,
					"id": 6,
					"interval": null,
					"links": [],
					"mappingType": 1,
					"mappingTypes": [
				{
					"name": "value to text",
					"value": 1
				},
				{
					"name": "range to text",
					"value": 2
				}
				],
					"maxDataPoints": 100,
					"nullPointMode": "connected",
					"nullText": null,
					"postfix": "",
					"postfixFontSize": "50%",
					"prefix": "",
					"prefixFontSize": "50%",
					"rangeMaps": [
				{
					"from": "null",
					"text": "0",
					"to": "null"
				}
				],
					"sparkline": {
					"fillColor": "rgba(31, 118, 189, 0.18)",
					"full": false,
					"lineColor": "rgb(31, 120, 193)",
					"show": false
				},
					"tableColumn": "",
					"targets": [
				{
					"expr": "sum(keycloak_failed_login_attempts{namespace=\"$namespace\"})",
					"format": "time_series",
					"intervalFactor": 2,
					"refId": "A"
				}
				],
					"thresholds": "",
					"timeFrom": "5s",
					"title": "Total Login Errors",
					"transparent": false,
					"type": "singlestat",
					"valueFontSize": "110%",
					"valueMaps": [
				{
					"op": "=",
					"text": "N/A",
					"value": "null"
				}
				],
					"valueName": "current"
				},
				{
					"cacheTimeout": null,
					"colorBackground": false,
					"colorValue": false,
					"colors": [
						"#299c46",
					"rgba(237, 129, 40, 0.89)",
					"#d44a3a"
				],
					"datasource": "Prometheus",
					"decimals": 1,
					"description": "Memory currently being used by Keycloak.",
					"format": "bytes",
					"gauge": {
					"maxValue": 2067000000,
					"minValue": 0,
					"show": true,
					"thresholdLabels": false,
					"thresholdMarkers": false
				},
					"gridPos": {
					"h": 5,
					"w": 5,
					"x": 15,
					"y": 0
				},
					"hideTimeOverride": true,
					"id": 5,
					"interval": null,
					"links": [],
					"mappingType": 1,
					"mappingTypes": [
				{
					"name": "value to text",
					"value": 1
				},
				{
					"name": "range to text",
					"value": 2
				}
				],
					"maxDataPoints": 100,
					"nullPointMode": "connected",
					"nullText": null,
					"postfix": "",
					"postfixFontSize": "50%",
					"prefix": "",
					"prefixFontSize": "50%",
					"rangeMaps": [
				{
					"from": "null",
					"text": "N/A",
					"to": "null"
				}
				],
					"sparkline": {
					"fillColor": "rgba(31, 118, 189, 0.18)",
					"full": false,
					"lineColor": "rgb(31, 120, 193)",
					"show": false
				},
					"tableColumn": "Value",
					"targets": [
				{
					"expr": "max(jvm_memory_bytes_used{area=\"heap\",namespace=\"$namespace\", service=\"sso\"}) + max(jvm_memory_bytes_used{area=\"nonheap\", namespace=\"$namespace\", service=\"sso\"})",
					"format": "time_series",
					"hide": false,
					"instant": false,
					"intervalFactor": 2,
					"legendFormat": "",
					"refId": "B"
				}
				],
					"thresholds": "682e6, 1364e6",
					"timeFrom": "5s",
					"title": "Current Memory",
					"type": "singlestat",
					"valueFontSize": "50%",
					"valueMaps": [
				{
					"op": "=",
					"text": "N/A",
					"value": "null"
				}
				],
					"valueName": "current"
				},
				{
					"aliasColors": {},
					"bars": false,
					"dashLength": 10,
					"dashes": false,
					"datasource": "Prometheus",
					"fill": 1,
					"gridPos": {
						"h": 7,
						"w": 20,
						"x": 0,
						"y": 5
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
						"expr": "keycloak_logins{namespace=\"$namespace\"}",
						"format": "time_series",
						"hide": false,
						"interval": "",
						"intervalFactor": 2,
						"refId": "A"
					}
				],
					"thresholds": [],
					"timeFrom": "1h",
					"timeRegions": [],
					"timeShift": null,
					"title": "Logins",
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
						"h": 7,
						"w": 20,
						"x": 0,
						"y": 12
					},
					"hideTimeOverride": false,
					"id": 7,
					"legend": {
						"alignAsTable": true,
						"avg": false,
						"current": false,
						"hideEmpty": false,
						"max": false,
						"min": false,
						"rightSide": true,
						"show": true,
						"sideWidth": null,
						"total": false,
						"values": false
					},
					"lines": true,
					"linewidth": 1,
					"links": [],
					"nullPointMode": "connected",
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
						"expr": "keycloak_failed_login_attempts{namespace=\"$namespace\"}",
						"format": "time_series",
						"hide": false,
						"interval": "",
						"intervalFactor": 2,
						"refId": "A"
					}
				],
					"thresholds": [],
					"timeFrom": "1h",
					"timeRegions": [],
					"timeShift": null,
					"title": "Login Errors",
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
					"fill": 2,
					"gridPos": {
						"h": 7,
						"w": 10,
						"x": 0,
						"y": 19
					},
					"hideTimeOverride": false,
					"id": 12,
					"legend": {
						"avg": false,
						"current": false,
						"max": false,
						"min": false,
						"rightSide": true,
						"show": true,
						"sideWidth": 70,
						"total": false,
						"values": false
					},
					"lines": true,
					"linewidth": 1,
					"links": [],
					"nullPointMode": "connected",
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
						"expr": "sum(jvm_memory_bytes_max{namespace=\"$namespace\",service=\"sso\"})",
						"format": "time_series",
						"instant": false,
						"intervalFactor": 1,
						"legendFormat": "Max",
						"refId": "A"
					},
					{
						"expr": "sum(jvm_memory_bytes_committed{namespace=\"$namespace\",service=\"sso\"})",
						"format": "time_series",
						"intervalFactor": 1,
						"legendFormat": "Committed",
						"refId": "C"
					},
					{
						"expr": "sum(jvm_memory_bytes_used{namespace=\"$namespace\",service=\"sso\"})",
						"format": "time_series",
						"instant": false,
						"intervalFactor": 1,
						"legendFormat": "Used",
						"refId": "B"
					}
				],
					"thresholds": [],
					"timeFrom": "1h",
					"timeRegions": [],
					"timeShift": null,
					"title": "Memory Usage",
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
				"refresh": "5s",
				"schemaVersion": 16,
				"style": "dark",
				"tags": [],
				"templating": {
				"list": [
			{
				"allValue": null,
				"current": {
				"tags": [],
				"text": "sso",
				"value": "sso"
			},
				"datasource": "Prometheus",
				"definition": "label_values(jvm_memory_bytes_used{service=\"sso\"}, namespace)",
				"hide": 0,
				"includeAll": false,
				"label": "Namespace",
				"multi": false,
				"name": "namespace",
				"options": [],
				"query": "label_values(jvm_memory_bytes_used{service=\"sso\"}, namespace)",
				"refresh": 1,
				"regex": "",
				"skipUrlSync": false,
				"sort": 1,
				"tagValuesQuery": "",
				"tags": [],
				"tagsQuery": "",
				"type": "query",
				"useTags": false
			}
			]
			},
				"time": {
				"from": "now/d",
				"to": "now"
			},
				"timepicker": {
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
				"timezone": "",
				"title": "Keycloak"
			}`,
			Name: "keycloak.json",
		},
	}
}

func GrafanaDashboardSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
