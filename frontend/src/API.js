// Thin wrapper over the browser's fetch API, keeping the axios-like
// response shape ({ data }) so callers can keep reading res.data.
// The server reports failures as 4xx/5xx with a {"error": "..."} body;
// those (and network errors) log the detail and reject, so callers'
// .then() never runs on a failed request.
async function request(method, url, args) {

  let res;
  try {
    res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest'
      },
      body: args === undefined ? undefined : JSON.stringify(args)
    });
  } catch (err) {
    console.error("API error:", err.message);
    throw err;
  }

  if (!res.ok) {
    let detail;
    try {
      const body = await res.json();
      detail = body && body.error;
    } catch {
      // not a JSON body — fall back to the status line
    }
    const err = new Error(detail || res.status + " " + res.statusText);
    err.status = res.status;
    console.error("API error:", err.message);
    throw err;
  }

  return { data: await res.json() };
}

class API {

  static post(url, args) {
    return request("POST", url, args);
  }

  static patch(url, args) {
    return request("PATCH", url, args);
  }

  static delete(url, args) {
    return request("DELETE", url, args);
  }
}

export default API;
