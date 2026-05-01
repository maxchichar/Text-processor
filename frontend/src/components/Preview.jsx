import { useState, useEffect, useRef } from 'react';
import { getPreview } from '../api/client';

function Preview({ fileId, filename }) {
  const [preview, setPreview] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [syncScroll, setSyncScroll] = useState(true);
  const originalRef = useRef(null);
  const processedRef = useRef(null);

  useEffect(() => {
    if (fileId) {
      console.log('Preview component - Loading preview for fileId:', fileId);
      loadPreview();
    }
  }, [fileId]);

  const loadPreview = async () => {
    setLoading(true);
    setError('');
    try {
      console.log('Calling getPreview for fileId:', fileId);
      const data = await getPreview(fileId);
      console.log('Preview data received:', data);
      setPreview(data);
    } catch (err) {
      console.error('Preview error:', err);
      setError(err.message || 'Failed to load preview');
    } finally {
      setLoading(false);
    }
  };

  const handleScroll = (source) => {
    if (!syncScroll) return;

    const sourceElement = source === 'original' ? originalRef.current : processedRef.current;
    const targetElement = source === 'original' ? processedRef.current : originalRef.current;

    if (sourceElement && targetElement) {
      const scrollPercentage = sourceElement.scrollTop / (sourceElement.scrollHeight - sourceElement.clientHeight);
      targetElement.scrollTop = scrollPercentage * (targetElement.scrollHeight - targetElement.clientHeight);
    }
  };

  if (loading) {
    console.log('Preview component - Loading state');
    return (
      <div className="preview loading" style={{background: '#d1ecf1', padding: '20px', border: '2px solid #0c5460', borderRadius: '5px'}}>
        <div className="spinner"></div>
        <p>🔄 Loading preview...</p>
      </div>
    );
  }

  if (error) {
    console.log('Preview component - Error state:', error);
    return (
      <div className="preview error" style={{background: '#f8d7da', padding: '20px', border: '2px solid #721c24', borderRadius: '5px'}}>
        <p>❌ Error: {error}</p>
      </div>
    );
  }

  if (!preview) {
    console.log('Preview component - Empty state, fileId:', fileId);
    return (
      <div className="preview empty" style={{background: '#fff3cd', padding: '20px', border: '2px solid #ffc107', borderRadius: '5px'}}>
        <p>⚠️ Upload and process a file to see the preview</p>
        <p style={{fontSize: '12px', marginTop: '10px'}}>Current fileId: {fileId || 'None'}</p>
      </div>
    );
  }

  console.log('Preview component - Rendering preview with', preview.original.length, 'original chars and', preview.processed.length, 'processed chars');

  const lineCount = Math.max(
    preview.original.split('\n').length,
    preview.processed.split('\n').length
  );

  return (
    <div className="preview" style={{background: '#d4edda', padding: '30px', border: '3px solid #28a745', borderRadius: '10px'}}>
      <div className="preview-header">
        <h2 style={{color: '#155724'}}>✅ Preview - File Successfully Processed!</h2>
        <label className="sync-scroll">
          <input
            type="checkbox"
            checked={syncScroll}
            onChange={(e) => setSyncScroll(e.target.checked)}
          />
          <span>Sync scroll</span>
        </label>
      </div>

      <div className="preview-content">
        <div className="preview-panel">
          <h3>Original</h3>
          <div className="preview-text" ref={originalRef} onScroll={() => handleScroll('original')}>
            {preview.original.split('\n').map((line, i) => (
              <div key={i} className="preview-line">
                <span className="line-number">{i + 1}</span>
                <span className="line-content">{line || ' '}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="preview-panel">
          <h3>Processed</h3>
          <div className="preview-text" ref={processedRef} onScroll={() => handleScroll('processed')}>
            {preview.processed.split('\n').map((line, i) => (
              <div key={i} className="preview-line">
                <span className="line-number">{i + 1}</span>
                <span className="line-content">{line || ' '}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="preview-stats">
        <p><strong>Lines:</strong> {lineCount}</p>
        <p><strong>Original size:</strong> {preview.original.length} characters</p>
        <p><strong>Processed size:</strong> {preview.processed.length} characters</p>
      </div>
    </div>
  );
}

export default Preview;