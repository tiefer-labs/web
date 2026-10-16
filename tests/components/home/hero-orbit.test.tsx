// SPDX-License-Identifier: MPL-2.0
import { act } from "@testing-library/react";
import { afterEach, expect, test, vi } from "vitest";
import { HeroOrbit, type HeroOrbitLabels } from "@/components/home/hero-orbit";
import { PASS } from "@/components/home/hero-scene";
import messages from "../../../messages/en.json";
import { renderClient } from "../../test-utils/render";

const { sentSize: _, ...hero } = messages.Hero;
const labels: HeroOrbitLabels = { ...hero, sentSizes: ["48 KB", "96 KB"] };
const { wildfire, oilSpill } = labels.events;

function preferReducedMotion(reduce: boolean) {
  vi.stubGlobal("matchMedia", (query: string) => ({
    matches: reduce && query.includes("reduce"),
    addEventListener: vi.fn<() => void>(),
    removeEventListener: vi.fn<() => void>(),
  }));
}

/** The status line that is showing: the others stay in the page, faded out. */
const status = (container: HTMLElement) =>
  [...container.querySelectorAll("text")].find(
    (text) => text.getAttribute("opacity") === "1" && text.classList.contains("text-[12px]"),
  )?.textContent;

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

test("visitors who prefer reduced motion see the still frame: the wildfire alert on its way down", () => {
  preferReducedMotion(true);
  const frame = vi.spyOn(window, "requestAnimationFrame");
  const { container } = renderClient(<HeroOrbit labels={labels} />);
  expect(frame).not.toHaveBeenCalled();
  expect(container.querySelector("svg")).toHaveAttribute("aria-hidden", "true");
  expect(status(container)).toBe(wildfire.detected);
});

test("starts with the wildfire pass and moves on to the oil spill", () => {
  preferReducedMotion(false);
  const callbacks: FrameRequestCallback[] = [];
  vi.spyOn(window, "requestAnimationFrame").mockImplementation((callback) => callbacks.push(callback));
  const { container } = renderClient(<HeroOrbit labels={labels} />);
  expect(status(container)).toBe(wildfire.scanning);

  act(() => callbacks.shift()?.(0));
  act(() => callbacks.shift()?.(PASS * 0.85));
  expect(status(container)).toBe(wildfire.alerted);
  act(() => callbacks.shift()?.(PASS * 1.3));
  expect(status(container)).toBe(oilSpill.scanning);
});
