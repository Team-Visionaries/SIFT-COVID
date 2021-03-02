docker rm -f siftclient

docker pull saashamor/siftclient

docker run -d \
    --name siftclient \
    -p 443:443 \
    -p 80:80 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    saashamor/siftclient

exit