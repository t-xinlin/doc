set -e
set +x
export DOCKER_LIST=docker.list
read -t 30 -p "input images name:" IMAGE_NAME
echo $IMAGE_NAME
if [ "$IMAGE_NAME" = "" ]
then
  echo "nothing input, exit it"
  exit 1
else
  echo ">>$IMAGE_NAME"
fi
docker images | grep "$IMAGE_NAME">$DOCKER_LIST
cat docker.list
while read line
do
    echo "============================="
    docker_name=`echo $line|awk -F " " '{print $1}'`
    docker_version=`echo $line|awk -F " " '{print $2}'`
    echo $docker_name
    echo $docker_version
    docker rmi -f ${docker_name}:${docker_version}
done<$DOCKER_LIST
rm -rf $DOCKER_LIST
echo "Done>>>>>>>>>>>>>>>>"
