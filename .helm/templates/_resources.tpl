{{- define "app_resources" }}
resources:
  requests:
    cpu: {{ pluck .Values.werf.env .Values.app.resources.cpu | first | default .Values.app.resources.cpu._default }}
    memory: {{ pluck .Values.werf.env .Values.app.resources.memory.requests | first | default .Values.app.resources.memory.requests._default }}
  limits:
    memory: {{ pluck .Values.werf.env .Values.app.resources.memory.limits | first | default .Values.app.resources.memory.limits._default }}
{{- end }}