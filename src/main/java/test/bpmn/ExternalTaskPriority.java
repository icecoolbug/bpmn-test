package test.bpmn;

import org.springframework.stereotype.Component;

@Component
public class ExternalTaskPriority {
	public long getTaskPriority() {
		return Long.MAX_VALUE-System.currentTimeMillis();
	}

}
