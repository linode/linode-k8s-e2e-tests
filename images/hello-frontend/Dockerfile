FROM debian:bullseye-20230202
RUN apt-get update -y && apt-get install curl -y && apt-get install nginx gettext-base -y

WORKDIR /etc/nginx/
COPY frontend.conf ./conf.d/frontend.conf
COPY nginx.conf ./nginx.conf
RUN rm -f /etc/nginx/conf.d/default.conf
EXPOSE 443/tcp 80/tcp
CMD ["nginx", "-g", "daemon off;"]
