function toggleEditForm(taskID) {
    const editForm = document.getElementById(`edit-task-form-${taskID}`);
    const taskText = document.getElementById(`task-text-${taskID}`);
    const buttons = document.getElementById(`task-buttons-${taskID}`);
    
    if (editForm.style.display === 'none') {
        editForm.style.display = 'block';
        taskText.style.display = 'none';
        buttons.style.display = 'none';  // Hide the Edit and Delete buttons
    } else {
        editForm.style.display = 'none';
        taskText.style.display = 'block';
        buttons.style.display = 'flex';  // Show the Edit and Delete buttons
    }
}
