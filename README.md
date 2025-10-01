# Wroxen (Go) — Telegram search bot (TDLib + Bot API)

This repo runs a Telegram bot (Bot API) + a user client (TDLib) to search messages in a channel/group and return paginated results.

## Setup
1. Fill `.env` from `.env.example` or set env variables in your host / Render:
   - APP_ID, APP_HASH, BOT_TOKEN, SEARCH_CHAT_ID, TDLIB_DB_DIR (optional)

2. Build & run with docker:
   - `docker build -t wroxen-go .`
   - `docker run -v /path/to/tdlib-db:/data/tdlib-db --env-file .env wroxen-go`

3. Authorize TDLib user session:
   - Run the container and complete the phone->code flow inside the container (shell). Once authorized, user session stored in mounted TDLib DB.
   - Alternatively authorize locally and upload TDLib DB folder to your server volume.

4. Add bot to the group where users will type search queries.
5. Ensure the user account (the TDLib account) is a member of `SEARCH_CHAT_ID` (the channel/group you want to search).

## Notes
- Building TDLib inside docker is heavy — build takes several minutes.
- Persist `/data/tdlib-db` volume — it stores your user session.
