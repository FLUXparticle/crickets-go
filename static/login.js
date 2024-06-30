async function login() {
    const form = document.getElementById('loginForm');
    const formData = new FormData(form);
    const json = JSON.stringify(Object.fromEntries(formData.entries()));
    const errorMessage = document.getElementById('error-message');

    const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: json
    });

    if (response.ok) {
        // Login successful, redirect to the main app
        window.location.href = '/app/';
    } else {
        // Show error message
        const result = await response.json();
        errorMessage.textContent = result.error;
    }
}