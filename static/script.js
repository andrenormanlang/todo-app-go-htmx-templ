function toggleEditForm(taskID) {
    const editForm = document.getElementById(`edit-task-form-${taskID}`);
    const taskText = document.getElementById(`task-text-${taskID}`);
    if (editForm.style.display === 'none') {
        editForm.style.display = 'block';
        taskText.style.display = 'none';
    } else {
        editForm.style.display = 'none';
        taskText.style.display = 'block';
    }
}
