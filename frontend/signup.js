document.getElementById('signupForm').onsubmit = async function(e) {
  e.preventDefault();
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const firstName = document.getElementById('firstName').value;
  const lastName = document.getElementById('lastName').value;
  const res = await fetch('/api/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, first_name: firstName, last_name: lastName, user_type: 'USER' })
  });
  const data = await res.json();
  document.getElementById('result').textContent = data.message ? 'Signup successful!' : (data.error || 'Signup failed');
};
