FROM nginx
COPY templates/index.html /usr/share/nginx/html/index.html
COPY assets/ /usr/share/nginx/html/assets/