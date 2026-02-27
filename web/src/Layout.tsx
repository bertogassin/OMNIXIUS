import { Link, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from './auth';
import './Layout.css';

export default function Layout() {
  const { token, user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="layout">
      <header className="layout-header">
        <div className="layout-header-inner">
          <Link to="/" className="layout-logo">OMNIXIUS</Link>
          <nav className="layout-nav">
            {token ? (
              <>
                <Link to="/">Dashboard</Link>
                <Link to="/marketplace">Marketplace</Link>
                <Link to="/orders">Orders</Link>
                <Link to="/find-professional">Find professional</Link>
                <Link to="/mail">Mail</Link>
                <Link to="/notifications">Notifications</Link>
                <Link to="/profile">Profile</Link>
                <Link to="/wallet">Wallet</Link>
                <Link to="/trade">Trade</Link>
                <Link to="/ai">AI</Link>
                <Link to="/settings">Settings</Link>
                <Link to="/admin">Admin</Link>
                <span className="layout-user">{user?.email ?? user?.name ?? 'User'}</span>
                <button type="button" className="layout-btn" onClick={handleLogout}>
                  Sign out
                </button>
              </>
            ) : (
              <Link to="/login">Sign in</Link>
            )}
          </nav>
        </div>
      </header>
      <main className="layout-main">
        <Outlet />
      </main>
    </div>
  );
}
