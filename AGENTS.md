# Kindle Clock agent notes

- The standard clock endpoint is intentionally independent of Nature Remo, Awair,
  SwitchBot, Redis, OpenWeatherMap, and their tokens.
- Kindle 4 output must remain exactly 600 x 800 pixels by default. Verify the
  decoded PNG dimensions, not only the HTTP status.
- Run formatting checks, `go vet ./...`, `go test ./...`, and `go build ./...`
  before handoff. Build and exercise the Docker image when Docker is available.
- Do not add secrets, tokens, or private device addresses to tracked files.
- `DISPLAY_ROTATION=90` is the default landscape mount configuration; both
  supported rotations (0 and 90) must retain the final configured dimensions.
