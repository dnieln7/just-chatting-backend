# Just Chatting Backend

Back end for Just Chatting mobile apps.

# Deploy with commands

Create image

```shell
docker build . -t dnieln7/just-chatting
```

Create and run container

```shell
docker run -d -p 4200:4444 --restart always --name just-chatting dnieln7/just-chatting
```

# Deploy with script

Give executable permissions

```shell
chmod +x update-docker.sh 
```

Run script

```shell
./update-docker.sh
```
