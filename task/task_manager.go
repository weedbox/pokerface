package task

type TaskManager struct {
	Tasks []interface{} `json:"tasks,omitempty"`
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make([]interface{}, 0),
	}
}

func (tm *TaskManager) AddTask(task Task) {
	if task == nil {
		return
	}

	tm.Tasks = append(tm.Tasks, task)
}

func (tm *TaskManager) GetTask(name string) Task {

	for _, t := range tm.Tasks {
		task := t.(Task)
		if task.GetName() == name {
			return task
		}
	}

	return nil
}

func (tm *TaskManager) TaskCount() int {
	return len(tm.Tasks)
}

func (tm *TaskManager) GetAvailableTask() Task {

	for _, t := range tm.Tasks {
		task := t.(Task)
		if !task.IsCompleted() {
			return task
		}
	}

	return nil
}

func (tm *TaskManager) Execute() {

	for {
		task := tm.GetAvailableTask()
		if task == nil {
			// No more task
			return
		}

		// Execute
		completed := task.Execute()
		if !completed {
			// Not yet completed
			return
		}
	}
}

func (tm *TaskManager) IsCompleted() bool {

	for _, t := range tm.Tasks {
		task := t.(Task)
		if !task.IsCompleted() {
			return false
		}
	}

	return true
}
