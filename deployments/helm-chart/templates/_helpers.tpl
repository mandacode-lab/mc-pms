{{/*
Expand the name of the chart.
*/}}
{{- define "mc-pms.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mc-pms.fullname" -}}
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
{{- define "mc-pms.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mc-pms.labels" -}}
helm.sh/chart: {{ include "mc-pms.chart" . }}
{{ include "mc-pms.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mc-pms.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mc-pms.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Component labels (add component name to labels)
Usage: {{ include "mc-pms.componentLabels" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.componentLabels" -}}
{{ include "mc-pms.labels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
Component selector labels
Usage: {{ include "mc-pms.componentSelectorLabels" (dict "root" $ "component" "admin") }}
*/}}
{{- define "mc-pms.componentSelectorLabels" -}}
{{ include "mc-pms.selectorLabels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "mc-pms.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "mc-pms.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database configuration helpers
*/}}
{{- define "mc-pms.database.host" -}}
{{- .Values.database.host | required "database.host is required" }}
{{- end }}

{{- define "mc-pms.database.port" -}}
{{- .Values.database.port | default 5432 }}
{{- end }}

{{- define "mc-pms.database.name" -}}
{{- .Values.database.name | required "database.name is required" }}
{{- end }}

{{- define "mc-pms.database.user" -}}
{{- .Values.database.user | required "database.user is required" }}
{{- end }}

{{- define "mc-pms.database.sslMode" -}}
{{- .Values.database.sslMode | default "require" }}
{{- end }}

{{/*
Global database secret name (for migration job)
Returns existingSecret if provided, otherwise generates a global secret name
*/}}
{{- define "mc-pms.database.secretName" -}}
{{- if .Values.database.existingSecret }}
{{- .Values.database.existingSecret }}
{{- else }}
{{- include "mc-pms.fullname" . }}-db
{{- end }}
{{- end }}

{{/*
Image helpers
Usage: {{ include "mc-pms.image" (dict "image" .Values.admin.image "chart" .Chart) }}
*/}}
{{- define "mc-pms.image" -}}
{{- $tag := .image.tag | default .chart.AppVersion }}
{{- printf "%s:%s" .image.repository $tag }}
{{- end }}

{{/*
HTTP configuration merge helper
Merges component-specific HTTP config with global config
Usage: {{ include "mc-pms.http.config" (dict "root" $ "component" .Values.admin) }}
*/}}
{{- define "mc-pms.http.config" -}}
{{- $global := .root.Values.global.http | default dict }}
{{- $component := .component.http | default dict }}
{{- $merged := merge $component $global }}
{{- toJson $merged }}
{{- end }}
