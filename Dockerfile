FROM nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY assets/ /usr/share/nginx/html/assets/
COPY templates/index.html /usr/share/nginx/html/index.html
