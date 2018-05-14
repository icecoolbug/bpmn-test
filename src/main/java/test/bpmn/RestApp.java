package test.bpmn;

import org.camunda.bpm.engine.ProcessEngine;
import org.camunda.bpm.engine.RuntimeService;
import org.camunda.bpm.spring.boot.starter.annotation.EnableProcessApplication;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.beans.factory.annotation.Autowired;

@SpringBootApplication
@EnableProcessApplication
public class RestApp implements CommandLineRunner {

	@Autowired
	private RuntimeService runtimeService;
	
	@Autowired
	private ProcessEngine processEngine;

	public static void main(final String... args) throws Exception {
		SpringApplication.run(RestApp.class, args);
	}

	@Override
	public void run(String... strings) throws Exception {
		for ( ; ; ) {
			Thread.sleep(3000);
		}
	}
}