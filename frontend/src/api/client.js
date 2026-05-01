const API_BASE_URL = '/api';

/**
 * Upload a file to the server
 * @param {File} file - The file to upload
 * @returns {Promise<Object>} - Response with fileId, filename, and size
 */
export async function uploadFile(file) {
  console.log('API: Uploading file:', file.name);
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: 'POST',
    body: formData,
  });

  console.log('API: Upload response status:', response.status);
  if (!response.ok) {
    const error = await response.json();
    console.error('API: Upload error:', error);
    throw new Error(error.error || 'Failed to upload file');
  }

  const result = await response.json();
  console.log('API: Upload result:', result);
  return result;
}

/**
 * Process a file with transformation options
 * @param {string} fileId - The ID of the file to process
 * @param {Object} transformations - Transformation options
 * @returns {Promise<Object>} - Response with fileId and status
 */
export async function processFile(fileId, transformations) {
  console.log('API: Processing file:', fileId, 'with transformations:', transformations);
  const response = await fetch(`${API_BASE_URL}/process/${fileId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(transformations),
  });

  console.log('API: Process response status:', response.status);
  if (!response.ok) {
    const error = await response.json();
    console.error('API: Process error:', error);
    throw new Error(error.error || 'Failed to process file');
  }

  const result = await response.json();
  console.log('API: Process result:', result);
  return result;
}

/**
 * Get preview of original and processed text
 * @param {string} fileId - The ID of the file
 * @returns {Promise<Object>} - Response with original and processed text
 */
export async function getPreview(fileId) {
  console.log('API: Getting preview for file:', fileId);
  const response = await fetch(`${API_BASE_URL}/preview/${fileId}`);

  console.log('API: Preview response status:', response.status);
  if (!response.ok) {
    const error = await response.json();
    console.error('API: Preview error:', error);
    throw new Error(error.error || 'Failed to get preview');
  }

  const result = await response.json();
  console.log('API: Preview result lengths - original:', result.original?.length, 'processed:', result.processed?.length);
  return result;
}

/**
 * Download processed file
 * @param {string} fileId - The ID of the file
 * @param {string} filename - Original filename
 */
export async function downloadFile(fileId, filename) {
  const response = await fetch(`${API_BASE_URL}/download/${fileId}`);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to download file');
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `processed-${filename}`;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

/**
 * Get processing history
 * @returns {Promise<Array>} - Array of history entries
 */
export async function getHistory() {
  const response = await fetch(`${API_BASE_URL}/history`);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to get history');
  }

  const data = await response.json();
  return data.history;
}

/**
 * Delete a file and its processed version
 * @param {string} fileId - The ID of the file to delete
 * @returns {Promise<Object>} - Response with success status
 */
export async function deleteFile(fileId) {
  const response = await fetch(`${API_BASE_URL}/files/${fileId}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete file');
  }

  return response.json();
}