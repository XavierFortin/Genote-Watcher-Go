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
              forceStartOnceScraper();
            }}
            className="force-start-once-scraper"
          >
            Force Start Scraper
          </button>
        </div>
        <div className="scraper-controls">
          <label htmlFor="interval" style={{ alignSelf: "center" }}>
            Interval:
          </label>
          <input
            name="interval"
            type="text"
            style={{ height: "30px", width: "50px", alignSelf: "center" }}
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
          <button
            onClick={() => {
              postSetInterval(newInterval);
            }}
            className="set-interval"
          >
            Set Interval
          </button>
        </div>
        {error && (
          <div
            className="error"
            style={{
              color: "red",
            }}
          >
            {error["response"].data}
          </div>
        )}
      </div>
    )
  );
}
