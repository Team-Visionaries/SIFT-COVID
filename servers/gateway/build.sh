GOOS=linux go build
docker build -t saashamor/siftgateway .
go clean

docker push saashamor/siftgateway 
ssh ec2-user@ec2-54-70-88-211.us-west-2.compute.amazonaws.com < deploy.sh