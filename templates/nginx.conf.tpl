#user  nobody;
worker_processes  1;

error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

# a file that will store the process ID of the main process
pid        logs/nginx.pid;


events {
    worker_connections  1024;
}

http {
  include mime.types;
  default_type  application/octet-stream;

  # define a log format named main and reference it at the access_log directive
  log_format  main '$remote_addr - $remote_user [$time_local]  $status '
    '"$request" $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"';
  access_log  logs/access.log  main;
  # Enabling the sendfile directive eliminates the step of copying the data into the buffer
  # and enables direct copying data from one file descriptor to another
  sendfile  on;
  # limit the amount of data transferred in a single sendfile() call
  sendfile_max_chunk 1m;
  # force the package to wait until it gets its maximum size (MSS) before sending it to the client.
  # This directive only works, when sendfile is on
  tcp_nopush  on;

  gzip_types text/plain text/css application/javascript application/json application/x-javascript text/xml application/xml application/xml+rss text/javascript;

  server {
      listen  80;
      root /usr/share/nginx/html;
      server_name mapi.ekwing.com;

      #charset koi8-r;
      #access_log  /var/log/nginx/host.access.log  main;

      location /teacher/class {
          if ($http_accept !~ text/html) {
            return 463;
          }
          try_files /class_list.html =404;
      }

      location ~ /code/mapi/teacher/[\d.]+/resource/class/(.*)$ {
          root /usr/share/nginx/assets;
          try_files /$1 =404;
      }

      location @upstream {
        proxy_pass http://{{.Upstream}}$request_uri;
        # 客户端的IP，进行反向代理的服务器对于目标代理服务器也是一个客户端
        proxy_set_header X-Real-IP $remote_addr;
        # $http_host：客户端请求头部的Host值，而$host表示的是客户端请求地址的hostname || 客户端请求头部的Host值 || nginx匹配到的server的name
        proxy_set_header Host $http_host;
        # 当只经过一次代理时，$proxy_add_x_forwarded_for的值与$remote_addr相等，每经过一次代理该变量都会在末尾拼接上" $remote_addr"
        proxy_set_header X-Forwarded-for $proxy_add_x_forwarded_for;
        proxy_set_header X-Nginx-Proxy true;
        proxy_set_header   Connection "";
        proxy_http_version 1.1;
        add_header Access-Control-Allow-Origin *;
      }

      error_page 463 =200 @upstream;

      #error_page  404              /404.html;

      # redirect server error pages to the static page /50x.html
      #
      error_page 500 502 503 504  /50x.html;
  }
}
