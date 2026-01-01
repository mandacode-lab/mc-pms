{{/*
Generic Service template
Usage: {{ include "pms-lib.service" (dict
  "component" "admin"
  "config" .Values.admin
  "root" $
)}}
*/}}
{{- define "pms-lib.service" -}}
{{- $component := .component -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- if $config.enabled }}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "pms-lib.fullname" $root }}-{{ $component }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
  {{- with $config.service.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  type: {{ $config.service.type }}
  ports:
  - port: {{ $config.service.port }}
    targetPort: http
    protocol: TCP
    name: http
  selector:
    {{- include "pms-lib.selectorLabels" (dict "component" $component "root" $root) | nindent 4 }}
{{- end }}
{{- end }}
