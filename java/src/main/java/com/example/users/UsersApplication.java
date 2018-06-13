package com.example.users;

import com.example.users.graceful.GracefulshutdownSpringApplication;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;

@SpringBootApplication
@EnableJpaAuditing
public class UsersApplication {

	public static void main(String[] args) {
		if (args != null && args.length > 0 && args[0].equals("graceful")) {
			System.out.println("Starting graceful Users Application");
			GracefulshutdownSpringApplication.run(UsersApplication.class, args);
		} else {
			System.out.println("Starting Users Application");
			SpringApplication.run(UsersApplication.class, args);
		}
	}
}
