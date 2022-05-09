FROM harbor.nova.net.cn/nova/alpine
COPY ./mybin ./
RUN chmod 0755 ./mybin
ENTRYPOINT [ "./mybin" ]