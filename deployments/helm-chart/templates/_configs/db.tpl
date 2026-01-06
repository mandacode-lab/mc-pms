{{/*
Database ConfigMap Name
Usage: {{ include "mc-pms.dbConfigMapName" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.dbConfigMapName" -}}
{{- printf "%s-%s-db" (include "mc-pms.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Database Secret Name
Returns existingSecret if provided, otherwise generates a secret name
Usage: {{ include "mc-pms.dbSecretName" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.dbSecretName" -}}
{{- if .root.Values.database.existingSecret }}
{{- .root.Values.database.existingSecret }}
{{- else }}
{{- printf "%s-%s-db" (include "mc-pms.fullname" .root) .component | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Database ConfigMap Template
Creates a ConfigMap with database configuration (non-sensitive data only) using prefix
Usage: {{ include "mc-pms.dbConfig" (dict "root" $ "component" "admin" "prefix" "DB_") }}
*/}}
{{- define "mc-pms.dbConfig" -}}
{{- $component := .component -}}
{{- $prefix := .prefix -}}
{{- $root := .root }}
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
  {{ $prefix }}SSL_MODE: {{ $root.Values.database.sslMode | quote }}
{{- end }}
