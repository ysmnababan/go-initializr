import { useState } from 'react';

function App() {
  const [projectName, setProjectName] = useState('learn-go-otel');
  const [db, setDb] = useState('postgres');
  const [framework, setFramework] = useState('echo');
  const [jwt, setJwt] = useState(false);
  const [swagger, setSwagger] = useState(false);
  const [redis, setRedis] = useState(false);
  const [validator, setValidator] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const payload = {
      project_name: projectName,
      db,
      framework,
      jwt,
      swagger,
      redis,
      validator,
    };

    try {
      const response = await fetch('/api/v1/initialize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error('Failed to generate project');

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);

      const a = document.createElement('a');
      a.href = url;
      a.download = `${projectName}.zip`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '2rem', fontFamily: 'sans-serif' }}>
      <h2>Go Project Initializr</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Project Name: </label>
          <input value={projectName} onChange={e => setProjectName(e.target.value)} />
        </div>

        <div>
          <label>DB: </label>
          <select value={db} onChange={e => setDb(e.target.value)}>
            <option value="postgres">Postgres</option>
            <option value="mysql">MySQL</option>
          </select>
        </div>

        <div>
          <label>Framework: </label>
          <select value={framework} onChange={e => setFramework(e.target.value)}>
            <option value="echo">Echo</option>
          </select>
        </div>

        <div>
          <label><input type="checkbox" checked={jwt} onChange={e => setJwt(e.target.checked)} /> JWT</label>
          <label><input type="checkbox" checked={swagger} onChange={e => setSwagger(e.target.checked)} /> Swagger</label>
          <label><input type="checkbox" checked={redis} onChange={e => setRedis(e.target.checked)} /> Redis</label>
          <label><input type="checkbox" checked={validator} onChange={e => setValidator(e.target.checked)} /> Validator</label>
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Generating...' : 'Generate ZIP'}
        </button>
      </form>
    </div>
  );
}

export default App;
