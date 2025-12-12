document.getElementById('loginForm').onsubmit = async function(e) {
  e.preventDefault();
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await res.json();
  document.getElementById('result').textContent = data.access_token ? 'Login successful!' : (data.error || 'Login failed');
};
