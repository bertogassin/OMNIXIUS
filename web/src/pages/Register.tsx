import { Link } from 'react-router-dom';
import '../pages/Login.css';

/** Registration stub — temporarily disabled. Full form kept in history; re-enable when needed. */
export default function Register() {
  return (
    <div className="page login-page">
      <div className="login-card">
        <h1>Create account</h1>
        <p className="login-sub">Registration is temporarily disabled.</p>
        <p className="login-links"><Link to="/login">← Sign in</Link></p>
      </div>
    </div>
  );
}
