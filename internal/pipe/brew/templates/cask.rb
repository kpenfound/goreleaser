# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
{{ if .CustomRequire -}}
require_relative "{{ .CustomRequire }}"
{{ end -}}
class {{ .Name }} < Formula
  desc "{{ .Desc }}"
  homepage "{{ .Homepage }}"
  version "{{ .Version }}"
  {{- if .License }}
  license "{{ .License }}"
  {{- end }}
  {{- with .Dependencies }}
  {{ range $index, $element := . }}
  depends_on "{{ .Name }}"
  {{- if .Type }} => :{{ .Type }}{{- else if .Version }} => "{{ .Version }}"{{- end }}
  {{- with .OS }} if OS.{{ . }}?{{- end }}
  {{- end }}
  {{- end -}}

  {{- if and (not .LinuxPackages) .MacOSPackages }}
  depends_on :macos
  {{- end }}
  {{- if and (not .MacOSPackages) .LinuxPackages }}
  depends_on :linux
  {{- end }}
  {{- printf "\n" }}

  {{- if and .MacOSPackages .LinuxPackages }}
  on_macos do
  {{- include "macos_packages" . | indent 2 }}
  end

  on_linux do
  {{- include "linux_packages" . | indent 2 }}
  end
  {{- end }}

  {{- if and (.MacOSPackages) (not .LinuxPackages) }}
  {{- template "macos_packages" . }}
  {{- end }}

  {{- if and (not .MacOSPackages) (.LinuxPackages) }}
  {{- template "linux_packages" . }}
  {{- end }}

  {{- with .Conflicts }}
  {{ range $index, $element := . }}
  conflicts_with "{{ . }}"
  {{- end }}
  {{- end }}

  {{- with .CustomBlock }}
  {{ range $index, $element := . }}
  {{ . }}
  {{- end }}
  {{- end }}

  {{- with .PostInstall }}

  def post_install
    {{- range . }}
    {{ . }}
    {{- end }}
  end
  {{- end -}}

  {{- with .Caveats }}

  def caveats
    <<~EOS
    {{- range $index, $element := . }}
      {{ . -}}
    {{- end }}
    EOS
  end
  {{- end -}}

  {{- with .Service }}

  service do
    {{- range . }}
    {{ . }}
    {{- end }}
  end
  {{- end -}}

  {{- if .Tests }}

  test do
    {{- range $index, $element := .Tests }}
    {{ . -}}
    {{- end }}
  end
  {{- end }}
end
