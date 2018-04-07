package com.example.users.graceful;

import org.springframework.boot.actuate.health.HealthIndicator;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConditionalOnClass(HealthIndicator.class)
public class GracefulShutdownAutoConfiguration {
    @Bean
    HealthIndicator gracefulShutdownHealthCheck() {
        return new GracefulShutdownHealthCheck();
    }
}
