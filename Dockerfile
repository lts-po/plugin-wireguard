FROM ubuntu:21.04 as builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y nftables iproute2 netcat inetutils-ping net-tools nano ca-certificates git curl
RUN mkdir /code
WORKDIR /code
ARG TARGETARCH
RUN curl -O https://dl.google.com/go/go1.17.linux-${TARGETARCH}.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.linux-${TARGETARCH}.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
COPY code/ /code/

#RUN --mount=type=tmpfs,target=/root/go/ (go build -ldflags "-s -w" -o /wireguard_plugin /code/wireguard_plugin.go && ./wireguard_plugin)
RUN (go build -ldflags "-s -w" -o /wireguard_plugin /code/wireguard_plugin.go)


FROM ubuntu:21.04
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y nftables iproute2 netcat inetutils-ping net-tools nano ca-certificates curl
RUN apt-get install -y wireguard-tools
COPY scripts /scripts/
COPY --from=builder /wireguard_plugin /
ENTRYPOINT ["/scripts/startup.sh"]