// SPDX-License-Identifier: MPL-2.0
"use client";

import { useEffect, useRef } from "react";

const MIN_THUMB = 32;

/** Draws the page scrollbar over the content, so dark sections reach the window edge. Mouse pointers only. */
export function PageScrollbar() {
  const track = useRef<HTMLDivElement>(null);
  const thumb = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const trackEl = track.current;
    const thumbEl = thumb.current;
    if (!trackEl || !thumbEl) return undefined;
    const root = document.documentElement;
    let drag: { y: number; top: number } | null = null;

    const geometry = () => {
      const view = root.clientHeight;
      const range = root.scrollHeight - view;
      const size = range > 0 ? Math.max(MIN_THUMB, (view / root.scrollHeight) * view) : 0;
      return { view, range, size, travel: view - size };
    };

    const update = () => {
      const { range, size, travel } = geometry();
      trackEl.toggleAttribute("data-empty", range <= 0);
      thumbEl.style.height = `${size}px`;
      thumbEl.style.transform = `translateY(${range > 0 ? (root.scrollTop / range) * travel : 0}px)`;
    };

    // "instant" so the thumb follows the pointer instead of the smooth scrolling set on html.
    const onPointerMove = (event: PointerEvent) => {
      if (!drag) return;
      const { range, travel } = geometry();
      window.scrollTo({ top: drag.top + ((event.clientY - drag.y) * range) / travel, behavior: "instant" });
    };

    const onPointerUp = () => {
      drag = null;
      delete trackEl.dataset.dragging;
    };

    const onThumbDown = (event: PointerEvent) => {
      if (event.button !== 0) return;
      event.preventDefault();
      event.stopPropagation();
      drag = { y: event.clientY, top: root.scrollTop };
      trackEl.dataset.dragging = "";
    };

    // Like a native track: a click pages towards the pointer.
    const onTrackDown = (event: PointerEvent) => {
      if (event.button !== 0) return;
      const { top } = thumbEl.getBoundingClientRect();
      window.scrollBy({ top: Math.sign(event.clientY - top) * root.clientHeight * 0.9 });
    };

    update();
    const resize = typeof ResizeObserver === "undefined" ? null : new ResizeObserver(update);
    resize?.observe(document.body);
    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);
    window.addEventListener("pointermove", onPointerMove);
    window.addEventListener("pointerup", onPointerUp);
    thumbEl.addEventListener("pointerdown", onThumbDown);
    trackEl.addEventListener("pointerdown", onTrackDown);
    return () => {
      resize?.disconnect();
      window.removeEventListener("scroll", update);
      window.removeEventListener("resize", update);
      window.removeEventListener("pointermove", onPointerMove);
      window.removeEventListener("pointerup", onPointerUp);
      thumbEl.removeEventListener("pointerdown", onThumbDown);
      trackEl.removeEventListener("pointerdown", onTrackDown);
    };
  }, []);

  // Hidden from assistive technology: keyboard and wheel scrolling stay native.
  return (
    <div
      ref={track}
      aria-hidden="true"
      data-print="hide"
      data-empty=""
      className="group fixed inset-y-0 right-0 z-20 hidden w-3 pointer-fine:not-data-empty:block"
    >
      <div
        ref={thumb}
        data-testid="page-scrollbar-thumb"
        className="absolute top-0 right-0.5 w-1.5 rounded-sm bg-ink-60 transition-[width] duration-150 group-hover:w-2 group-data-dragging:w-2 forced-colors:bg-[CanvasText]"
      />
    </div>
  );
}
