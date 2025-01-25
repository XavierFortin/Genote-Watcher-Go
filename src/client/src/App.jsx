import "./App.css";
import { LogViewer } from "./components/LogViewer/LogViewer.jsx";
import { ScraperControls } from "./components/ScraperControls/ScraperControls.jsx";

function App() {
  return (
    <>
      <h1 className="app-title">
        <span className="color-fade">Genote Watcher</span>
      </h1>
      <ScraperControls />

      <LogViewer />
    </>
  );
}

export default App;
