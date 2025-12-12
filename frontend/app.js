document.getElementById('loginBtn').onclick = function() {
  window.location.href = './login.html';
  // If index.html is in /frontend/, this will go to /frontend/login.html
};
document.getElementById('signupBtn').onclick = function() {
  window.location.href = './signup.html';
  // If index.html is in /frontend/, this will go to /frontend/signup.html
};
document.getElementById('googleLoginBtn').onclick = function() {
    window.location.href = 'http://localhost:8000/api/google-login';
};
