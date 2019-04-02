set +x
docker images | grep "2018">docker.list
cat docker.list
while read line
do
    echo "============================="
    docker_name=`echo $line|awk -F " " '{print $1}'`
    docker_version=`echo $line|awk -F " " '{print $2}'`
    echo $docker_name
    echo $docker_version
    docker rmi -f ${docker_name}:${docker_version}
done<$WORKSPACE/docker.list