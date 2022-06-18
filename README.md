# arduino-help-bot
Bot for help


```
docker run --rm -it -v $(pwd)/docs:/docs -e BOT_TOKEN="tokenhere" -e BOT_ADMIN_ROLES="admins,roles,ids,seperated,by,commad" -e "exe,ini,ino" ghcr.io/9glt/arduino-help-bot:latest
```


for local develpment:
```
./build-and-run.sh [bottokengoeshere] [admin,roles,ids,for,admin,commands] [exe,ini,txt,etc...]
```