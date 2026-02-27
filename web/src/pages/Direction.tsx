import { useParams, Link } from 'react-router-dom';

const titles: Record<string, string> = {
  ixi: 'IXI · Blockchain',
  learning: 'Ascent · Learning',
  media: 'Lens · Media',
  rewards: 'Bounty · Rewards',
  startups: 'Flare · Startups',
  repositorium: 'Ark · Repositorium',
};

export default function Direction() {
  const { name } = useParams<{ name: string }>();
  const title = (name && titles[name]) || name || 'Direction';
  return (
    <div className="page">
      <header className="page-header">
        <h1>{title}</h1>
        <p className="page-intro">Content and API to be wired.</p>
      </header>
      <div className="page-content">
        <p className="page-back"><Link to="/">← Dashboard</Link></p>
      </div>
    </div>
  );
}
