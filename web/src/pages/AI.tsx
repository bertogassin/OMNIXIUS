import { useState } from 'react';

export default function AI() {
  const [message, setMessage] = useState('');
  const [reply, setReply] = useState('');

  const send = () => {
    setReply('AI chat: backend ai/ (Python) — connect API and token. Screen preserved.');
  };

  return (
    <div className="page">
      <header className="page-header">
        <h1>Oracle · AI</h1>
        <p className="page-intro">Root of our AI: chat, then own models.</p>
      </header>
      <div className="page-content">
        <section className="page-section">
          <div className="page-form-row">
            <label>
              <input value={message} onChange={(e) => setMessage(e.target.value)} placeholder="Message" />
            </label>
            <button type="button" className="btn btn-primary" onClick={send}>Send</button>
          </div>
          {reply && <p>{reply}</p>}
        </section>
      </div>
    </div>
  );
}
