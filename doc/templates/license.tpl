# Dependency licenses
{{ range . }}
 - {{.Name}} {{.Version}} ([{{.LicenseName}}]({{.LicenseURL}}))
{{- end }}

