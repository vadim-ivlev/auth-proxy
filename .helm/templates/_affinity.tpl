{{- define "affinity" }}

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 5
      podAffinityTerm:
        topologyKey: "kubernetes.io/hostname"
        labelSelector:
          matchLabels:
             app: {{ .Chart.Name }}
{{- if eq .Values.werf.env "production" }}
  nodeAffinity:
{{- toYaml (pluck .Values.werf.env .Values.nodeaffinity | first | default .Values.nodeaffinity._default) | nindent 8 }}
{{- end }}
{{- end }}
