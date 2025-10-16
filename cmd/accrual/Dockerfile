FROM ubuntu:22.04

COPY cmd/accrual/accrual_linux_amd64 /server
RUN chmod +x /server

EXPOSE 8081
CMD ["/server"]