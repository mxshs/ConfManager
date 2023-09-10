FROM ubuntu

WORKDIR /app

RUN apt update
RUN apt install -y curl
RUN apt install -y git

RUN mkdir /root/.ssh/
RUN echo "${SSH_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

COPY . .

CMD ["./confmanager", "confmanager", "nvim_conf"]
