{{/*
HTTP ConfigMap Name Generator
Returns the ConfigMap name for HTTP configuration
Usage: {{ include "pms-lib.configs.http.configMapName" (dict "component" "admin" "root" $) }}
*/}}
{{- define "pms-lib.configs.http.configMapName" -}}
{{- printf "%s-%s-http-config" (include "pms-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
HTTP ConfigMap Generator
Generates ConfigMap with HTTP configuration using prefix pattern
Automatically merges global values from root context

Usage: {{ include "pms-lib.config.http" (dict
  "component" "admin"
  "prefix" "HTTP_"
  "config" .Values
  "root" $
)}}
*/}}
{{- define "pms-lib.config.http" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- $global := $root.Values.global.http | default dict -}}
{{- $httpConfig := $config.http | default dict -}}

{{- if $config.enabled }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "pms-lib.configs.http.configMapName" (dict "component" $component "root" $root) }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
    app.kubernetes.io/config: http
data:
  {{ $prefix }}HOST: {{ $httpConfig.host | default $global.host | default "0.0.0.0" | quote }}
  {{ $prefix }}PORT: {{ $httpConfig.port | default $global.port | default "8080" | quote }}
  {{ $prefix }}READ_TIMEOUT: {{ $httpConfig.readTimeout | default $global.readTimeout | default "10s" | quote }}
  {{ $prefix }}WRITE_TIMEOUT: {{ $httpConfig.writeTimeout | default $global.writeTimeout | default "10s" | quote }}
  {{ $prefix }}IDLE_TIMEOUT: {{ $httpConfig.idleTimeout | default $global.idleTimeout | default "60s" | quote }}
  {{ $prefix }}SHUTDOWN_TIMEOUT: {{ $httpConfig.shutdownTimeout | default $global.shutdownTimeout | default "5s" | quote }}
  {{- $corsConfig := $httpConfig.cors | default dict -}}
  {{- $globalCors := $global.cors | default dict -}}
  {{- $corsPrefix := printf "%sCORS_" $prefix }}
  {{ $corsPrefix }}ENABLED: {{ $corsConfig.enabled | default $globalCors.enabled | default false | quote }}
  {{ $corsPrefix }}ALLOWED_ORIGINS: {{ $corsConfig.allowedOrigins | default $globalCors.allowedOrigins | default (list "*") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_METHODS: {{ $corsConfig.allowedMethods | default $globalCors.allowedMethods | default (list "GET" "POST" "PUT" "PATCH" "DELETE" "OPTIONS") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_HEADERS: {{ $corsConfig.allowedHeaders | default $globalCors.allowedHeaders | default (list "Origin" "Content-Length" "Content-Type" "Authorization" "Accept" "X-Requested-With") | join "," | quote }}
  {{- if or $corsConfig.exposeHeaders $globalCors.exposeHeaders }}
  {{ $corsPrefix }}EXPOSE_HEADERS: {{ $corsConfig.exposeHeaders | default $globalCors.exposeHeaders | join "," | quote }}
  {{- end }}
  {{ $corsPrefix }}ALLOW_CREDENTIALS: {{ $corsConfig.allowCredentials | default $globalCors.allowCredentials | default false | quote }}
  {{ $corsPrefix }}MAX_AGE: {{ $corsConfig.maxAge | default $globalCors.maxAge | default 43200 | quote }}
{{- end }}
{{- end }}
