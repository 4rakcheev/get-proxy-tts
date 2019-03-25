# get-proxy-tts
Service provides HTTP GET interface for access to text to speech functionality like AWS Polly, Google TTS or others services.

## Overview
This can be useful for easy integration TTS to your legacy systems that integrating with audio files by URL.
Service caching generated speeches with ident parameters to local folder.

## Usage

```bash
go get github.com/4rakcheev/get-proxy-tts
go run get-proxy-tts.go
```

## Text-to-speech generating
Simple generating text with simple text type (without SSML) like:

http://localhost:80/polly?access_key=AKIAJC56K3IEGLBH72NA&secret_key=OQgTW1yigF%2BjPE1yDOa7Zh6%2BQ/TLk5QQufaK4d0I&voice=Brian&text=Test%20message%20for%20live%20demo%21&language=en-US

## Supported parameters
- `access_key` AWS IAM access key
- `secret_key` AWS IAM Secret key. Don't forget replace `+` symbol to `%2B` in this string. TDB: create live helper to generate correct URL
- `voice` [Polly Voice ID](https://docs.aws.amazon.com/en_us/polly/latest/dg/API_SynthesizeSpeech.html#polly-SynthesizeSpeech-request-LanguageCode) 
- `language` [Polly Language ID](https://docs.aws.amazon.com/en_us/polly/latest/dg/API_SynthesizeSpeech.html#polly-SynthesizeSpeech-request-LanguageCode)
- `text` String to speech