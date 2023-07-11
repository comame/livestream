FROM ubuntu

RUN apt update && apt install -y build-essential libpcre3 libpcre3-dev libssl-dev libz-dev curl git

WORKDIR /root

RUN curl -LO http://nginx.org/download/nginx-1.24.0.tar.gz \
    && tar xvzf nginx-1.24.0.tar.gz \
    && git clone https://github.com/arut/nginx-rtmp-module.git

WORKDIR /root/nginx-1.24.0

RUN ./configure --with-http_ssl_module --add-module=../nginx-rtmp-module \
    && make \
    && make install

RUN mkdir -p /var/local/www/hls && chmod 777 /var/local/www/hls

COPY ./nginx.conf /usr/local/nginx/conf/nginx.conf

EXPOSE 1935
EXPOSE 8080

CMD /usr/local/nginx/sbin/nginx; while : ; do sleep 1; done
