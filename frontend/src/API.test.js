import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import API from "./API";

function jsonResponse(body, ok = true, status = 200) {
  return {
    ok,
    status,
    statusText: ok ? "OK" : "Internal Server Error",
    json: () => Promise.resolve(body),
  };
}

beforeEach(() => {
  vi.spyOn(console, "error").mockImplementation(() => {});
});

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("API", () => {
  it("resolves with an axios-like { data } shape", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ groups: [] }));
    vi.stubGlobal("fetch", fetchMock);

    const res = await API.post("/api/v1/groups/view", { paging: {} });

    expect(res.data).toEqual({ groups: [] });
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe("/api/v1/groups/view");
    expect(init.method).toBe("POST");
    expect(init.headers["Content-Type"]).toBe("application/json");
    expect(JSON.parse(init.body)).toEqual({ paging: {} });
  });

  it("sends the method matching the wrapper", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({}));
    vi.stubGlobal("fetch", fetchMock);

    await API.patch("/api/v1/groups/rename", { groupId: 1, name: "x" });
    await API.delete("/api/v1/groups/delete", { groupId: 1 });

    expect(fetchMock.mock.calls[0][1].method).toBe("PATCH");
    expect(fetchMock.mock.calls[1][1].method).toBe("DELETE");
  });

  it("rejects with the server's error detail on non-2xx", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(jsonResponse({ error: "boom" }, false, 500))
    );

    await expect(
      API.post("/api/v1/groups/view", {})
    ).rejects.toThrow("boom");
  });

  it("rejects with the status line when the error body is not JSON", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 502,
        statusText: "Bad Gateway",
        json: () => Promise.reject(new SyntaxError("not json")),
      })
    );

    await expect(API.post("/api/v1/groups/view", {})).rejects.toThrow(
      "502 Bad Gateway"
    );
  });

  it("rejects on network error", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("Failed to fetch"))
    );

    await expect(API.post("/api/v1/groups/view", {})).rejects.toThrow(
      "Failed to fetch"
    );
  });
});
