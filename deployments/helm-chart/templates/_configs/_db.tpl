{{/*
Database ConfigMap Name
*/}}
{{- define "mc-pms.configs.db.configMapName" -}}
{{- printf "%s-%s-db-config" (include "mc-helm-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Database Secret Name
Uses existingSecret if provided (component > global)
*/}}
{{- define "mc-pms.configs.db.secretName" -}}
{{- $global := .root.Values.global | default dict -}}
{{- $existingSecret := .config.database.existingSecret | default $global.existingSecret | default "" -}}
{{- if $existingSecret -}}
{{- $existingSecret -}}
{{- else -}}
{{- printf "%s-%s-db-secret" (include "mc-helm-lib.fullname" .root) .component | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end }}

{{/*
Database ConfigMap & Secret
*/}}
{{- define "mc-pms.config.db" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $config := .config -}}
{{- $root := .root -}}
{{- $global := $root.Values.global | default dict -}}
{{- $db := $config.database -}}
{{- $existingSecret := $db.existingSecret | default $global.existingSecret | default "" -}}

{{- if $db.enabled }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "mc-pms.configs.db.configMapName" (dict "component" $component "root" $root) }}
  labels:
    {{- include "mc-helm-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
data:
  {{ $prefix }}HOST: {{ $db.host | default $global.host | default "localhost" | quote }}
  {{ $prefix }}PORT: {{ $db.port | default $global.port | default "5432" | quote }}
  {{ $prefix }}DATABASE: {{ $db.database | default $global.database | default "postgres" | quote }}
  {{ $prefix }}USERNAME: {{ $db.username | default $global.username | default "postgres" | quote }}
  {{ $prefix }}SSL_MODE: {{ $db.sslMode | default $global.sslMode | default "disable" | quote }}

{{- if not $existingSecret }}
---
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "mc-pms.configs.db.secretName" (dict "component" $component "config" $config "root" $root) }}
  labels:
    {{- include "mc-helm-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
type: Opaque
stringData:
  {{ $prefix }}PASSWORD: {{ $db.password | default $global.password | default "postgres" | quote }}
{{- end }}
{{- end }}
{{- end }}
