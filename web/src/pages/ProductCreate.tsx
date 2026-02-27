import { Link } from 'react-router-dom';

export default function ProductCreate() {
  return (
    <div className="page">
      <p className="page-back"><Link to="/marketplace">← Marketplace</Link></p>
      <header className="page-header">
        <h1>Create listing</h1>
        <p className="page-intro">New product or service. API: POST /api/products.</p>
      </header>
      <div className="page-content">
        <p>Form to be wired.</p>
      </div>
    </div>
  );
}
