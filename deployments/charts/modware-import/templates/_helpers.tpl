{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "modware-import.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "modware-import.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "modware-import.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Generate a random container name
*/}}
{{- define "modware-import.containerName" -}}
{{- randAlpha 10 | printf "%s-%s-%s" .Release.Name "importer" | trimSuffix "-" | lower -}}
{{- end -}}

{{/*
Generates image and imagePullPolicy manifest lines
*/}}
{{- define "modware-import.imageManifest" -}}
image: {{ .Values.image.repository }}:{{ .Values.image.tag}}
imagePullPolicy: {{ .Values.image.pullPolicy }}
{{- end -}}

{{/*
Generate env manifests
*/}}
{{- define "modware-import.envManifest" -}}
env:
- name: SECRET_KEY
  valueFrom:
    secretKeyRef:
       name: dictybase-configuration
       key: minio.secretkey
- name: ACCESS_KEY
  valueFrom:
    secretKeyRef:
       name: dictybase-configuration
       key: minio.accesskey
- name: LOG_LEVEL
  value: {{ .Values.log.level }}
{{- end -}}
