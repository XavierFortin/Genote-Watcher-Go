async function Get(url) {
  var result = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return await result.json();
}

export function helloWorld() {
  return Get("http://localhost:4000/api/hello");
}

export function stopScraper() {
  return Get("http://localhost:4000/api/scraper/stop");
}

export function startScraper() {
  return Get("http://localhost:4000/api/scraper/start");
}

export function forceStartOnceScraper() {
  return Get("http://localhost:4000/api/scraper/force-start");
}

export function restartScraper() {
  return Get("http://localhost:4000/api/scraper/restart");
}
