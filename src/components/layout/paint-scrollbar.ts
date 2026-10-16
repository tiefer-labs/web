// SPDX-License-Identifier: MPL-2.0

/** Sizes and places the thumb. Self-contained, because PageScrollbar also inlines it as a script. */
export function paintScrollbar(track: HTMLElement, thumb: HTMLElement) {
  const root = document.documentElement;
  const view = root.clientHeight;
  const range = root.scrollHeight - view;
  const size = range > 0 ? Math.max(32, (view / root.scrollHeight) * view) : 0;
  track.toggleAttribute("data-empty", range <= 0);
  thumb.style.height = `${size}px`;
  thumb.style.transform = `translateY(${range > 0 ? (root.scrollTop / range) * (view - size) : 0}px)`;
  return { range, travel: view - size };
}
