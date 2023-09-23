// Define a function to get the value of a cookie by name
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

$(document).ready(function () {
    // Check for JWT token in the cookie when the page loads
    const token = getCookie('jwt-token');
    if (token) {
        // Circular Button Click Event
        $('#profileButton').click(function () {
            // Send an AJAX request to fetch user details
            $.ajax({
                url: 'http://localhost:8080/profile',
                type: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }, // Include the JWT token in the request headers
                success: function (response) {
                    // Handle the response and display user details as needed
                    alert('User Details: ' + JSON.stringify(response));
                },
                error: function () {
                    // Handle error by displaying an error message
                    alert('Failed to fetch user details. Please try again.');
                }
            });
        });
        // If token exists, fetch tasks and display them
        fetchTasks(token);
    } else {
        // If no token, redirect to login or handle as needed
        window.location.href = 'http://localhost:8080'; // Update the URL as needed
    }
});

// Add Task Form
$('#saveTaskBtn').click(function () {
    const title = $('#title').val();
    if (title.trim() === '') {
        $('#feedbackMessage').text('Task title cannot be empty.');
        return;
    }

    const token = getCookie('jwt-token'); // Get the token from the cookie

    // Send a POST request to create a new task with the token in the request headers
    $.ajax({
        url: 'http://localhost:8080/task/create',
        type: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        data: JSON.stringify({ title }),
        contentType: 'application/json',
        success: function () {
            fetchTasks(); // Refresh the task list
            $('#title').val(''); // Clear the input field
        },
        error: function () {
            $('#feedbackMessage').text('Failed to create a task.');
        }
    });
});

// Get All Tasks Button
$('#getTasksBtn').click(function () {
    fetchTasks();
});

// Function to fetch tasks and update the table
function fetchTasks() {
    const token = getCookie('jwt-token'); // Get the token from the cookie
    $.ajax({
        url: 'http://localhost:8080/tasks',
        type: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        success: function (response) {
            const data = response.data;
            const taskTable = $('#taskTable tbody');
            taskTable.empty();
            data.forEach(function (task, index) {
                taskTable.append(`
            <tr>
                <th scope="row">${index + 1}</th>
                <td>${task.title}</td>
                <td>${task.status}</td>
                <td>${task.priority}</td>
                <td>
                    <button class="btn btn-danger delete-task" data-id="${task.id}">Delete</button>
                    <button class="btn btn-primary update-task" data-id="${task.id}">Update Status</button>
                    <button class="btn btn-info update-priority" data-id="${task.id}">Update Priority</button>
                </td>
            </tr>
        `);
            });

            // Delete Task Button
            $('.delete-task').click(function () {
                const taskId = $(this).data('id');
                const token = getCookie('jwt-token'); // Get the token from the cookie
                $.ajax({
                    url: `http://localhost:8080/task/delete/${taskId}`,
                    type: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    },
                    success: function () {
                        fetchTasks();
                    },
                    error: function () {
                        $('#feedbackMessage').text('Failed to delete the task.');
                    }
                });
            });

            // Update Task Button (open modal or handle as needed)
            $('.update-task').click(function () {
                const taskId = $(this).data('id');
                const token = getCookie('jwt-token'); // Get the token from the cookie
                $('#updateStatusModal').modal('show');
                $('#updateStatusBtn').click(function () {
                    const newStatus = $('#statusSelect').val();
                    $.ajax({
                        url: `http://localhost:8080/task/update/${taskId}`,
                        type: 'PUT',
                        headers: {
                            'Authorization': `Bearer ${token}`
                        },
                        data: JSON.stringify({ status: newStatus }),
                        contentType: 'application/json',
                        success: function () {
                            fetchTasks();
                            $('#updateStatusModal').modal('hide');
                        },
                        error: function () {
                            $('#updateStatusError').text('Failed to update status. Please try again.');
                        }
                    });
                });
            });

            // Update Priority Button (open modal for priority update)
            $('.update-priority').click(function () {
                const taskId = $(this).data('id');
                const token = getCookie('jwt-token'); // Get the token from the cookie
                $('#updatePriorityModal').modal('show');
                $('#updatePriorityBtn').click(function () {
                    const newPriority = $('#prioritySelect').val();
                    $.ajax({
                        url: `http://localhost:8080/task/update/${taskId}`,
                        type: 'PUT',
                        headers: {
                            'Authorization': `Bearer ${token}`
                        },
                        data: JSON.stringify({ priority: newPriority }),
                        contentType: 'application/json',
                        success: function () {
                            fetchTasks();
                            $('#updatePriorityModal').modal('hide');
                        },
                        error: function () {
                            $('#updatePriorityError').text('Failed to update priority. Please try again.');
                        }
                    });
                });
            });
        },
        error: function () {
            $('#feedbackMessage').text('Failed to fetch tasks.');
        }
    });
}
