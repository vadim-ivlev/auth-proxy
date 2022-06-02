# vi:syntax=yaml
# vi:filetype=yaml

{{- define "tolerations" }}
tolerations:
{{- toYaml (pluck .Values.werf.env .Values.tolerations | first | default .Values.tolerations._default) | nindent 8 }}
{{- end }}
