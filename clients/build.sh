docker build -t saashamor/siftclient .

docker push saashamor/siftclient 
ssh ec2-user@isift.info < deploy.sh