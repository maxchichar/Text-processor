import { useState, useEffect } from 'react';
import { getHistory, deleteFile } from '../api/client';

function HistoryList({ onHistoryItemClick }) {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await getHistory();
      setHistory(data);
    } catch (err) {
      setError(err.message || 'Failed to load history');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (fileId) => {
    if (!confirm('Are you sure you want to delete this file?')) return;

    try {
      await deleteFile(fileId);
      setHistory(history.filter((item) => item.fileId !== fileId));
    } catch (err) {
      setError(err.message || 'Failed to delete file');
    }
  };

  const formatDate = (timestamp) => {
    return new Date(timestamp).toLocaleString();
  };

  const formatSize = (bytes) => {
    return (bytes / 1024).toFixed(2) + ' KB';
  };

  const getEnabledTransformations = (transformations) => {
    return Object.entries(transformations)
      .filter(([_, enabled]) => enabled)
      .map(([key, _]) => key)
      .join(', ');
  };

  if (loading) {
    return (
      <div className="history-list loading">
        <div className="spinner"></div>
        <p>Loading history...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="history-list error">
        <p>{error}</p>
        <button onClick={loadHistory} className="retry-btn">Retry</button>
      </div>
    );
  }

  if (history.length === 0) {
    return (
      <div className="history-list empty">
        <p>No processing history yet</p>
      </div>
    );
  }

  return (
    <div className="history-list">
      <h2>Processing History</h2>
      <div className="history-items">
        {history.map((item) => (
          <div key={item.fileId} className="history-item">
            <div className="history-item-header">
              <h3>{item.filename}</h3>
              <span className="timestamp">{formatDate(item.timestamp)}</span>
            </div>

            <div className="history-item-details">
              <p><strong>Original:</strong> {formatSize(item.originalSize)}</p>
              <p><strong>Processed:</strong> {formatSize(item.processedSize)}</p>
              <p><strong>Transformations:</strong> {getEnabledTransformations(item.transformations)}</p>
            </div>

            <div className="history-item-actions">
              <button
                onClick={() => onHistoryItemClick(item)}
                className="action-btn view-btn"
              >
                View
              </button>
              <button
                onClick={() => handleDelete(item.fileId)}
                className="action-btn delete-btn"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default HistoryList;