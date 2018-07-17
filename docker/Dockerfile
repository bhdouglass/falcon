FROM clickable/ubuntu-sdk:16.04-armhf

RUN wget https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz && \
    tar -xvf go1.5.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    mv go /usr/local && \
    rm go1.5.linux-amd64.tar.gz
