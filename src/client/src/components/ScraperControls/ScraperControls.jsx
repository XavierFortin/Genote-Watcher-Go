import {
  useStatus,
  usePostStartScraper,
  usePostStopScraper,
  usePostForceStartOnceScraper,
  usePostChangeInterval,
} from "@/services/hooks/scraperController.js";
import "./ScraperControls.css";
import { useState } from "react";
import { useEffect } from "react";

export function ScraperControls() {
  const { data, isLoading } = useStatus();

  const { mutate: startScraper } = usePostStartScraper();
  const { mutate: stopScraper } = usePostStopScraper();
  const { mutate: forceStartOnceScraper } = usePostForceStartOnceScraper();
  const { mutate: postSetInterval, error } = usePostChangeInterval();

  const { isRunning, interval } = data || {};
  const [newInterval, setNewInterval] = useState(interval);

  useEffect(() => {
    setNewInterval(interval);
  }, [interval]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    !isLoading && (
      <>
        <div className="scraper-controls-container">
          <div className="scraper-controls">
            {isRunning ? (
              <>
                <button
                  onClick={() => {
                    stopScraper();
                  }}
                  className="scraper-control-button stop-scraper"
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
                  className="scraper-control-button start-scraper"
                >
                  Start scraper
                </button>
              </>
            )}

            <button
              onClick={() => {
                forceStartOnceScraper();
              }}
              className="scraper-control-button force-start-once-scraper"
            >
              Force Start Scraper
            </button>

            <div>
              <label
                htmlFor="interval"
                className="block text-sm font-medium text-white"
              >
                Time Interval
              </label>
              <input
                type="text"
                id="interval"
                className="bg-gray-700 border border-gray-600 text-white text-sm rounded-lg block w-full p-2.5 
                 placeholder-gray-400 focus:ring-blue-500 focus:border-blue-500"
                placeholder="30s"
                required
                defaultValue={interval}
                onChange={({ target: { value } }) => {
                  if (value) {
                    setNewInterval(value);
                    console.log(value);
                  } else {
                    setNewInterval(0);
                  }
                }}
              />
            </div>
            <button
              onClick={() => {
                postSetInterval(newInterval);
              }}
              className="scraper-control-button set-interval"
            >
              Set Interval
            </button>
          </div>
          {error && (
            <div className="text-red-700">{error["response"].data}</div>
          )}
        </div>
      </>
    )
  );
}
