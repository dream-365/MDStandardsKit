#!/bin/sh
# Behavior tests for bin/doc-writer-tool. No network required.
set -eu

root=$(cd -P "$(dirname "$0")/.." && pwd)
tool="$root/bin/doc-writer-tool"
work=$(mktemp -d "${TMPDIR:-/tmp}/dw-shim-test.XXXXXX")
trap 'rm -rf "$work"' EXIT HUP INT TERM

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

release="$work/release"
mkdir -p "$release"
printf '#!/bin/sh\necho "FAKE_MDS $*"\n' >"$release/mdstandards-darwin-arm64"
printf '#!/bin/sh\necho "FAKE_QC $*"\n' >"$release/qualitycheck-darwin-arm64"
chmod +x "$release"/mdstandards-darwin-arm64 "$release"/qualitycheck-darwin-arm64

if command -v sha256sum >/dev/null 2>&1; then
  hash() { sha256sum "$1" | awk '{print $1}'; }
else
  hash() { shasum -a 256 "$1" | awk '{print $1}'; }
fi

write_checksums() {
  {
    printf '%s  mdstandards-darwin-arm64\n' "$(hash "$release/mdstandards-darwin-arm64")"
    printf '%s  qualitycheck-darwin-arm64\n' "$(hash "$release/qualitycheck-darwin-arm64")"
  } >"$release/checksums.txt"
}
write_checksums

cache="$work/cache"
mkdir -p "$cache"

# 1. download + execute
out=$(XDG_CACHE_HOME="$cache" MDSTANDARDS_RELEASE_BASE_URL="file://$release" "$tool" mdstandards --version) ||
  fail "download path exited non-zero"
printf '%s\n' "$out" | grep -q '^FAKE_MDS' || fail "download path did not execute the downloaded asset"

# 2. cache hit does not need the release
out=$(XDG_CACHE_HOME="$cache" MDSTANDARDS_RELEASE_BASE_URL="file://$work/missing" "$tool" mdstandards --version) ||
  fail "cache hit exited non-zero"
printf '%s\n' "$out" | grep -q '^FAKE_MDS' || fail "cache hit did not use the cached asset"

# 3. checksum mismatch fails closed and never executes
printf 'deadbeef  mdstandards-darwin-arm64\n' >"$release/checksums.txt"
rm -rf "$cache"
mkdir -p "$cache"
if XDG_CACHE_HOME="$cache" MDSTANDARDS_RELEASE_BASE_URL="file://$release" "$tool" mdstandards --version >"$work/mismatch.out" 2>"$work/mismatch.err"; then
  fail "checksum mismatch exited 0"
fi
grep -q 'checksum mismatch' "$work/mismatch.err" || fail "checksum mismatch did not report the cause"
grep -q 'FAKE_MDS' "$work/mismatch.out" && fail "checksum mismatch executed the asset"
write_checksums

# 4. build fallback with Go
rm -rf "$cache"
mkdir -p "$cache"
out=$(XDG_CACHE_HOME="$cache" MDSTANDARDS_SKIP_DOWNLOAD=1 "$tool" mdstandards --version) ||
  fail "build fallback exited non-zero"
printf '%s\n' "$out" | grep -q '^mdstandards ' || fail "build fallback did not execute the built tool"

# 5. no Go and no download fails with a clear message
safe="$work/safe"
mkdir -p "$safe"
for t in dirname readlink sed awk head mktemp mkdir rm mv chmod curl sha256sum shasum uname; do
  p=$(command -v "$t" 2>/dev/null) && ln -sf "$p" "$safe/$t"
done
ln -sf /bin/sh "$safe/sh"
rm -rf "$cache"
mkdir -p "$cache"
if PATH="$safe" XDG_CACHE_HOME="$cache" MDSTANDARDS_RELEASE_BASE_URL="file://$work/missing" "$tool" mdstandards --version >"$work/nogo.out" 2>"$work/nogo.err"; then
  fail "no-Go path exited 0"
fi
grep -q 'Go 1.27.1' "$work/nogo.err" || fail "no-Go path did not mention Go 1.27.1"

# 6. unsupported platform
fakebin="$work/fakebin"
mkdir -p "$fakebin"
printf '#!/bin/sh\ncase "$1" in -s) echo Windows;; -m) echo x86_64;; *) echo Windows;; esac\n' >"$fakebin/uname"
chmod +x "$fakebin/uname"
if PATH="$fakebin:$safe" XDG_CACHE_HOME="$cache" "$tool" mdstandards --version >"$work/platform.out" 2>"$work/platform.err"; then
  fail "unsupported platform exited 0"
fi
grep -q 'unsupported platform' "$work/platform.err" || fail "unsupported platform did not report the cause"

# 7. policy override requires the explicit allow flag
if MDSTANDARDS_POLICY_DIR=/tmp/policy XDG_CACHE_HOME="$cache" "$tool" mdstandards --version >"$work/policy.out" 2>"$work/policy.err"; then
  fail "policy override without allow flag exited 0"
fi
grep -q 'MDSTANDARDS_ALLOW_POLICY_OVERRIDE' "$work/policy.err" || fail "policy guard did not name the allow flag"

# 8. symlinked dispatcher resolves the plugin root
ln -sf "$tool" "$work/doc-writer-tool"
out=$(XDG_CACHE_HOME="$cache" MDSTANDARDS_SKIP_DOWNLOAD=1 "$work/doc-writer-tool" mdstandards --version) ||
  fail "symlinked dispatcher exited non-zero"
printf '%s\n' "$out" | grep -q '^mdstandards ' || fail "symlinked dispatcher did not run the tool"

echo "shim-test: all checks passed"
