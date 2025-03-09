FROM hseeberger/scala-sbt:11.0.11_1.5.5_2.13.6

WORKDIR /app
COPY . .

RUN sbt stage

RUN apt-get update && apt-get install -y wget unzip && \
    wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip && \
    unzip ngrok-stable-linux-amd64.zip && \
    mv ngrok /usr/local/bin && \
    rm ngrok-stable-linux-amd64.zip

EXPOSE 9000
CMD ["target/universal/stage/bin/play-scala-app"]
