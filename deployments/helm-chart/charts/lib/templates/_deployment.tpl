{{/*
Generic Deployment template
Usage: {{ include "pms-lib.deployment" (dict
  "component" "admin"
  "config" .Values.admin
  "configRefs" (list
    (dict "type" "configMap" "name" "admin-http-config" "template" "pms-lib.config.http")
    (dict "type" "secret" "name" "admin-db-secret" "template" "pms-lib.config.db")
  )
  "root" $
)}}
*/}}
{{- define "pms-lib.deployment" -}}
{{- $component := .component -}}
{{- $config := .config -}}
{{- $configRefs := .configRefs | default list -}}
{{- $root := .root -}}
{{- if $config.enabled }}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "pms-lib.fullname" $root }}-{{ $component }}
  labels:
    {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 4 }}
spec:
  {{- if not $config.autoscaling.enabled }}
  replicas: {{ $config.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "pms-lib.selectorLabels" (dict "component" $component "root" $root) | nindent 6 }}
  template:
    metadata:
      annotations:
        {{- range $configRef := $configRefs }}
        checksum/{{ $configRef.name }}: {{ include $configRef.template $configRef.params | sha256sum }}
        {{- end }}
        {{- with $config.podAnnotations }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
      labels:
        {{- include "pms-lib.labels" (dict "component" $component "root" $root) | nindent 8 }}
        {{- with $config.podLabels }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      {{- with $root.Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "pms-lib.serviceAccountName" $root }}
      securityContext:
        {{- toYaml $config.podSecurityContext | nindent 8 }}
      containers:
      - name: {{ $component }}
        securityContext:
          {{- toYaml $config.securityContext | nindent 12 }}
        image: "{{ $config.image.repository }}:{{ $config.image.tag | default $root.Chart.AppVersion }}"
        imagePullPolicy: {{ $config.image.pullPolicy }}
        ports:
        - name: http
          containerPort: {{ $config.service.targetPort }}
          protocol: TCP
        {{- if $configRefs }}
        envFrom:
        {{- range $configRef := $configRefs }}
        - {{ $configRef.type }}Ref:
            name: {{ $configRef.name }}
        {{- end }}
        {{- end }}
        {{- with $config.env }}
        env:
          {{- toYaml . | nindent 10 }}
        {{- end }}
        livenessProbe:
          {{- toYaml $config.livenessProbe | nindent 10 }}
        readinessProbe:
          {{- toYaml $config.readinessProbe | nindent 10 }}
        resources:
          {{- toYaml $config.resources | nindent 10 }}
        {{- with $config.volumeMounts }}
        volumeMounts:
          {{- toYaml . | nindent 10 }}
        {{- end }}
      {{- with $config.volumes }}
      volumes:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with $config.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with $config.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with $config.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
{{- end }}
{{- end }}
