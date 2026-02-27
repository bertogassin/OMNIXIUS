import { Link } from 'react-router-dom';

export default function Trade() {
  return (
    <div className="page">
      <header className="page-header">
        <h1>Forge · Trade</h1>
        <p className="page-intro">Exchange, wallets, copy trading.</p>
      </header>
      <div className="page-content">
        <p><Link to="/wallet">Wallet</Link> — balances and transfers.</p>
      </div>
    </div>
  );
}
