# lineTest



use post 127.0.0.1:3000/sendMessage to send message

"resText" can send at most 5 message one time
ex:
{
    "resText":["hello!","how are you?"],
    "sendTo":"{userid}"
}

use get 127.0.0.1:3000/getMessages?id={userid} to get user messages




