FROM scratch
COPY build/linux/mssql/* /
ENV PATH /:$PATH
CMD ["broker-server", "--host","0.0.0.0","--port","61000"]
EXPOSE 61000
