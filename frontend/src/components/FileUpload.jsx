import { useState } from 'react';
import { uploadFile } from '../api/client';

function FileUpload({ onFileUploaded }) {
  const [dragActive, setDragActive] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  const [uploadedFile, setUploadedFile] = useState(null);

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFile(e.dataTransfer.files[0]);
    }
  };

  const handleChange = (e) => {
    e.preventDefault();
    if (e.target.files && e.target.files[0]) {
      handleFile(e.target.files[0]);
    }
  };

  const handleFile = async (file) => {
    setError('');
    setUploading(true);

    // Validate file type
    if (!file.name.endsWith('.txt')) {
      setError('Only .txt files are allowed');
      setUploading(false);
      return;
    }

    try {
      const result = await uploadFile(file);
      console.log('File upload result:', result);
      setUploadedFile(result);
      onFileUploaded(result);
      console.log('Called onFileUploaded with:', result);
    } catch (err) {
      console.error('File upload error:', err);
      setError(err.message || 'Failed to upload file');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="file-upload">
      <h2>Upload File</h2>
      <div
        className={`drop-zone ${dragActive ? 'active' : ''}`}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
      >
        <input
          type="file"
          id="file-upload"
          accept=".txt"
          onChange={handleChange}
          disabled={uploading}
        />
        <label htmlFor="file-upload" className="drop-zone-content">
          {uploading ? (
            <div className="uploading">
              <div className="spinner"></div>
              <p>Uploading...</p>
            </div>
          ) : (
            <div className="upload-prompt">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="17 8 12 3 7 8" />
                <line x1="12" y1="3" x2="12" y2="15" />
              </svg>
              <p>Drag and drop a .txt file here, or click to select</p>
            </div>
          )}
        </label>
      </div>

      {error && <div className="error">{error}</div>}

      {uploadedFile && (
        <div className="uploaded-file">
          <h3>File Uploaded</h3>
          <p><strong>Filename:</strong> {uploadedFile.filename}</p>
          <p><strong>Size:</strong> {(uploadedFile.size / 1024).toFixed(2)} KB</p>
          <p><strong>File ID:</strong> {uploadedFile.fileId}</p>
        </div>
      )}
    </div>
  );
}

export default FileUpload;