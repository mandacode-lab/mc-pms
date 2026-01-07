{{/*
HTTP ConfigMap Name
Usage: {{ include "mc-pms.httpConfigMapName" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.httpConfigMapName" -}}
{{- printf "%s-%s-http" (include "mc-pms.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Database ConfigMap Name
Usage: {{ include "mc-pms.dbConfigMapName" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.dbConfigMapName" -}}
{{- printf "%s-%s-db" (include "mc-pms.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
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
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "mc-pms.httpConfigMapName" (dict "root" $root "component" $component) }}
  labels:
    {{- include "mc-pms.componentLabels" (dict "root" $root "component" $component) | nindent 4 }}
data:
  {{ $prefix }}HOST: {{ $http.host | default "0.0.0.0" | quote }}
  {{ $prefix }}PORT: {{ $http.port | default 8080 | quote }}
  {{ $prefix }}READ_TIMEOUT: {{ $http.readTimeout | default 60 | quote }}
  {{ $prefix }}WRITE_TIMEOUT: {{ $http.writeTimeout | default 60 | quote }}
  {{ $prefix }}IDLE_TIMEOUT: {{ $http.idleTimeout | default 120 | quote }}
  {{ $prefix }}SHUTDOWN_TIMEOUT: {{ $http.shutdownTimeout | default 20 | quote }}
  {{- $cors := $http.cors | default dict }}
  {{- $corsPrefix := printf "%sCORS_" $prefix }}
  {{ $corsPrefix }}ENABLED: {{ $cors.enabled | default false | quote }}
  {{- if $cors.allowedOrigins }}
  {{ $corsPrefix }}ALLOWED_ORIGINS: {{ $cors.allowedOrigins | join "," | quote }}
  {{- end }}
  {{- if $cors.allowedMethods }}
  {{ $corsPrefix }}ALLOWED_METHODS: {{ $cors.allowedMethods | join "," | quote }}
  {{- end }}
  {{- if $cors.allowedHeaders }}
  {{ $corsPrefix }}ALLOWED_HEADERS: {{ $cors.allowedHeaders | join "," | quote }}
  {{- end }}
  {{- if $cors.exposeHeaders }}
  {{ $corsPrefix }}EXPOSE_HEADERS: {{ $cors.exposeHeaders | join "," | quote }}
  {{- end }}
  {{ $corsPrefix }}ALLOW_CREDENTIALS: {{ $cors.allowCredentials | default false | quote }}
  {{ $corsPrefix }}MAX_AGE: {{ $cors.maxAge | default 43200 | quote }}
{{- end }}

{{/*
Database ConfigMap Template
Creates a ConfigMap with database configuration (non-sensitive data) using prefix
Usage: {{ include "mc-pms.dbConfigMap" (dict "root" $ "component" "admin" "prefix" "DB_") }}
*/}}
{{- define "mc-pms.dbConfigMap" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $root := .root -}}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "mc-pms.dbConfigMapName" (dict "root" $root "component" $component) }}
  labels:
    {{- include "mc-pms.componentLabels" (dict "root" $root "component" $component) | nindent 4 }}
data:
  {{ $prefix }}HOST: {{ $root.Values.database.host | quote }}
  {{ $prefix }}PORT: {{ $root.Values.database.port | quote }}
  {{ $prefix }}DATABASE: {{ $root.Values.database.name | quote }}
  {{ $prefix }}USERNAME: {{ $root.Values.database.username | quote }}
  {{ $prefix }}SSL_MODE: {{ $root.Values.database.sslMode | quote }}
{{- end }}
