// SPDX-License-Identifier: MPL-2.0
import { PageScrollbarController } from "./page-scrollbar-controller";
import { paintScrollbar } from "./paint-scrollbar";

const TRACK_ID = "page-scrollbar";
const THUMB_ID = "page-scrollbar-thumb";

// Runs as soon as the markup is parsed, so the thumb is there on the first paint instead of after hydration.
const paintNow = `(${paintScrollbar.toString()})(document.getElementById("${TRACK_ID}"),document.getElementById("${THUMB_ID}"))`;

/**
 * Draws the page scrollbar over the content, so dark sections reach the window edge. Mouse pointers only; hidden
 * from assistive technology because keyboard and wheel scrolling stay native.
 */
export function PageScrollbar() {
  return (
    <>
      {/* The inline script changes these attributes before React hydrates them. */}
      <div
        id={TRACK_ID}
        aria-hidden="true"
        data-print="hide"
        data-empty=""
        suppressHydrationWarning
        className="group fixed inset-y-0 right-0 z-20 hidden w-3 pointer-fine:not-data-empty:block"
      >
        <div
          id={THUMB_ID}
          suppressHydrationWarning
          className="absolute top-0 right-0.5 w-1.5 rounded-sm bg-ink-60 transition-[width] duration-150 group-hover:w-2 group-data-dragging:w-2 forced-colors:bg-[CanvasText]"
        />
      </div>
      <script dangerouslySetInnerHTML={{ __html: paintNow }} />
      <PageScrollbarController trackId={TRACK_ID} thumbId={THUMB_ID} />
    </>
  );
}
