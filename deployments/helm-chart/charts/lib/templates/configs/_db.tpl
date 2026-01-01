{{/*
Database ConfigMap Name Generator
Returns the ConfigMap name for database configuration
Usage: {{ include "pms-lib.configs.db.configMapName" (dict "component" "admin" "root" $) }}
*/}}
{{- define "pms-lib.configs.db.configMapName" -}}
{{- printf "%s-%s-db-config" (include "pms-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Database Secret Name Generator
Returns the Secret name for database configuration
Automatically uses existingSecret if provided (component-level takes precedence over global)
Automatically reads global from root context

Usage: {{ include "pms-lib.configs.db.secretName" (dict "component" "admin" "config" .Values "root" $) }}
*/}}
{{- define "pms-lib.configs.db.secretName" -}}
{{- $dbConfig := .config.database | default dict -}}
{{- $globalDb := .root.Values.global.database | default dict -}}
{{- $existingSecret := $dbConfig.existingSecret | default $globalDb.existingSecret | default "" -}}
{{- if $existingSecret -}}
{{- $existingSecret -}}
{{- else -}}
{{- printf "%s-%s-db-secret" (include "pms-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end }}

{{/*
Database ConfigMap & Secret Generator
Generates ConfigMap and Secret with database configuration using prefix pattern
Automatically merges global values from root context

Usage: {{ include "pms-lib.config.db" (dict
  "component" "admin"
  "prefix" "DB_"
  "config" .Values
  "root" $
)}}
*/}}
{{- define "pms-lib.config.db" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- $global := $root.Values.global.database | default dict -}}
{{- $dbConfig := $config.database | default dict -}}
{{- $existingSecret := $dbConfig.existingSecret | default $global.existingSecret | default "" -}}
{{- $createSecret := and (not $existingSecret) ($dbConfig.enabled | default false) -}}

{{- if $dbConfig.enabled }}
---
# ConfigMap for non-sensitive database configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "pms-lib.configs.db.configMapName" (dict "component" $component "root" $root) }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
    app.kubernetes.io/config: database
data:
  {{ $prefix }}HOST: {{ $dbConfig.host | default $global.host | default "localhost" | quote }}
  {{ $prefix }}PORT: {{ $dbConfig.port | default $global.port | default "5432" | quote }}
  {{ $prefix }}DATABASE: {{ $dbConfig.database | default $global.database | default "postgres" | quote }}
  {{ $prefix }}USERNAME: {{ $dbConfig.username | default $global.username | default "postgres" | quote }}
  {{ $prefix }}SSL_MODE: {{ $dbConfig.sslMode | default $global.sslMode | default "disable" | quote }}

{{- if $createSecret }}
---
# Secret for sensitive database configuration
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "pms-lib.configs.db.secretName" (dict "component" $component "config" $config "root" $root) }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
    app.kubernetes.io/config: database
type: Opaque
stringData:
  {{ $prefix }}PASSWORD: {{ $dbConfig.password | default $global.password | default "postgres" | quote }}
{{- end }}
{{- end }}
{{- end }}
