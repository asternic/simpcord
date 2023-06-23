# simpcord

A simple Discord to HTTP/Rest gateway using the awesome [discordgo](https://github.com/bwmarrin/discordgo) library

This process connects to Discord via websockets and calls the specified webhook whenever a message is received. It also accepts POST requests to send replies or messages back to Discord.

## Usage

Usage of ./simpcord:
```
  -a string
        Auth Token
  -l string
        Bind IP Address
  -p string
        Bind Port
  -t string
        Bot Token
  -w string
        Webhook URL
```

### Auth Token
A security token that is checked when receiving a POST request to `/send/{channelid}`. The token can be passed via an HTTP header named Token or through URL parameters.

### Bind IP Address
The IP address to bind to for accepting POST requests to send messages. Default value: `127.0.0.1`

### Bind Port
The TCP port to bind to. Default value: `8001`

### Bot Token
The Discord Bot Token obtained when creating a Bot for your Discord application.

### Webhook URL
The URL to send a webhook whenever a message is received. Message details will be sent in the `jsonData` parameter.

```json
{
  "author": {
    "id": "315436374476766171",
    "email": "",
    "username": "someusername",
    "avatar": "ad04f394f610c0f0422fbe3d3e014d5e",
    "locale": "",
    "discriminator": "0",
    "token": "",
    "verified": false,
    "mfa_enabled": false,
    "banner": "",
    "accent_color": 0,
    "bot": false,
    "public_flags": 0,
    "premium_type": 0,
    "system": false,
    "flags": 0
  },
  "avatar": "https://cdn.discordapp.com/avatars/715406378767876178/ad326f690f610c0f0452fbe3d1e01413e.png",
  "body": "Hello World",
  "channelid": "4448063866649532308",
  "event": "Message"
}
```

## Sending messages to Discord

To send a message as your bot to a specific channel, send a POST request to `/send/{channelid}` using simpcord. You can obtain the channelid in Discord by right-clicking on the channel name in your server and getting the invite link (the last long number in the URL).

Example:

```
curl -X POST http://localhost:8001/send/4448063866649532308 -H "Content-Type: application/json" -H "Token: 1234" --data-raw '{"body":"Hello world"}'
```

## License

Copyright &copy; 2023 Nicolás Gudiño

[MIT](https://choosealicense.com/licenses/mit/)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Legal

This code is in no way affiliated with, authorized, maintained, sponsored or
endorsed by Discord or any of its affiliates or subsidiaries. This is an
independent and unofficial software. Use at your own risk.
