package com.example.users.graceful;

import org.springframework.boot.SpringApplication;
import org.springframework.context.ConfigurableApplicationContext;

public class GracefulshutdownSpringApplication {
    public static void run(Class<?> appClazz, String... args) {
        SpringApplication app = new SpringApplication(appClazz);
        app.setRegisterShutdownHook(false);
        ConfigurableApplicationContext applicationContext = app.run(args);
        Runtime.getRuntime().addShutdownHook(new Thread(new GracefulShutdownHook(applicationContext)));
    }
}
