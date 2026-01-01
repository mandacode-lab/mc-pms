{{/*
HTTP ConfigMap Name
*/}}
{{- define "lib.configs.http.configMapName" -}}
{{- printf "%s-%s-http-config" (include "mc-helm-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
HTTP ConfigMap
*/}}
{{- define "lib.config.http" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- $global := $root.Values.global | default dict -}}
{{- $http := $config.http -}}

{{- if $config.enabled }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "lib.configs.http.configMapName" (dict "component" $component "root" $root) }}
  labels:
    {{- include "mc-helm-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
data:
  {{ $prefix }}HOST: {{ $http.host | default $global.host | default "0.0.0.0" | quote }}
  {{ $prefix }}PORT: {{ $http.port | default $global.port | default "8080" | quote }}
  {{ $prefix }}READ_TIMEOUT: {{ $http.readTimeout | default $global.readTimeout | default "10s" | quote }}
  {{ $prefix }}WRITE_TIMEOUT: {{ $http.writeTimeout | default $global.writeTimeout | default "10s" | quote }}
  {{ $prefix }}IDLE_TIMEOUT: {{ $http.idleTimeout | default $global.idleTimeout | default "60s" | quote }}
  {{ $prefix }}SHUTDOWN_TIMEOUT: {{ $http.shutdownTimeout | default $global.shutdownTimeout | default "5s" | quote }}
  {{- $cors := $http.cors | default dict }}
  {{- $corsPrefix := printf "%sCORS_" $prefix }}
  {{ $corsPrefix }}ENABLED: {{ $cors.enabled | default false | quote }}
  {{ $corsPrefix }}ALLOWED_ORIGINS: {{ $cors.allowedOrigins | default (list "*") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_METHODS: {{ $cors.allowedMethods | default (list "GET" "POST" "PUT" "PATCH" "DELETE" "OPTIONS") | join "," | quote }}
  {{ $corsPrefix }}ALLOWED_HEADERS: {{ $cors.allowedHeaders | default (list "Origin" "Content-Length" "Content-Type" "Authorization" "Accept" "X-Requested-With") | join "," | quote }}
  {{- if $cors.exposeHeaders }}
  {{ $corsPrefix }}EXPOSE_HEADERS: {{ $cors.exposeHeaders | join "," | quote }}
  {{- end }}
  {{ $corsPrefix }}ALLOW_CREDENTIALS: {{ $cors.allowCredentials | default false | quote }}
  {{ $corsPrefix }}MAX_AGE: {{ $cors.maxAge | default 43200 | quote }}
{{- end }}
{{- end }}
