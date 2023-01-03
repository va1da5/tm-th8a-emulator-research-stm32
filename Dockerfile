FROM tinygo/tinygo

USER root

RUN apt update \
    && apt install -y git ssh make libtool pkgconf autoconf automake \
    texinfo libusb-1.0-0 libusb-1.0-0-dev libjaylink0 \
    && apt -y clean \
    && apt -y autoremove \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /

RUN git clone https://github.com/STMicroelectronics/OpenOCD.git \
    && cd /OpenOCD \
    && ./bootstrap \
    && ./configure --enable-stlink \
    && make \
    && make install \
    && rm -rf /OpenOCD

RUN go install -v golang.org/x/tools/gopls@latest  \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && bash -c "ln -s $(tinygo info | grep GOROOT | cut -d':' -f2 | xargs) /root/.cache/tinygo/goroot"
# && go install -v github.com/ramya-rao-a/go-outline \

WORKDIR /
