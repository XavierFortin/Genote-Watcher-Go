async function Get(url) {
  var result = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return await result.json();
}

export function stopScraper() {
  return Get("http://localhost:4000/api/scraper/stop");
}
