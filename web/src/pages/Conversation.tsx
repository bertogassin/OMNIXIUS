import { useParams, Link } from 'react-router-dom';

export default function Conversation() {
  const { id } = useParams<{ id: string }>();
  return (
    <div className="page">
      <p className="page-back"><Link to="/mail">← Mail</Link></p>
      <header className="page-header">
        <h1>Conversation {id}</h1>
        <p className="page-intro">Chat with seller/buyer.</p>
      </header>
      <div className="page-content">
        <p>Messages API wired — UI to be completed.</p>
      </div>
    </div>
  );
}
