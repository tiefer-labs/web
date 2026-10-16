// SPDX-License-Identifier: MPL-2.0
import { fireEvent } from "@testing-library/react";
import { afterEach, beforeEach, expect, test, vi } from "vitest";
import { PageScrollbar } from "@/components/layout/page-scrollbar";
import { renderClient, renderServer } from "../../test-utils/render";

const root = document.documentElement;
const track = () => document.getElementById("page-scrollbar");
const thumb = () => document.getElementById("page-scrollbar-thumb");

beforeEach(() => {
  vi.spyOn(root, "clientHeight", "get").mockReturnValue(800);
  vi.spyOn(root, "scrollHeight", "get").mockReturnValue(3200);
  root.scrollTop = 0;
});

afterEach(() => {
  vi.restoreAllMocks();
  document.body.replaceChildren();
});

test("the inline script paints the thumb before hydration", async () => {
  const { container } = await renderServer(<PageScrollbar />);
  root.scrollTop = 1200;
  // oxlint-disable-next-line typescript/no-implied-eval
  new Function(container.querySelector("script")?.textContent ?? "")();
  expect(track()).toHaveAttribute("aria-hidden", "true");
  expect(track()).not.toHaveAttribute("data-empty");
  expect(thumb()?.style.height).toBe("200px");
  expect(thumb()?.style.transform).toBe("translateY(300px)");
});

test("follows scrolling once hydrated", () => {
  renderClient(<PageScrollbar />);
  root.scrollTop = 1200;
  fireEvent.scroll(window);
  expect(thumb()?.style.transform).toBe("translateY(300px)");
});

test("hides itself when the page does not scroll", () => {
  vi.spyOn(root, "scrollHeight", "get").mockReturnValue(800);
  renderClient(<PageScrollbar />);
  expect(track()).toHaveAttribute("data-empty");
});

test("dragging the thumb scrolls the page in proportion", () => {
  const scrollTo = vi.spyOn(window, "scrollTo").mockImplementation(() => {});
  renderClient(<PageScrollbar />);
  const element = thumb();
  if (!element) throw new Error("no thumb");

  fireEvent.pointerDown(element, { button: 0, clientY: 100 });
  fireEvent.pointerMove(window, { clientY: 160 });
  fireEvent.pointerUp(window);

  // 60px of thumb travel over 600px of track is a tenth of the 2400px range.
  expect(scrollTo).toHaveBeenCalledWith({ top: 240, behavior: "instant" });
});
