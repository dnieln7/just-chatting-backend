#echo "Fetching latest changes from git..."
#git pull
echo "Trying to stop old container"
docker stop just-chatting
echo "Deleting old container"
docker container rm just-chatting --force
echo "Deleting old image"
docker rmi dnieln7/just-chatting
echo "Creating new image"
docker build . -t dnieln7/just-chatting
echo "Creating new container"
docker run -d -p 4200:4444 --restart always --name just-chatting dnieln7/just-chatting
echo "Deleting dangling images"
docker image prune --force
