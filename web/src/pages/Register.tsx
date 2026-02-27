import { useState, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { api, ApiError } from '../api';
import '../pages/Login.css';

const MIN_PASSWORD = 8;
const MAX_PASSWORD = 128;
const MAX_NAME = 200;

const emailRegex = /^[a-zA-Z0-9.+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const passwordLetter = /[a-zA-Z]/;
const passwordNumber = /[0-9]/;
const passwordSymbol = /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?~\s]/;

function validateEmail(v: string): boolean {
  return v.length <= 255 && emailRegex.test(v.trim());
}

function passwordChecks(pwd: string) {
  return {
    length: pwd.length >= MIN_PASSWORD && pwd.length <= MAX_PASSWORD,
    letter: passwordLetter.test(pwd),
    number: passwordNumber.test(pwd),
    symbol: passwordSymbol.test(pwd),
  };
}

function allPasswordOk(checks: ReturnType<typeof passwordChecks>) {
  return checks.length && (checks.letter || checks.number || checks.symbol);
}

export default function Register() {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [name, setName] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const emailValid = useMemo(() => (email.trim() ? validateEmail(email) : null), [email]);
  const pwdChecks = useMemo(() => passwordChecks(password), [password]);
  const pwdOk = useMemo(() => allPasswordOk(pwdChecks), [pwdChecks]);
  const confirmOk = useMemo(() => !confirm || password === confirm, [password, confirm]);
  const nameOk = useMemo(() => name.length <= MAX_NAME, [name]);

  const canSubmit =
    email.trim() &&
    validateEmail(email) &&
    pwdOk &&
    password === confirm &&
    nameOk &&
    !submitting;

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!canSubmit) return;
    setErr(null);
    setSubmitting(true);
    api.auth
      .register(email.trim(), password, name.trim() || '')
      .then((r: { token?: string; user?: unknown }) => {
        if (r.token) api.setToken(r.token);
        if (r.user) api.user = r.user as { id: number; email?: string; name?: string; role?: string };
        navigate('/');
      })
      .catch((e: ApiError) => setErr(e.data?.error || 'Registration failed'))
      .finally(() => setSubmitting(false));
  };

  return (
    <div className="page login-page">
      <div className="login-card register-card">
        <h1>Create account</h1>
        <p className="login-sub">Create your OMNIXIUS account. All fields validated.</p>
        <form onSubmit={submit} className="login-form">
          {err && <p className="login-error">{err}</p>}

          <label>
            Email
            <input
              type="email"
              placeholder="name@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className={`login-input ${emailValid === false ? 'login-input-invalid' : ''} ${emailValid === true ? 'login-input-valid' : ''}`}
              autoComplete="email"
              maxLength={255}
            />
            {emailValid === false && <span className="register-hint">Valid email: letters, numbers, @ and a dot.</span>}
            {emailValid === true && <span className="register-hint register-hint-ok">Valid email</span>}
          </label>

          <label>
            Password
            <div className="login-pwd-wrap">
              <input
                type={showPassword ? 'text' : 'password'}
                placeholder={`${MIN_PASSWORD}–${MAX_PASSWORD} characters`}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={MIN_PASSWORD}
                maxLength={MAX_PASSWORD}
                className={`login-input ${password ? (pwdOk ? 'login-input-valid' : 'login-input-invalid') : ''}`}
                autoComplete="new-password"
              />
              <button type="button" className="login-pwd-toggle" onClick={() => setShowPassword((s) => !s)} aria-label={showPassword ? 'Hide password' : 'Show password'}>
                {showPassword ? 'Hide' : 'Show'}
              </button>
            </div>
            <ul className="register-requirements">
              <li className={pwdChecks.length ? 'register-ok' : ''}>{MIN_PASSWORD}–{MAX_PASSWORD} characters</li>
              <li className={pwdChecks.letter ? 'register-ok' : ''}>Letters (A–Z, a–z)</li>
              <li className={pwdChecks.number ? 'register-ok' : ''}>Numbers (0–9)</li>
              <li className={pwdChecks.symbol ? 'register-ok' : ''}>Symbols (!@#$% etc.)</li>
            </ul>
          </label>

          <label>
            Confirm password
            <input
              type={showPassword ? 'text' : 'password'}
              placeholder="Repeat password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
              className={`login-input ${confirm ? (confirmOk ? 'login-input-valid' : 'login-input-invalid') : ''}`}
              autoComplete="new-password"
            />
            {confirm && !confirmOk && <span className="register-hint">Passwords must match.</span>}
            {confirm && confirmOk && <span className="register-hint register-hint-ok">Passwords match</span>}
          </label>

          <label>
            Name (optional)
            <input
              type="text"
              placeholder="Your name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className={`login-input ${name.length > 0 && !nameOk ? 'login-input-invalid' : ''}`}
              autoComplete="name"
              maxLength={MAX_NAME}
            />
            {name.length > 0 && <span className="register-hint">{name.length}/{MAX_NAME} characters</span>}
          </label>

          <button type="submit" className="login-btn" disabled={!canSubmit}>
            {submitting ? 'Creating account…' : 'Create account'}
          </button>
        </form>
        <p className="login-links">
          <Link to="/login">Already have an account? Sign in</Link>
        </p>
      </div>
    </div>
  );
}
