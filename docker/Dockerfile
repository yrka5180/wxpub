FROM python:2

WORKDIR /app
ADD n9e /app

COPY ./mybin ./
RUN chmod 0755 ./mybin
ENTRYPOINT [ "./mybin" ]