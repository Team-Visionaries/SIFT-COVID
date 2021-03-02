# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM nginx
ADD index.html /etc/nginx/html/index.html
ADD scripts.js /etc/nginx/html/scripts.js
ADD default.conf /etc/nginx/conf.d/default.conf
COPY . /usr/share/nginx/html
EXPOSE 443
