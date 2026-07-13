#!/bin/sh
# Intended to be called by kindle-dash as: fetch-dashboard.sh /path/to/image.png
set -eu

output_path=${1:?usage: fetch-dashboard.sh OUTPUT_PATH}
: "${CLOCK_URL:?CLOCK_URL must point to the /clock PNG endpoint}"

output_dir=$(dirname "$output_path")
temporary_path="$output_dir/.clock-download-$$.png"
trap 'rm -f "$temporary_path"' EXIT HUP INT TERM

# kindle-dash ships xh, which is preferable to Kindle's legacy wget for HTTPS.
if [ -n "${DASHBOARD_HTTP_CLIENT:-}" ]; then
  "$DASHBOARD_HTTP_CLIENT" -d -q -o "$temporary_path" get "$CLOCK_URL"
else
  echo "DASHBOARD_HTTP_CLIENT is required (set it to kindle-dash's xh binary)." >&2
  exit 2
fi

test -s "$temporary_path"
mv -f "$temporary_path" "$output_path"
trap - EXIT HUP INT TERM
