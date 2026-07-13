# Kindle client integration

Copy `fetch-dashboard.sh` to the `local/` directory of an existing
[kindle-dash](https://github.com/pascalw/kindle-dash) installation. Set
`CLOCK_URL` and `DASHBOARD_HTTP_CLIENT` in its environment. The script downloads
to a sibling temporary file and only replaces the current image after a successful,
non-empty response.

Use kindle-dash's bundled `xh` binary for HTTPS. Do not rely on the old Kindle
`wget` binary to support current TLS certificates. This repository intentionally
does not redistribute kindle-dash or any of its binaries.
