function TransformationConfig({ options, onOptionsChange }) {
  const handleChange = (transformation) => {
    onOptionsChange({
      ...options,
      [transformation]: !options[transformation],
    });
  };

  const transformations = [
    { id: 'tokenize', name: 'Tokenization', description: 'Split text into tokens' },
    { id: 'hex', name: 'Hex to Decimal', description: 'Convert hex numbers to decimal' },
    { id: 'bin', name: 'Binary to Decimal', description: 'Convert binary numbers to decimal' },
    { id: 'case', name: 'Case Transformation', description: 'Apply case changes (up/low/cap)' },
    { id: 'quote', name: 'Quote Normalization', description: 'Normalize quote marks' },
    { id: 'punctuation', name: 'Punctuation Normalization', description: 'Fix punctuation spacing' },
    { id: 'article', name: 'Article Correction', description: 'Fix a/an usage' },
  ];

  return (
    <div className="transformation-config">
      <h2>Transformation Options</h2>
      <div className="transformations-list">
        {transformations.map((trans) => (
          <div key={trans.id} className="transformation-item">
            <label className="checkbox-label">
              <input
                type="checkbox"
                checked={options[trans.id]}
                onChange={() => handleChange(trans.id)}
              />
              <span className="checkbox-custom"></span>
              <div className="transformation-info">
                <span className="transformation-name">{trans.name}</span>
                <span className="transformation-description">{trans.description}</span>
              </div>
            </label>
          </div>
        ))}
      </div>
      <div className="transformations-summary">
        <p>
          <strong>{Object.values(options).filter(Boolean).length}</strong> of{' '}
          {Object.keys(options).length} transformations enabled
        </p>
      </div>
    </div>
  );
}

export default TransformationConfig;