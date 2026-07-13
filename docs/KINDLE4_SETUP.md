# Kindle 4 clock setup

## Architecture

The server renders a grayscale PNG. A Kindle client (for example kindle-dash)
downloads `/clock`, displays it with `eips`, and schedules the next refresh.
The server has no dependency on room sensors or external tokens for this route.

## Run locally

```powershell
$env:APP_TIMEZONE = 'Europe/Warsaw'
$env:DISPLAY_ROTATION = '90'
go run .
Invoke-WebRequest http://127.0.0.1:8080/health
Invoke-WebRequest http://127.0.0.1:8080/clock -OutFile .\clock.png
```

`/clock` and `/` return an uncached `image/png`. Defaults are `DISPLAY_WIDTH=600`,
`DISPLAY_HEIGHT=800`, and `DISPLAY_ROTATION=90`. Rotation 90 renders a landscape
canvas and rotates it into a 600 x 800 PNG; rotation 0 renders portrait directly.

## Docker

```powershell
docker build -t kindle-clock-k4 .
docker run --rm -p 8080:8080 -e APP_TIMEZONE=Europe/Warsaw kindle-clock-k4
```

In another terminal, download `http://127.0.0.1:8080/clock` and inspect it before
connecting a Kindle. The image must be 600 x 800 and contain both black text and
a white background.

## Render

`render.yaml` is a Blueprint for a Docker web service with `/health` as its health
check. Create the Blueprint from the fork in the Render dashboard. Render provides
`PORT`; do not configure it manually. Set a unique service name if `kindle-clock-k4`
is already taken. The free plan can sleep when idle, so a first refresh can be slow.

## Kindle client

Use the files under `kindle-client/` with a separately installed kindle-dash.
Configure `CLOCK_URL` with the public `/clock` URL and set
`DASHBOARD_HTTP_CLIENT` to the bundled `xh` executable. The script preserves the
last known-good image when a download fails. Older Kindles can fail modern HTTPS
when using their built-in `wget`; use the dashboard client's HTTPS-capable binary.

## Rollback and troubleshooting

To roll back, deploy the previous Git commit in Render or point the Kindle client
at the previous `/clock` URL. Check `/health` first, then download `/clock` from a
computer on the same network. If the PNG is wrong-sized, confirm all three
`DISPLAY_*` variables. If time is wrong, set `APP_TIMEZONE=Europe/Warsaw` and
restart the service. Invalid timezone or display values deliberately prevent startup
or return an HTTP error instead of silently producing an unsuitable dashboard.

## Device modification warning

Jailbreaking and installing third-party software on a Kindle are neither official
nor risk-free. Confirm the exact Kindle model and firmware against current community
instructions before copying any modification package. This repository only supplies
the clock server and client integration example.
