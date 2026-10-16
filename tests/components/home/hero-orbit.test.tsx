// SPDX-License-Identifier: MPL-2.0
import { act } from "@testing-library/react";
import { afterEach, expect, test, vi } from "vitest";
import { HeroOrbit, type HeroOrbitLabels } from "@/components/home/hero-orbit";
import messages from "../../../messages/en.json";
import { renderClient } from "../../test-utils/render";

const labels: HeroOrbitLabels = messages.Hero;

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

test("visitors who prefer reduced motion see the still frame: the alert on its way down", () => {
  preferReducedMotion(true);
  const frame = vi.spyOn(window, "requestAnimationFrame");
  const { container } = renderClient(<HeroOrbit labels={labels} />);
  expect(frame).not.toHaveBeenCalled();
  expect(container.querySelector("svg")).toHaveAttribute("aria-hidden", "true");
  expect(container).toHaveTextContent(labels.detected);
  expect(container).toHaveTextContent(labels.alertTag);
});

test("starts at the beginning of the pass and plays it frame by frame", () => {
  preferReducedMotion(false);
  const callbacks: FrameRequestCallback[] = [];
  vi.spyOn(window, "requestAnimationFrame").mockImplementation((callback) => callbacks.push(callback));
  const { container } = renderClient(<HeroOrbit labels={labels} />);
  expect(container).toHaveTextContent(labels.scanning);

  act(() => callbacks.shift()?.(0));
  act(() => callbacks.shift()?.(13_500));
  expect(container).toHaveTextContent(labels.alerted);
});
