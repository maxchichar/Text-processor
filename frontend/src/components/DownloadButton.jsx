import { downloadFile } from '../api/client';

function DownloadButton({ fileId, filename, disabled }) {
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState('');

  const handleDownload = async () => {
    if (!fileId || disabled) return;

    setDownloading(true);
    setError('');

    try {
      await downloadFile(fileId, filename);
    } catch (err) {
      setError(err.message || 'Failed to download file');
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="download-button">
      <button
        onClick={handleDownload}
        disabled={!fileId || disabled || downloading}
        className="download-btn"
      >
        {downloading ? (
          <>
            <div className="spinner-small"></div>
            Downloading...
          </>
        ) : (
          <>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="7 10 12 15 17 10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
            Download Processed File
          </>
        )}
      </button>

      {error && <div className="error">{error}</div>}
    </div>
  );
}

export default DownloadButton;