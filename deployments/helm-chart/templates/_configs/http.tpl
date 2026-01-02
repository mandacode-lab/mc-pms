{{/*
HTTP ConfigMap Name
Usage: {{ include "mc-pms.httpConfigMapName" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.httpConfigMapName" -}}
{{- printf "%s-%s-http" (include "mc-pms.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
HTTP ConfigMap Template
Creates a ConfigMap with HTTP configuration using prefix
Usage: {{ include "mc-pms.httpConfigMap" (dict "root" $ "component" "admin" "prefix" "HTTP_" "config" .Values.admin) }}
*/}}
{{- define "mc-pms.httpConfigMap" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- $http := $config.http | default dict -}}
{{- $cors := $http.cors | default dict -}}
{{- $corsPrefix := printf "%sCORS_" $prefix -}}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "mc-pms.httpConfigMapName" (dict "root" $root "component" $component) }}
  labels:
    {{- include "mc-pms.componentLabels" (dict "root" $root "component" $component) | nindent 4 }}
data:
  # HTTP Server Configuration
  {{ $prefix }}HOST: {{ $http.host | default "0.0.0.0" | quote }}
  {{ $prefix }}PORT: {{ $http.port | default 8080 | quote }}
  {{ $prefix }}READ_TIMEOUT: {{ $http.readTimeout | default 60 | quote }}
  {{ $prefix }}WRITE_TIMEOUT: {{ $http.writeTimeout | default 60 | quote }}
  {{ $prefix }}IDLE_TIMEOUT: {{ $http.idleTimeout | default 120 | quote }}
  {{ $prefix }}SHUTDOWN_TIMEOUT: {{ $http.shutdownTimeout | default 20 | quote }}

  # CORS Configuration (nested with CORS_ prefix)
  {{ $corsPrefix }}ENABLED: {{ $cors.enabled | default false | quote }}
  {{ $corsPrefix }}ALLOWED_ORIGINS: {{ $cors.allowedOrigins | default (list "*") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_METHODS: {{ $cors.allowedMethods | default (list "GET" "POST" "PUT" "PATCH" "DELETE" "OPTIONS") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_HEADERS: {{ $cors.allowedHeaders | default (list "Origin" "Content-Length" "Content-Type" "Authorization" "Accept" "X-Requested-With") | join "," | quote }}
  {{- if $cors.exposeHeaders }}
  {{ $corsPrefix }}EXPOSE_HEADERS: {{ $cors.exposeHeaders | join "," | quote }}
  {{- else }}
  {{ $corsPrefix }}EXPOSE_HEADERS: ""
  {{- end }}
  {{ $corsPrefix }}ALLOW_CREDENTIALS: {{ $cors.allowCredentials | default false | quote }}
  {{ $corsPrefix }}MAX_AGE: {{ $cors.maxAge | default 43200 | quote }}
{{- end }}
