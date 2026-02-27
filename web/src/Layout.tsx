import { Link, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from './auth';
import './Layout.css';

const root = typeof window !== 'undefined' ? window.location.origin : '';

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
            <Link to="/">Dashboard</Link>
            <Link to="/marketplace">Trove · marketplace</Link>
            <Link to="/mail">Relay · mail</Link>
            <Link to="/orders">Orders</Link>
            <Link to="/find-professional">Find professional</Link>
            <Link to="/notifications">Notifications</Link>
            <Link to="/profile">Profile</Link>
            <Link to="/wallet">Wallet</Link>
            <Link to="/remittances">Remittances</Link>
            <Link to="/trade">Trade</Link>
            <Link to="/connect">Connect</Link>
            <Link to="/vault">Vault</Link>
            <Link to="/direction/ixi">IXI</Link>
            <Link to="/direction/learning">Learning</Link>
            <Link to="/direction/repositorium">Repositorium</Link>
            <Link to="/direction/startups">Startups</Link>
            <Link to="/direction/media">Media</Link>
            <Link to="/direction/rewards">Rewards</Link>
            <Link to="/ai">AI</Link>
            <Link to="/settings">Settings</Link>
            <Link to="/admin">Admin</Link>
            {token ? (
              <>
                <span className="layout-user">{user?.email ?? user?.name ?? 'User'}</span>
                <button type="button" className="layout-btn" onClick={handleLogout}>Sign out</button>
              </>
            ) : (
              <Link to="/login">Sign in</Link>
            )}
          </nav>
          <div className="layout-nav-external">
            <a href={`${root}/ecosystem.html`}>Ecosystem</a>
            <a href={`${root}/architecture.html`}>Architecture</a>
            <a href={`${root}/contact.html`}>Contact</a>
          </div>
        </div>
      </header>
      <main className="layout-main">
        <Outlet />
      </main>
    </div>
  );
}
