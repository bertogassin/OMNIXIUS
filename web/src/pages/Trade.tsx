import { Link } from 'react-router-dom';

export default function Trade() {
  return (
    <div>
      <h1>Forge · Trade</h1>
      <p>Exchange, wallets. <Link to="/wallet">Wallet</Link></p>
    </div>
  );
}
