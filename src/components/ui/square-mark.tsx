// SPDX-License-Identifier: MPL-2.0
import { cn } from "@/lib/cn";

/** The 12 px square node used on the flow and roadmap tracks. */
export function SquareMark({ className, rectClassName }: { className?: string; rectClassName?: string }) {
  return (
    <svg
      className={cn("absolute", className)}
      viewBox="0 0 12 12"
      width="12"
      height="12"
      aria-hidden="true"
      focusable="false"
    >
      <rect x="1" y="1" width="10" height="10" strokeWidth="1.25" className={rectClassName} />
    </svg>
  );
}
