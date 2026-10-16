// SPDX-License-Identifier: MPL-2.0
"use client";

import { type RefObject, useEffect, useState } from "react";
import { PERIOD, POSTER } from "./hero-scene";

/**
 * Plays the pass from its start, in a loop while the hero is on screen. Visitors who prefer reduced motion get the
 * still frame instead.
 */
export function usePassProgress(target: RefObject<Element | null>) {
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    const element = target.current;
    const reduce = window.matchMedia("(prefers-reduced-motion: reduce)");
    let frame = 0;
    let last: number | null = null;
    let current = 0;
    let onScreen = true;

    const tick = (now: number) => {
      if (last !== null) current = (current + (now - last) / PERIOD) % 1;
      last = now;
      setProgress(current);
      frame = requestAnimationFrame(tick);
    };
    const sync = () => {
      cancelAnimationFrame(frame);
      frame = 0;
      last = null;
      if (reduce.matches) {
        current = POSTER;
        setProgress(POSTER);
      } else if (onScreen) frame = requestAnimationFrame(tick);
    };

    const observer =
      element && typeof IntersectionObserver !== "undefined"
        ? new IntersectionObserver(([entry]) => {
            onScreen = entry?.isIntersecting ?? true;
            sync();
          })
        : null;
    if (element) observer?.observe(element);
    reduce.addEventListener("change", sync);
    sync();
    return () => {
      cancelAnimationFrame(frame);
      observer?.disconnect();
      reduce.removeEventListener("change", sync);
    };
  }, [target]);

  return progress;
}
