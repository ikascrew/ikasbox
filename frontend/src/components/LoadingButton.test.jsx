import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import LoadingButton from "./LoadingButton.jsx";

describe("LoadingButton", () => {
  it("calls onClick and disables itself while the returned promise is pending", async () => {
    let resolvePromise;
    const onClick = vi.fn(() => new Promise((resolve) => { resolvePromise = resolve; }));

    render(<LoadingButton onClick={onClick}>Save</LoadingButton>);

    const button = screen.getByRole("button", { name: "Save" });
    fireEvent.click(button);

    expect(onClick).toHaveBeenCalledTimes(1);
    expect(button).toBeDisabled();

    resolvePromise();
    await waitFor(() => expect(button).not.toBeDisabled());
  });
});
