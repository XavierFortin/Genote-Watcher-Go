import {
  useStatus,
  usePostStartScraper,
  usePostStopScraper,
  usePostForceStartOnceScraper,
  usePostRestartScraper,
} from "@/services/hooks/scraperController.js";
import "./ScraperControls.css";

export function ScraperControls() {
  const { data, isLoading } = useStatus();

  const { isRunning, interval } = data || {};
  const { mutate: startScraper } = usePostStartScraper();
  const { mutate: stopScraper } = usePostStopScraper();
  const { mutate: forceStartOnceScraper } = usePostForceStartOnceScraper();
  const { mutate: restartScraper } = usePostRestartScraper();

  return (
    !isLoading && (
      <div className="scraper-controls-container">
        <div className="scraper-controls">
          {isRunning ? (
            <>
              <button
                onClick={() => {
                  stopScraper();
                }}
                className="stop-scraper"
              >
                Stop scraper
              </button>
            </>
          ) : (
            <>
              <button
                onClick={() => {
                  startScraper();
                }}
                className="start-scraper"
              >
                Start scraper
              </button>
            </>
          )}

          <button
            onClick={() => {
              restartScraper();
            }}
            className="restart-scraper"
          >
            Restart Scraper
          </button>

          <button
            onClick={() => {
              forceStartOnceScraper();
            }}
            className="force-start-once-scraper"
          >
            Force Start Scraper
          </button>
        </div>
        <div className="scraper-controls">
          <input name="interval" type="text" defaultValue={interval} />
        </div>
      </div>
    )
  );
}
