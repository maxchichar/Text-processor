import { useState } from 'react';
import FileUpload from './components/FileUpload';
import TransformationConfig from './components/TransformationConfig';
import Preview from './components/Preview';
import DownloadButton from './components/DownloadButton';
import HistoryList from './components/HistoryList';
import { processFile } from './api/client';
import './App.css';

function App() {
  const [currentFile, setCurrentFile] = useState(null);
  const [processedFile, setProcessedFile] = useState(null);
  const [transformations, setTransformations] = useState({
    tokenize: true,
    hex: true,
    bin: true,
    case: true,
    quote: true,
    punctuation: true,
    article: true,
  });
  const [processing, setProcessing] = useState(false);
  const [error, setError] = useState('');

  const handleFileUploaded = (fileData) => {
    console.log('App: File uploaded:', fileData);
    setCurrentFile(fileData);
    setProcessedFile(null);
    setError('');
  };

  const handleProcess = async () => {
    if (!currentFile) {
      setError('Please upload a file first');
      return;
    }

    setProcessing(true);
    setError('');

    try {
      console.log('App: Processing file:', currentFile.fileId);
      const result = await processFile(currentFile.fileId, transformations);
      console.log('App: Process result:', result);
      const newProcessedFile = {
        ...result,
        filename: currentFile.filename,
      };
      console.log('App: Setting processedFile:', newProcessedFile);
      setProcessedFile(newProcessedFile);
    } catch (err) {
      console.error('App: Process error:', err);
      setError(err.message || 'Failed to process file');
    } finally {
      setProcessing(false);
    }
  };

  const handleHistoryItemClick = (item) => {
    console.log('App: History item clicked:', item);
    setCurrentFile({
      fileId: item.fileId,
      filename: item.filename,
      size: item.originalSize,
    });
    const newProcessedFile = {
      fileId: item.fileId,
      filename: item.filename,
      status: 'processed',
    };
    console.log('App: Setting processedFile from history:', newProcessedFile);
    setProcessedFile(newProcessedFile);
    setTransformations(item.transformations);
    setError('');
  };

  const handleReset = () => {
    console.log('App: Resetting');
    setCurrentFile(null);
    setProcessedFile(null);
    setError('');
  };

  console.log('App: Render state - currentFile:', currentFile, 'processedFile:', processedFile);

  return (
    <div className="app">
      <header className="app-header">
        <h1>Text Processor</h1>
        <p>Upload, transform, and download your text files</p>
      </header>

      <main className="app-main">
        {error && <div className="error-banner">{error}</div>}

        {!currentFile ? (
          <div className="upload-section">
            <FileUpload onFileUploaded={handleFileUploaded} />
          </div>
        ) : (
          <div className="process-section">
            <div className="file-info">
              <h2>Current File</h2>
              <p><strong>Filename:</strong> {currentFile.filename}</p>
              <p><strong>Size:</strong> {(currentFile.size / 1024).toFixed(2)} KB</p>
              <button onClick={handleReset} className="reset-btn">
                Upload Different File
              </button>
            </div>

            <TransformationConfig
              options={transformations}
              onOptionsChange={setTransformations}
            />

            <div className="process-actions">
              <button
                onClick={handleProcess}
                disabled={processing}
                className="process-btn"
              >
                {processing ? (
                  <>
                    <div className="spinner-small"></div>
                    Processing...
                  </>
                ) : (
                  'Process File'
                )}
              </button>

              {processedFile && (
                <>
                  {console.log('App: Rendering DownloadButton')}
                  <DownloadButton
                    fileId={processedFile.fileId}
                    filename={processedFile.filename}
                    disabled={false}
                  />
                </>
              )}
            </div>

            {/* Debug info */}
            <div style={{background: '#f0f0f0', padding: '10px', margin: '10px 0', borderRadius: '5px', fontSize: '12px'}}>
              <strong>Debug Info:</strong>
              <p>Current File: {currentFile ? currentFile.filename : 'None'}</p>
              <p>Processed File: {processedFile ? processedFile.fileId : 'None'}</p>
              <p>Processing: {processing ? 'Yes' : 'No'}</p>
            </div>

            {processedFile && (
              <>
                {console.log('App: Rendering Preview with fileId:', processedFile.fileId)}
                <Preview
                  fileId={processedFile.fileId}
                  filename={processedFile.filename}
                />
              </>
            )}
          </div>
        )}

        <div className="history-section">
          <HistoryList onHistoryItemClick={handleHistoryItemClick} />
        </div>
      </main>

      <footer className="app-footer">
        <p>Text Processor - Transform your text with powerful processing options</p>
      </footer>
    </div>
  );
}

export default App;