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
{{- end }}

{{- define "affinity_consumers" }}

{{- $ := index . 0 }}
{{- $consumername := index . 1 }}

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 5
      podAffinityTerm:
        topologyKey: "kubernetes.io/hostname"
        labelSelector:
          matchLabels:
             app: {{ $.Chart.Name }}-consumer-{{ $consumername | replace "_" "-" }}
{{- end }}
