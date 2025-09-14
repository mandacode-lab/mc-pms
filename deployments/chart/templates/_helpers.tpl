{{/*
Expand the name of the chart.
*/}}
{{- define "defined.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "defined.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "defined.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "defined.labels" -}}
helm.sh/chart: {{ include "defined.chart" . }}
{{ include "defined.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "defined.selectorLabels" -}}
app.kubernetes.io/name: {{ include "defined.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Auth service labels
*/}}
{{- define "defined.auth.labels" -}}
{{ include "defined.labels" . }}
app.kubernetes.io/component: auth
{{- end }}

{{/*
Auth selector labels
*/}}
{{- define "defined.auth.selectorLabels" -}}
{{ include "defined.selectorLabels" . }}
app.kubernetes.io/component: auth
{{- end }}

{{/*
Core service labels
*/}}
{{- define "defined.core.labels" -}}
{{ include "defined.labels" . }}
app.kubernetes.io/component: core
{{- end }}

{{/*
Core selector labels
*/}}
{{- define "defined.core.selectorLabels" -}}
{{ include "defined.selectorLabels" . }}
app.kubernetes.io/component: core
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "defined.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "defined.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database Secret name
*/}}
{{- define "defined.databaseSecretName" -}}
{{- if .Values.database.existingSecret }}
{{- .Values.database.existingSecret }}
{{- else }}
{{- printf "%s-database" (include "defined.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Redis Secret Name
*/}}
{{- define "defined.redisSecretName" -}}
{{- if .Values.redis.existingSecret }}
{{- .Values.redis.existingSecret }}
{{- else }}
{{- printf "%s-redis" (include "defined.fullname" .) }}
{{- end }}
{{- end }}

{{/*
KEK Secret Name
*/}}
{{- define "defined.kekSecretName" -}}
{{- if .Values.kek.existingSecret }}
{{- .Values.kek.existingSecret }}
{{- else }}
{{- printf "%s-kek" (include "defined.fullname" .) }}
{{- end }}
{{- end }}

{{/*
ConfigMap names
*/}}
{{- define "defined.auth.configMapName" -}}
{{- printf "%s-auth-config" (include "defined.fullname" .) }}
{{- end }}

{{- define "defined.core.configMapName" -}}
{{- printf "%s-core-config" (include "defined.fullname" .) }}
{{- end }}