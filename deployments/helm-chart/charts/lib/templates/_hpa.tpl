{{/*
Generic HorizontalPodAutoscaler template
Usage: {{ include "pms-lib.hpa" (dict
  "component" "admin"
  "config" .Values.admin
  "root" $
)}}
*/}}
{{- define "pms-lib.hpa" -}}
{{- $component := .component -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- if and $config.enabled $config.autoscaling.enabled }}
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "pms-lib.fullname" $root }}-{{ $component }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "pms-lib.fullname" $root }}-{{ $component }}
  minReplicas: {{ $config.autoscaling.minReplicas }}
  maxReplicas: {{ $config.autoscaling.maxReplicas }}
  metrics:
  {{- if $config.autoscaling.targetCPUUtilizationPercentage }}
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: {{ $config.autoscaling.targetCPUUtilizationPercentage }}
  {{- end }}
  {{- if $config.autoscaling.targetMemoryUtilizationPercentage }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ $config.autoscaling.targetMemoryUtilizationPercentage }}
  {{- end }}
{{- end }}
{{- end }}
