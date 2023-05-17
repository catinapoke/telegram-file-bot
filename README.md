# Telegram file bot
It's in WIP state and going to be used for save data transfer, so you will be able to send crypted files by link and send password via different way.

## Development stages

1. Core  
    1.1 Write basic commands
    - ✔️start
    - ✔️load
    - get
    - delete
    - list

    1.2 Add async in go-microservice  
    1.3 Errors handling with errors.As/Is in Handler, so users will understand errors  
    1.4 Autocreating database schema at new machines  
    1.5 Admin features  
  
2. Encoding / decoding
  - setkey
  - erasekey
  - encoding on load, decoding on get
3. Share settings
  - /changeshare {name}
  - Sending file state after loading
4. [Microservice] Web telegram bot, so user will be able to skip telegram file storage and don't show raw file

## Run
Create secrets in folder ./tmp and run docker compose via 
`docker build . && docker compose up`