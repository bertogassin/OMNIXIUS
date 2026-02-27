import { useState } from 'react';

export default function AI() {
  const [message, setMessage] = useState('');
  const [reply, setReply] = useState('');

  const send = () => {
    setReply('AI chat: backend ai/ (Python) — connect API and token. Screen preserved.');
  };

  return (
    <div>
      <h1>Oracle · AI</h1>
      <p>Root of our AI: chat, then own models.</p>
      <input value={message} onChange={(e) => setMessage(e.target.value)} placeholder="Message" />
      <button type="button" className="btn btn-primary" onClick={send}>Send</button>
      {reply && <p>{reply}</p>}
    </div>
  );
}
