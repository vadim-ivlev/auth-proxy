{{- define "app_probes" }}
livenessProbe:
  httpGet:
    path: /ping
    port: {{ pluck .Values.werf.env .Values.app.port | first | default .Values.app.port._default }}
readinessProbe:
  httpGet:
    path: /ping
    port: {{ pluck .Values.werf.env .Values.app.port | first | default .Values.app.port._default }}
  failureThreshold: 1
  initialDelaySeconds: 3
{{- end }}