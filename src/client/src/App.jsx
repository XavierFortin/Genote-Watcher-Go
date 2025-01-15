import "./App.css";
import {
  forceStartOnceScraper,
  restartScraper,
  startScraper,
  stopScraper,
} from "./services/GenoteController";

function App() {
  return (
    <>
      <h1>Vite + React</h1>
      <div className="card">
        <button
          onClick={() => {
            startScraper();
          }}
          id="start-scraper"
        >
          Start scraper
        </button>

        <button
          onClick={() => {
            stopScraper();
          }}
          id="stop-scraper"
        >
          Stop scraper
        </button>

        <button
          onClick={() => {
            restartScraper();
          }}
          id="restart-scraper"
        >
          Restart Scraper
        </button>

        <button
          onClick={() => {
            forceStartOnceScraper();
          }}
          id="force-start-once-scraper"
        >
          Force Start Scraper
        </button>
      </div>
    </>
  );
}

export default App;
