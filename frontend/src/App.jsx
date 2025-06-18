import './index.css'; 
import { useState } from 'react';

function App() {
  // State variables for storing user inputs and UI state
  const [projectName, setProjectName] = useState('learn-go');
  const [moduleName, setModuleName] = useState('learn-go');
  const [db, setDb] = useState('postgres');
  const [framework, setFramework] = useState('echo');
  const [jwt, setJwt] = useState(false);
  const [swagger, setSwagger] = useState(false);
  const [redis, setRedis] = useState(false);
  const [validator, setValidator] = useState(false);
  const [loading, setLoading] = useState(false);
  const [darkMode, setDarkMode] = useState(false); // Toggle for dark/light mode

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const payload = {
      module_name: moduleName,
      project_name: projectName,
      db,
      framework,
      jwt,
      swagger,
      redis,
      validator,
    };

    console.log('Submitting payload:', payload);

    try {
      const response = await fetch('/api/v1/initialize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      console.log('Received response:', response);

      if (!response.ok) throw new Error('Failed to generate project');

      const blob = await response.blob();
      console.log('Received blob:', blob);

      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${projectName}.zip`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Error occurred:', err);
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`${darkMode ? 'bg-gray-900 text-white' : 'bg-gray-100 text-black'} min-h-screen p-8 font-sans`}>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-center w-full">Go Project Initializr</h1>
        <button
          onClick={() => setDarkMode(!darkMode)}
          className="absolute right-8 bg-gray-300 dark:bg-gray-700 text-sm px-4 py-1 rounded"
        >
          {darkMode ? 'Light Mode' : 'Dark Mode'}
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Form Section */}
        <form onSubmit={handleSubmit} className={`${darkMode ? 'bg-gray-800' : 'bg-white'} p-6 rounded-xl shadow-md`}>
          <div className="mb-4">
            <label className="block mb-1 font-medium">Project Name</label>
            <input
              value={projectName}
              onChange={e => setProjectName(e.target.value)}
              className="w-full border px-3 py-2 rounded text-black"
            />
          </div>

          <div className="mb-4">
            <label className="block mb-1 font-medium">Module Name</label>
            <input
              value={moduleName}
              onChange={e => setModuleName(e.target.value)}
              className="w-full border px-3 py-2 rounded text-black"
            />
          </div>

          <div className="mb-4">
            <label className="block mb-1 font-medium">Database</label>
            <select
              value={db}
              onChange={e => setDb(e.target.value)}
              className="w-full border px-3 py-2 rounded text-black"
            >
              <option value="postgres">Postgres</option>
              <option value="mysql">MySQL</option>
            </select>
          </div>

          <div className="mb-4">
            <label className="block mb-1 font-medium">Framework</label>
            <select
              value={framework}
              onChange={e => setFramework(e.target.value)}
              className="w-full border px-3 py-2 rounded text-black"
            >
              <option value="echo">Echo</option>
            </select>
          </div>

          <div className="mb-4 grid grid-cols-2 gap-2">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={jwt}
                onChange={e => setJwt(e.target.checked)}
              /> JWT
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={swagger}
                onChange={e => setSwagger(e.target.checked)}
              /> Swagger
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={redis}
                onChange={e => setRedis(e.target.checked)}
              /> Redis
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={validator}
                onChange={e => setValidator(e.target.checked)}
              /> Validator
            </label>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="mt-4 w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
          >
            {loading ? 'Generating...' : 'Generate ZIP'}
          </button>
        </form>

        {/* Right Panel: Setup Instructions */}
        <div className={`${darkMode ? 'bg-gray-800' : 'bg-white'} p-6 rounded-xl shadow-md text-sm leading-relaxed`}>
          <h2 className="text-xl font-bold mb-4">Project Setup Instructions</h2>
          <p>
            1. Unzip the downloaded archive.<br />
            2. Navigate into the project folder.<br />
            3. Run <code>go mod tidy</code> to install dependencies.<br />
            4. Setup your database based on <code>.env</code> configuration.<br />
            5. Run <code>go run main.go</code> to start the server.<br /><br />
            Optional:<br />
            - Use Swagger UI at <code>/swagger/index.html</code> if enabled.<br />
            - Make sure Redis is running if selected.<br />
            - JWT support requires token handling middleware configuration.
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;
