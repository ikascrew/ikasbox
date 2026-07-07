import { describe, it, expect } from "vitest";
import Paging from "./Paging.js";

describe("Paging.create", () => {
  it("defaults to a limit of 10", () => {
    const paging = Paging.create();
    expect(paging).toEqual({ current: 1, count: 0, limit: 10 });
  });

  it("accepts a custom limit", () => {
    const paging = Paging.create(100);
    expect(paging.limit).toBe(100);
    expect(paging.current).toBe(1);
    expect(paging.count).toBe(0);
  });
});
