# shellcheck shell=sh
# Sourced by bin/doc-writer-tool. Expects _root and bin to be set.

_tool=${1:-}
if [ -z "$_tool" ]; then
  echo "doc-writer toolchain: missing tool name" >&2
  exit 2
fi
shift

case $_tool in
  mdstandards|qualitycheck|taskscheck) ;;
  *)
    echo "doc-writer toolchain: unknown tool '$_tool'" >&2
    exit 2
    ;;
esac

if [ -n "${MDSTANDARDS_POLICY_DIR:-}" ]; then
  if [ "${MDSTANDARDS_ALLOW_POLICY_OVERRIDE:-}" != "1" ]; then
    echo "doc-writer toolchain: MDSTANDARDS_POLICY_DIR requires MDSTANDARDS_ALLOW_POLICY_OVERRIDE=1" >&2
    exit 2
  fi
  echo "doc-writer toolchain: warning: using policy override $MDSTANDARDS_POLICY_DIR" >&2
fi

_manifest="$_root/.claude-plugin/plugin.json"
if [ ! -f "$_manifest" ]; then
  echo "doc-writer toolchain: missing $_manifest" >&2
  exit 2
fi
_version=$(sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$_manifest" | head -1)
if [ -z "$_version" ]; then
  echo "doc-writer toolchain: cannot read version from $_manifest" >&2
  exit 2
fi

_os=$(uname -s)
_arch=$(uname -m)
case "$_os:$_arch" in
  Darwin:arm64) _platform=darwin-arm64 ;;
  Darwin:x86_64) _platform=darwin-amd64 ;;
  Linux:x86_64) _platform=linux-amd64 ;;
  Linux:aarch64) _platform=linux-arm64 ;;
  *)
    echo "doc-writer toolchain: unsupported platform $_os/$_arch; supported: darwin-arm64, darwin-amd64, linux-amd64, linux-arm64. Build manually with Go 1.27.1+ or use a supported platform." >&2
    exit 2
    ;;
esac

_cache="${XDG_CACHE_HOME:-$HOME/.cache}/mdstandards/$_version/$_platform"
_bin="$_cache/$_tool"

if [ -n "${MDSTANDARDS_BIN:-}" ]; then
  if [ -x "${MDSTANDARDS_BIN}/$_tool" ]; then
    exec "${MDSTANDARDS_BIN}/$_tool" "$@"
  fi
  echo "doc-writer toolchain: MDSTANDARDS_BIN is set but ${MDSTANDARDS_BIN}/$_tool is missing or not executable" >&2
  exit 2
fi

if [ -x "$_bin" ]; then
  exec "$_bin" "$@"
fi

if ! mkdir -p "$_cache" 2>/dev/null; then
  echo "doc-writer toolchain: cannot create cache dir $_cache" >&2
  exit 2
fi

if [ "${MDSTANDARDS_SKIP_DOWNLOAD:-}" != "1" ]; then
  _base="${MDSTANDARDS_RELEASE_BASE_URL:-https://github.com/dream-365/MDStandardsKit/releases/download/v$_version}"
  if command -v curl >/dev/null 2>&1; then
    _fetch() { curl -fLsS "$1"; }
  elif command -v wget >/dev/null 2>&1; then
    _fetch() { wget -qO- "$1"; }
  else
    _fetch() { return 1; }
  fi

  if _checksums=$(_fetch "$_base/checksums.txt" 2>/dev/null); then
    _asset="$_tool-$_platform"
    _sha=$(printf '%s\n' "$_checksums" | awk -v a="$_asset" '$2 == a { print $1; exit }')
    if [ -n "$_sha" ]; then
      _tmp=$(mktemp "$_cache/.download.XXXXXX") || {
        echo "doc-writer toolchain: mktemp failed in $_cache" >&2
        exit 2
      }
      if _fetch "$_base/$_asset" >"$_tmp" 2>/dev/null; then
        if command -v sha256sum >/dev/null 2>&1; then
          _actual=$(sha256sum "$_tmp" | awk '{print $1}')
        elif command -v shasum >/dev/null 2>&1; then
          _actual=$(shasum -a 256 "$_tmp" | awk '{print $1}')
        else
          rm -f "$_tmp"
          echo "doc-writer toolchain: no sha256 tool (sha256sum or shasum)" >&2
          exit 2
        fi
        if [ "$_actual" != "$_sha" ]; then
          rm -f "$_tmp"
          echo "doc-writer toolchain: checksum mismatch for $_asset" >&2
          exit 2
        fi
        if ! chmod 0755 "$_tmp" 2>/dev/null; then
          rm -f "$_tmp"
          echo "doc-writer toolchain: chmod failed for $_tmp" >&2
          exit 2
        fi
        if ! mv -f "$_tmp" "$_bin" 2>/dev/null; then
          rm -f "$_tmp"
          echo "doc-writer toolchain: mv failed for $_bin" >&2
          exit 2
        fi
        exec "$_bin" "$@"
      fi
      rm -f "$_tmp"
    fi
  fi
fi

if ! command -v go >/dev/null 2>&1; then
  echo "doc-writer toolchain: download unavailable and Go 1.27.1+ is not installed" >&2
  exit 2
fi

case $_tool in
  mdstandards) _pkg="." ;;
  qualitycheck) _pkg="./tools/qualitycheck" ;;
  taskscheck) _pkg="./tools/taskscheck" ;;
esac

_tmp=$(mktemp "$_cache/.build.XXXXXX") || {
  echo "doc-writer toolchain: mktemp failed in $_cache" >&2
  exit 2
}
if ! (cd "$_root" && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$_tmp" "$_pkg"); then
  rm -f "$_tmp"
  echo "doc-writer toolchain: go build failed for $_tool" >&2
  exit 2
fi
chmod 0755 "$_tmp" 2>/dev/null || true
if ! mv -f "$_tmp" "$_bin" 2>/dev/null; then
  rm -f "$_tmp"
  echo "doc-writer toolchain: mv failed for $_bin" >&2
  exit 2
fi
exec "$_bin" "$@"
