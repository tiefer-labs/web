// SPDX-License-Identifier: MPL-2.0
import { fireEvent, screen } from "@testing-library/react";
import { afterEach, beforeEach, expect, test, vi } from "vitest";
import { PageScrollbar } from "@/components/layout/page-scrollbar";
import { renderClient } from "../../test-utils/render";

const root = document.documentElement;

beforeEach(() => {
  vi.spyOn(root, "clientHeight", "get").mockReturnValue(800);
  vi.spyOn(root, "scrollHeight", "get").mockReturnValue(3200);
  root.scrollTop = 0;
});

afterEach(() => {
  vi.restoreAllMocks();
});

test("sizes the thumb to the visible share of the page and follows scrolling", () => {
  const { container } = renderClient(<PageScrollbar />);
  const thumb = screen.getByTestId("page-scrollbar-thumb");
  expect(container.firstElementChild).toHaveAttribute("aria-hidden", "true");
  expect(thumb.style.height).toBe("200px");

  root.scrollTop = 1200;
  fireEvent.scroll(window);
  expect(thumb.style.transform).toBe("translateY(300px)");
});

test("hides itself when the page does not scroll", () => {
  vi.spyOn(root, "scrollHeight", "get").mockReturnValue(800);
  const { container } = renderClient(<PageScrollbar />);
  expect(container.firstElementChild).toHaveAttribute("data-empty");
});

test("dragging the thumb scrolls the page in proportion", () => {
  const scrollTo = vi.spyOn(window, "scrollTo").mockImplementation(() => {});
  renderClient(<PageScrollbar />);
  const thumb = screen.getByTestId("page-scrollbar-thumb");

  fireEvent.pointerDown(thumb, { button: 0, clientY: 100 });
  fireEvent.pointerMove(window, { clientY: 160 });
  fireEvent.pointerUp(window);

  // 60px of thumb travel over 600px of track is a tenth of the 2400px range.
  expect(scrollTo).toHaveBeenCalledWith({ top: 240, behavior: "instant" });
});
