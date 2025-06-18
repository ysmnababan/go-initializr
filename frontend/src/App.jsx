import './index.css'; 
import { useState, useEffect } from 'react';
import { Toaster, toast } from 'react-hot-toast';

function App() {
  const [projectName, setProjectName] = useState('learn-go');
  const [moduleName, setModuleName] = useState('learn-go');
  const [db, setDb] = useState('postgres');
  const [framework, setFramework] = useState('echo');
  const [jwt, setJwt] = useState(false);
  const [swagger, setSwagger] = useState(false);
  const [redis, setRedis] = useState(false);
  const [validator, setValidator] = useState(false);
  const [loading, setLoading] = useState(false);
  const [darkMode, setDarkMode] = useState(false);

  const validProjectName = /^[a-zA-Z0-9_-]{3,64}$/; // validation regex for project name
  const validModuleName = /^[a-z0-9\-\/.]{3,64}$/; // validation regex for module name

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!validProjectName.test(projectName)) {
      toast.error('Project name must be 3-64 characters and contain only letters, numbers, underscores, or dashes.');
      return;
    }

    if (!validModuleName.test(moduleName)) {
      toast.error('Module name must be 3-64 characters and contain only lowercase letters, numbers, hyphens, slashes, or dots.');
      return;
    }

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

    try {
      const response = await fetch('/api/v1/initialize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const resData = await response.json().catch(() => ({}));
        const msg = resData?.meta?.message || resData?.message || 'Something went wrong';
        throw new Error(response.status >= 500 ? 'Internal server error' : msg);
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${projectName}.zip`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      toast.error(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.ctrlKey && e.key === 'Enter') {
        handleSubmit(e);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [projectName, moduleName, db, framework, jwt, swagger, redis, validator]);

  return (
    <div className={`${darkMode ? 'bg-gray-900 text-white' : 'bg-gray-100 text-black'} min-h-screen p-4 md:p-8 font-sans`}>
      <Toaster position="top-right" reverseOrder={false} />
      <div className="relative mb-6">
        <div className="text-left">
          <h1 className="text-4xl font-bold flex items-center">
            <img src="/go-initializr-icon.svg" alt="Go Initializr Logo" className="h-24 w-24" />
            Go Initializr
          </h1>
        </div>
        <div className="absolute top-0 right-0">
          <button
            onClick={() => setDarkMode(!darkMode)}
            className="bg-gray-300 dark:bg-gray-700 text-sm px-4 py-1 rounded"
          >
            {darkMode ? 'Light Mode' : 'Dark Mode'}
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-8">
        <form onSubmit={handleSubmit} className={`${darkMode ? 'bg-gray-800' : 'bg-white'} p-4 md:p-6 rounded-xl shadow-md`}>
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
            className="mt-4 w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 flex items-center justify-center gap-3 text-xl"
          >
            {loading ? 'Generating...' : 'Generate'}
            {!loading && (
              <span className="ml-2 flex items-center gap-1">
                <kbd className="bg-gray-200 text-gray-800 px-2 py-0.5 rounded border border-gray-300 text-xs font-mono shadow-inner">Ctrl</kbd>
                <span className="text-xs font-mono">+</span>
                <kbd className="bg-gray-200 text-gray-800 px-2 py-0.5 rounded border border-gray-300 text-xs font-mono shadow-inner">Enter</kbd>
              </span>
            )}
          </button>
        </form>

        <div className={`${darkMode ? 'bg-gray-800' : 'bg-white'} p-4 md:p-6 rounded-xl shadow-md text-sm leading-relaxed`}>
          <h2 className="text-xl font-bold mb-4">Project Setup Instructions</h2>
          <p>
            1. Unzip the downloaded archive.<br />
            2. Navigate into the project folder.<br />
            3. Run <code>go mod tidy</code> to install dependencies.<br />
            4. Setup your database based on <code>.env</code> configuration.<br />
            5. Run <code>go run main.go</code> to start the server.<br /><br />
            Optional:<br />
            - Run <code>go get github.com/swaggo/swag@latest</code> <strong>only if</strong> Swagger capability was selected.<br />
            - Run <code>swag init</code> <strong>only if</strong> Swagger capability was selected.<br />
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
