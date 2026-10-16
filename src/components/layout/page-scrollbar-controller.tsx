// SPDX-License-Identifier: MPL-2.0
"use client";

import { useEffect } from "react";
import { paintScrollbar } from "./paint-scrollbar";

/** Keeps the thumb from PageScrollbar in step with the page, and makes it draggable. */
export function PageScrollbarController({ trackId, thumbId }: { trackId: string; thumbId: string }) {
  useEffect(() => {
    const track = document.getElementById(trackId);
    const thumb = document.getElementById(thumbId);
    if (!track || !thumb) return undefined;
    const root = document.documentElement;
    let drag: { y: number; top: number } | null = null;

    const update = () => paintScrollbar(track, thumb);

    // "instant" so the thumb follows the pointer instead of the smooth scrolling set on html.
    const onPointerMove = (event: PointerEvent) => {
      if (!drag) return;
      const { range, travel } = update();
      window.scrollTo({ top: drag.top + ((event.clientY - drag.y) * range) / travel, behavior: "instant" });
    };

    const onPointerUp = () => {
      drag = null;
      delete track.dataset.dragging;
    };

    const onThumbDown = (event: PointerEvent) => {
      if (event.button !== 0) return;
      event.preventDefault();
      event.stopPropagation();
      drag = { y: event.clientY, top: root.scrollTop };
      track.dataset.dragging = "";
    };

    // Like a native track: a click pages towards the pointer.
    const onTrackDown = (event: PointerEvent) => {
      if (event.button !== 0) return;
      const { top } = thumb.getBoundingClientRect();
      window.scrollBy({ top: Math.sign(event.clientY - top) * root.clientHeight * 0.9 });
    };

    update();
    const resize = typeof ResizeObserver === "undefined" ? null : new ResizeObserver(update);
    resize?.observe(document.body);
    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);
    window.addEventListener("pointermove", onPointerMove);
    window.addEventListener("pointerup", onPointerUp);
    thumb.addEventListener("pointerdown", onThumbDown);
    track.addEventListener("pointerdown", onTrackDown);
    return () => {
      resize?.disconnect();
      window.removeEventListener("scroll", update);
      window.removeEventListener("resize", update);
      window.removeEventListener("pointermove", onPointerMove);
      window.removeEventListener("pointerup", onPointerUp);
      thumb.removeEventListener("pointerdown", onThumbDown);
      track.removeEventListener("pointerdown", onTrackDown);
    };
  }, [trackId, thumbId]);

  return null;
}
