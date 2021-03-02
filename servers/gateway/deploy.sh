docker rm -f siftgateway

docker pull saashamor/siftgateway

export TLSCERT=/etc/letsencrypt/live/api.info441saasha.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.info441saasha.me/privkey.pem

docker run -d \
    --name siftgateway \
    -p 443:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    saashamor/siftgateway

exit