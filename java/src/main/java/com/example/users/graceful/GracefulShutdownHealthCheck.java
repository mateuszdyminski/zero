package com.example.users.graceful;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.boot.actuate.health.Health;
import org.springframework.boot.actuate.health.HealthIndicator;


public class GracefulShutdownHealthCheck implements HealthIndicator, IProbeController {
    private static final Log log = LogFactory.getLog(GracefulShutdownHealthCheck.class);

    public static final String GRACEFULSHUTDOWN = "Gracefulshutdown";
    private Health health;

    GracefulShutdownHealthCheck() {
        setReady(true);
    }

    public Health health() {
        return health;
    }

    public void setReady(boolean ready) {
        if (ready) {
            health = new Health.Builder().withDetail(GRACEFULSHUTDOWN, "application up").up().build();
            log.info("Gracefulshutdown healthcheck up");
        } else {
            health = new Health.Builder().withDetail(GRACEFULSHUTDOWN, "gracefully shutting down").down().build();
            log.info("Gracefulshutdown healthcheck down");
        }
    }
}


