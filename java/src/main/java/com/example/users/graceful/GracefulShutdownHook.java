package com.example.users.graceful;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.util.StringUtils;

import java.util.Map;

class GracefulShutdownHook implements Runnable {
    protected static final String GRACEFUL_SHUTDOWN_WAIT_SECONDS = "estaGracefulShutdownWaitSeconds";
    private static final String DEFAULT_GRACEFUL_SHUTDOWN_WAIT_SECONDS = "5";

    private static final Log log = LogFactory.getLog(GracefulShutdownHook.class);

    private final ConfigurableApplicationContext applicationContext;

    GracefulShutdownHook(ConfigurableApplicationContext applicationContext) {
        this.applicationContext = applicationContext;
    }

    public void run() {
        setReadynessToFalse();
        delayShutdownSpringContext();
        shutdownSpringContext();
    }

    private void shutdownSpringContext() {
        log.info("Spring Application context starting to shutdown");
        applicationContext.close();
        log.info("Spring Application context is shutdown");
    }

    private void setReadynessToFalse() {
        log.info("Setting readyness for application to false, so the application doesn't receive new connections from Openshift");
        final Map<String, IProbeController> probeControllers = applicationContext.getBeansOfType(IProbeController.class);
        if (probeControllers.size() < 1) {
            log.error("Could not find a ProbeController Bean. Your ProbeController needs to implement the Interface: " + IProbeController.class.getName());
        }
        if (probeControllers.size() > 1) {
            log.warn("You have more than one ProbeController for Readyness-Check registered. " +
                    "Most probably one as Rest service and one in automatically configured as Actuator health check.");
        }
        for (IProbeController probeController : probeControllers.values()) {
            probeController.setReady(false);
        }
    }

    private void delayShutdownSpringContext() {
        try {
            int shutdownWaitSeconds = getShutdownWaitSeconds();
            log.info("Gonna wait for " + shutdownWaitSeconds + " seconds before shutdown SpringContext!");
            Thread.sleep(shutdownWaitSeconds * 1000);
        } catch (InterruptedException e) {
            log.error("Error while gracefulshutdown Thread.sleep", e);
        }
    }

    /**
     * Tries to get the value from the Systemproperty estaGracefulShutdownWaitSeconds,
     * otherwise it tries to read it from the application.yml, if there also not found 20 is returned
     *
     * @return The configured seconds, default 20
     */
    private int getShutdownWaitSeconds() {
        String waitSeconds = System.getProperty(GRACEFUL_SHUTDOWN_WAIT_SECONDS);
        if (StringUtils.isEmpty(waitSeconds)) {
            waitSeconds = applicationContext.getEnvironment().getProperty(GRACEFUL_SHUTDOWN_WAIT_SECONDS, DEFAULT_GRACEFUL_SHUTDOWN_WAIT_SECONDS);
        }
        return Integer.parseInt(waitSeconds);
    }
}