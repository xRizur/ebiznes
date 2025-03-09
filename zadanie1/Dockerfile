FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    python3.11 \
    openjdk-8-jdk \
    kotlin \
    wget \
    unzip \
    curl

ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64
ENV PATH=$JAVA_HOME/bin:$PATH

ENV GRADLE_VERSION=8.13
RUN wget https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip && \
    unzip gradle-${GRADLE_VERSION}-bin.zip -d /opt/gradle && \
    rm gradle-${GRADLE_VERSION}-bin.zip
ENV GRADLE_HOME=/opt/gradle/gradle-${GRADLE_VERSION}
ENV PATH=$GRADLE_HOME/bin:$PATH

WORKDIR /app
COPY . /app

RUN gradle build

CMD ["java", "-cp", "build/libs/hello-world-1.0.0.jar", "com.example.HelloWorld"]

