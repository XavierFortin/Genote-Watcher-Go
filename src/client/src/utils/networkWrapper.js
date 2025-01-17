export async function get(url) {
  var result = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return await result.json();
}

export async function post(url, data = null) {
  console.log(data);
  var result = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return result;
}
