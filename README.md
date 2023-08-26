# Go Whatsapp Bot 

#### Minimalist Whatsapp Bot API using the official [Whatsapp Cloud API](https://developers.facebook.com/docs/whatsapp/cloud-api/) from Meta (Facebook).


## Logic
#### On the user's first interaction, the bot will send a prompt asking the user to choose a preferred language. After the language has been chosen, the bot control will be set to `false` allowing an authorized admin to take control over the conversation. This approach allows you to have multiple admins being able to chat with clients through the same bot account.


## Extending
#### There are many ways the bot can be extended, you can automate by creating message, media, links templates and much more. I recommend looking into the [cloud appi guides docs](https://developers.facebook.com/docs/whatsapp/cloud-api/guides) and [references](https://developers.facebook.com/docs/whatsapp/cloud-api/reference).


https://github.com/4d3v/whatsapp-cloudapi-bot/assets/55705104/58ad6d91-006f-40d2-b4d4-634ac78cb130


## Setup And External Packages
#### For the setup such as Authentication, Database, etc. There are many options, but I decided to go with Firebase's because I never used it before and wanted to learn, and for its cool features such as realtime updates.
| <!-- --> | <!-- --> | <!-- --> |
| --- | --- | --- |
| fasthttp | [github.com/valyala/fasthttp](https://github.com/valyala/fasthttp) | Fast HTTP implementation for Go. |
| firebase | [firebase.google.com/go](https://firebase.google.com/go) | Firebase Admin Go SDK |
| jsoniter | [github.com/json-iterator/go](https://github.com/json-iterator/go) | Drop-in replacement of "encoding/json" |
| godotenv | [github.com/joho/godotenv](https://github.com/joho/godotenv) | Loads env vars from a .env file |


> **Note**
>
> I actually first built this bot on NodeJS and then rewrote to Golang. the Node's implementation will be available soon as well, which is pretty much using the same logic as this one.


## How To Run
Assuming you have a Facebook dev account, you will have to first log in here [https://developers.facebook.com](https://developers.facebook.com) and create an app. Then I would recommend you [start here](https://developers.facebook.com/docs/whatsapp/cloud-api/). 

The second part would be to create a [Firestore Database](https://firebase.google.com/docs/firestore), grab the serviceAccount.json config and put it on the project's root directory. The database is optional, you can of course replace it with anything you want. ``Take a look into /pkg/db /pkg/driver /pkg/models/userData.go``

To make it short, you need to have a working webserver with https and connect using the Meta's provided token and then configure the webhook to receive message notifications.
Also create a .env file and fill the variables `ACCESS_TOKEN= APP_SECRET= HUB_TOKEN= WHATSAPP_PHONE_ID= WHATSAPP_BUSINESS_ID=` with the Meta's provided values.