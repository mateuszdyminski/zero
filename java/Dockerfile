FROM openjdk:8-jre-alpine

RUN mkdir -p /usr/share/zero

ADD target/users-1.0.0.jar /usr/share/zero/users-1.0.0.jar

EXPOSE 8080

ENTRYPOINT [ "java", "-Djava.security.egd=file:/dev/./urandom", "-jar", "/usr/share/zero/users-1.0.0.jar" ]