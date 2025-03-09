FROM gradle:8.2.1-jdk17 AS builder
WORKDIR /home/gradle/project

COPY --chown=gradle:gradle . .

RUN ./gradlew shadowJar --no-daemon

FROM openjdk:17-jdk-slim
WORKDIR /app

COPY --from=builder /home/gradle/project/build/libs/app.jar app.jar

EXPOSE 8080

CMD ["java", "-jar", "app.jar"]
