import React, { useState } from 'react';
import ScrollToBottom, { useScrollToBottom, useSticky } from 'react-scroll-to-bottom';
import './App.css';

const MsgList = ({ history, client }) => {
  const scrollToBottom = useScrollToBottom();
  const [sticky] = useSticky();

  return (
    <React.Fragment>
      <ul className="msg-list">
        {history.map((m,j) =>
          <li key={`msg-${j}`} className="msg-list-item">
            <span style={m.authenticated ? {} : {color: "red"}}>{`${client}: ${new Date(m.pinged_at).toISOString()}`}</span>
          </li>
        )}
       { !sticky && <button onClick={ scrollToBottom }>Click me to scroll to bottom</button> }
      </ul>
    </React.Fragment>
  );
}

const fetchHistory = setHistory => {
  setTimeout(() => {
    const url = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080'
    fetch(`${url}/history`, {method: "GET"})
      .then(h => h.json()).then(h => setHistory(h))
      .catch(e => console.log(e))
  }, 1000);
}

function App() {
  const [history, setHistory] = useState({})
  fetchHistory(setHistory)

  return (
    <div className="container">
        {Object.keys(history).map((c,i) =>
          <div key={`pub-${i}`}>
            <span>{c}</span>
            <div className="msg-list-wrapper">
              <ScrollToBottom>
                <MsgList history={history[c]} client={c}/>
              </ScrollToBottom>
            </div>
          </div>
        )}
    </div>
  );
}

export default App;
