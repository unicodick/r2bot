<p align="center"><a href="https://github.com/unicodick/r2bot">
 <picture>
   <source srcset=".github/logo.png" />
   <img alt="r2bot" src="https://github.com/unicodick/r2bot" />
 </picture>
</a></p>

## deployment

1. copy `.env.example` to `.env` and set required values.
2. run with docker:

```bash
docker build -t r2bot .
docker run --env-file .env r2bot
```

## configuration

all configuration is done via environment variables:

- `BOT_TOKEN` - telegram bot token
- `ALLOWED_IDS` - telegram id
- `R2_ACCOUNT_ID` - cloudflare r2 account id
- `R2_ACCESS_KEY` - r2 access key
- `R2_SECRET_KEY` - r2 secret key
- `R2_BUCKET` - bucket name
- `R2_PUBLIC_URL` - public url bucket

## usage

- `/start` - show bot info
- send any file as a document - it will be uploaded to R2
- send any dircet URL - file will be downloaded and uploaded to R2
- receive a public link with a send button
- custom file names:
  - **with caption**: the caption will be used as the filename
  - if caption has no extension, the original file extension will be added
  - example: caption `examplenameforfile` + original file `document.pdf` → saved as `examplenameforfile.pdf`
  - **without caption**: falls back to the original filename
