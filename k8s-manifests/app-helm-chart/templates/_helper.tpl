# common templates
{{- define "common.name" -}}
{{- printf "%s-%s" .Values.envPrefix .Values.appSlugName -}}
{{- end -}}

{{- define "secret.name" -}}
{{- printf "%s-secrets" (include "common.name" .) -}}
{{- end -}}

{{- define "deployment.name" -}}
{{- printf "%s-deploy" (include "common.name" .) -}}
{{- end -}}

{{- define "service.name" -}}
{{- printf "%s-svc" (include "common.name" .) -}}
{{- end -}}

{{- define "ingress.name" -}}
{{- printf "%s-ingress" (include "common.name" .) -}}
{{- end -}}

# app labels templates
{{- define "common.labels" -}}
env: {{ .Values.environment }}
app: {{ .Values.appName }}
version: {{ .Chart.AppVersion }}
chart-version: {{ .Chart.Version }}
product: {{ .Values.productName }}
{{- end -}}

# container templates
{{- define "container.resources" -}}
{{- if hasKey . "resources" -}}
resources:
{{- if hasKey .resources "requests" }}
  requests:
    memory: {{ .resources.requests.memory }}
    cpu: {{ .resources.requests.cpu }}
{{- end -}}
{{- if hasKey .resources "limits" }}
  limits:
    memory: {{ .resources.limits.memory }}
    cpu: {{ .resources.limits.cpu }}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "container.probes.readiness" -}}
{{- if hasKey . "readinessProbe" -}}
readinessProbe: {{- toYaml .readinessProbe | nindent 2 -}}
{{- end -}}
{{- end -}}

{{- define "container.probes.liveness" -}}
{{- if hasKey . "livenessProbe" -}}
livenessProbe: {{- toYaml .livenessProbe | nindent 2 -}}
{{- end -}}
{{- end -}}

{{- define "container.probes.startup" -}}
{{- if hasKey . "startupProbe" -}}
startupProbe: {{- toYaml .startupProbe | nindent 2 -}}
{{- end -}}
{{- end -}}

{{- define "container.probes" -}}
{{ include "container.probes.startup" . }}
{{ include "container.probes.liveness" . }}
{{ include "container.probes.readiness" . }}
{{- end -}}

# tekton resources related
{{- define "tekton.event-listener.name" -}}
{{- printf "%s-el" (include "common.name" .) -}}
{{- end -}}

{{- define "tekton.github-token-secret.name" -}}
{{- printf "%s-github-token-secret" (include "common.name" .) -}}
{{- end -}}

{{- define "tekton.tt.github.push.name" -}}
{{- printf "%s-github-push-tt" (include "common.name" .) -}}
{{- end -}}

{{- define "tekton.ingress.name" -}}
{{- printf "%s-webhook-ingress" (include "common.name" .) -}}
{{- end -}}

{{- define "tekton.ingress.path" -}}
{{- printf "/%s/%s" .Values.namespace (include "common.name" .) -}}
{{- end -}}
