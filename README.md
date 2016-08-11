# Telegram bot for [animeshot](https://as.bitinn.net)

## Live

You can directly send message to @animeshotbot in telegram to search [animeshot](https://as.bitinn.net)

or use inline mode to send photo from [animeshot](https://as.bitinn.net) to chat

## Run bot yourself

use go get to install:
    
    go get github.com/merrickluo/animeshotbot
    
run:

    usage: animeshotbot [<flags>]

    Flags:
          --help          Show context-sensitive help (also try --help-long and --help-man).
      -m, --mode="fetch"  start bot with fetch methods
      -p, --port=8185     listen port, use only with webhook mode
      -d, --debug         enable debug mode
          --version       Show application version.

Note that you must set ANIMESHOTBOT\_TG\_TOKEN (your bot token) environment variable before run this bot

## TODO

- [ ] full mode show result with offset and limit
- [ ] switch mode use command
- [ ] inline mode offset
