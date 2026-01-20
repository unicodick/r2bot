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
- receive a public link with a send button
