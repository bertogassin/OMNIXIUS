import { useParams } from 'react-router-dom';

export default function Conversation() {
  const { id } = useParams<{ id: string }>();
  return (
    <div>
      <h1>Conversation {id}</h1>
      <p>Chat with seller/buyer. Messages API wired — UI to be completed. Screen preserved.</p>
    </div>
  );
}
