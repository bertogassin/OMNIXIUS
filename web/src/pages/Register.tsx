import { Link } from 'react-router-dom';

/** Registration temporarily disabled. Re-enable later by restoring the form and api.auth.register. */
export default function Register() {
  return (
    <div className="page login-page">
      <div className="login-card">
        <h1>Create account</h1>
        <p className="page-intro">Registration is temporarily disabled.</p>
        <p><Link to="/login" className="btn btn-primary">← Sign in</Link></p>
      </div>
    </div>
  );
}
