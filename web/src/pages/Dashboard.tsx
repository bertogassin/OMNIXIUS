import { Link } from 'react-router-dom';
import { useAuth } from '../auth';
import '../pages/Dashboard.css';

export default function Dashboard() {
  const { user } = useAuth();

  return (
    <div className="dashboard-page">
      <h1>Dashboard</h1>
      <p className="dashboard-welcome">
        Welcome, {user?.name || user?.email || 'User'}.
      </p>
      <div className="dashboard-actions">
        <Link to="/marketplace" className="dashboard-card">
          <h2>Marketplace</h2>
          <p>Browse products and services (Trove).</p>
        </Link>
        <a href="/app" className="dashboard-card dashboard-card-out">
          Full app (static) →
        </a>
      </div>
    </div>
  );
}
