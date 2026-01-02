{{/*
Migration ConfigMap Name
Usage: {{ include "mc-pms.migrationConfigMapName" . }}
*/}}
{{- define "mc-pms.migrationConfigMapName" -}}
{{- printf "%s-migration" (include "mc-pms.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Migration ConfigMap Template
Creates a ConfigMap with migration configuration (non-sensitive data only)
Usage: {{ include "mc-pms.migrationConfig" . }}
*/}}
{{- define "mc-pms.migrationConfig" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "mc-pms.migrationConfigMapName" . }}
  labels:
    {{- include "mc-pms.componentLabels" (dict "root" . "component" "migration") | nindent 4 }}
data:
  # Database connection parameters (non-sensitive)
  HOST: {{ include "mc-pms.database.host" . | quote }}
  PORT: {{ include "mc-pms.database.port" . | quote }}
  DB_NAME: {{ include "mc-pms.database.name" . | quote }}
  SSL_MODE: {{ include "mc-pms.database.sslMode" . | quote }}

  # Migration options
  {{- if .Values.migration.allowDirty }}
  ALLOW_DIRTY: "true"
  {{- end }}
  {{- if .Values.migration.dryRun }}
  DRY_RUN: "true"
  {{- end }}
{{- end }}
