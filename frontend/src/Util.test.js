import { describe, it, expect } from "vitest";
import Util from "./Util.js";

describe("Util.formatDate", () => {
  it("formats a date/time string as yyyy/MM/dd HH:mm:ss", () => {
    const formatted = Util.formatDate("2026-03-05T09:07:03");
    expect(formatted).toBe("2026/03/05 09:07:03");
  });

  it("zero-pads single-digit month/day/time components", () => {
    const formatted = Util.formatDate("2026-01-02T03:04:05");
    expect(formatted).toBe("2026/01/02 03:04:05");
  });
});
