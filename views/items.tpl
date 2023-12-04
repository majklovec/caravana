{{- range $v := .templates }}
{{- $val := $v | split "|" -}}
<option value="{{- index $val "_0" -}}">{{ index $val "_1" }} | {{ index $val "_2" }}</option>
{{- end }}