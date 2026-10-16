// SPDX-License-Identifier: MPL-2.0
import { act } from "@testing-library/react";
import { afterEach, expect, test, vi } from "vitest";
import { HeroOrbit } from "@/components/home/hero-orbit";
import { heroScene, POSTER } from "@/components/home/hero-scene";
import { renderClient } from "../../test-utils/render";

function preferReducedMotion(reduce: boolean) {
  vi.stubGlobal("matchMedia", (query: string) => ({
    matches: reduce && query.includes("reduce"),
    addEventListener: vi.fn<() => void>(),
    removeEventListener: vi.fn<() => void>(),
  }));
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

const satellite = (container: HTMLElement) =>
  container.querySelector('g[opacity] > g[transform^="translate"]')?.getAttribute("transform");

test("stays on the still frame for visitors who prefer reduced motion", () => {
  preferReducedMotion(true);
  const frame = vi.spyOn(window, "requestAnimationFrame");
  const { container } = renderClient(<HeroOrbit stored="KEPT ON BOARD" alertTag="ALERT 48 KB" />);
  expect(frame).not.toHaveBeenCalled();
  const { sat } = heroScene(POSTER);
  expect(satellite(container)).toBe(`translate(${sat.x} ${sat.y}) rotate(${sat.angle})`);
  expect(container.querySelector("svg")).toHaveAttribute("aria-hidden", "true");
  expect(container).toHaveTextContent("ALERT 48 KB");
});

test("starts at the beginning of the pass and moves the satellite frame by frame", () => {
  preferReducedMotion(false);
  const callbacks: FrameRequestCallback[] = [];
  vi.spyOn(window, "requestAnimationFrame").mockImplementation((callback) => callbacks.push(callback));
  const { container } = renderClient(<HeroOrbit stored="KEPT ON BOARD" alertTag="ALERT 48 KB" />);
  const start = satellite(container);
  const { sat } = heroScene(0);
  expect(start).toBe(`translate(${sat.x} ${sat.y}) rotate(${sat.angle})`);

  act(() => callbacks.shift()?.(0));
  act(() => callbacks.shift()?.(1000));
  expect(satellite(container)).not.toBe(start);
});
