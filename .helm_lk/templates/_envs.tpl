{{- define "app_envs" }}
- name: PG_DATABASE
  value: {{ tpl (pluck .Values.werf.env .Values.infra.pg.database | first | default .Values.infra.pg.database._default ) . | quote }}
- name: PG_HOST
  value: {{ tpl (pluck .Values.werf.env .Values.infra.pg.host | first | default .Values.infra.pg.host._default ) . | quote }}
- name: PG_USER
  value: {{ tpl (pluck .Values.werf.env .Values.infra.pg.user | first | default .Values.infra.pg.user._default ) . | quote }}
- name: PG_PASSWORD
  value: {{ pluck .Values.werf.env .Values.infra.pg.password | first | default .Values.infra.pg.password._default | quote }}
- name: PG_PORT
  value: {{ pluck .Values.werf.env .Values.infra.pg.port | first | default .Values.infra.pg.port._default | quote }}
- name: PG_SSLMODE
  value: {{ pluck .Values.werf.env .Values.infra.pg.sslmode | first | default .Values.infra.pg.sslmode._default | quote }}
- name: PG_SEARCH_PATH
  value: {{ pluck .Values.werf.env .Values.infra.pg.searchpath | first | default .Values.infra.pg.searchpath._default | quote }}
- name: cookie_name
  value: {{ .Chart.Name }}-{{ .Values.werf.env }}
- name: app_name
  value: {{ .Chart.Name }}
- name: tls
  value: {{ pluck .Values.werf.env .Values.app.tls | first | default .Values.app.tls._default | squote }}
- name: secure
  value: {{ pluck .Values.werf.env .Values.app.secure | first | default .Values.app.secure._default | squote }}
- name: selfreg
  value: {{ pluck .Values.werf.env .Values.app.selfreg | first | default .Values.app.selfreg._default | squote }}
- name: use_captcha
  value: {{ pluck .Values.werf.env .Values.app.use_captcha | first | default .Values.app.use_captcha._default | squote }}
- name: use_pin
  value: {{ pluck .Values.werf.env .Values.app.use_pin | first | default .Values.app.use_pin._default | squote }}
- name: max_attempts
  value: {{ pluck .Values.werf.env .Values.app.max_attempts | first | default .Values.app.max_attempts._default | squote }}
- name: reset_time
  value: {{ pluck .Values.werf.env .Values.app.reset_time | first | default .Values.app.reset_time._default | squote }}
- name: confirm_email_url
  value: {{ tpl (pluck .Values.werf.env .Values.app.confirm_email_url | first | default .Values.app.confirm_email_url._default) . | squote }}
- name: entry_point
  value: {{ tpl (pluck .Values.werf.env .Values.app.entry_point | first | default .Values.app.entry_point._default) . | squote }}
- name: smtp_address
  value: {{ pluck .Values.werf.env .Values.app.smtp_address | first | default .Values.app.smtp_address._default | squote }}
- name: from
  value: {{ pluck .Values.werf.env .Values.app.mail_from | first | default .Values.app.mail_from._default | squote }}
- name: admin_url
  value: {{  tpl (pluck .Values.werf.env .Values.app.admin_url | first | default .Values.app.admin_url._default) . | squote }}
- name: mail_tmpl_path
  value: {{  tpl (pluck .Values.werf.env .Values.app.mail_tmpl_path | first | default .Values.app.mail_tmpl_path._default) . | squote }}
- name: graphql_test_url
  value: {{  tpl (pluck .Values.werf.env .Values.app.graphql_test_url | first | default .Values.app.graphql_test_url._default) . | squote }}
- name: GIN_MODE
  value: "release"
{{- end }}
