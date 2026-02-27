import { Link } from 'react-router-dom';

export default function ProductCreate() {
  return (
    <div>
      <h1>Create listing</h1>
      <p>New product or service. API: POST /api/products. Form to be wired. Screen preserved.</p>
      <p><Link to="/marketplace">← Marketplace</Link></p>
    </div>
  );
}
