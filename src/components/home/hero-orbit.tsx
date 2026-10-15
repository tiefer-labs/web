// SPDX-License-Identifier: MPL-2.0
import { cn } from "@/lib/cn";

const stars: [number, number, number][] = [
  [612, 64, 1],
  [742, 148, 1.2],
  [826, 42, 0.8],
  [918, 118, 1],
  [1004, 58, 1.2],
  [1188, 96, 0.8],
  [1262, 36, 1],
  [1342, 132, 1.2],
  [1402, 74, 0.8],
  [1310, 262, 1],
  [1396, 318, 0.8],
  [1226, 352, 1],
  [1122, 402, 0.8],
  [1374, 436, 1],
  [866, 292, 0.8],
  [958, 372, 1],
  [700, 330, 0.8],
  [782, 452, 1],
  [1032, 486, 0.8],
  [1290, 520, 1],
  [668, 226, 1],
];

const line = "fill-none stroke-on-ink-muted";
const tele = "fill-on-ink-muted font-text text-[11px] leading-none font-medium tracking-[0.14em] opacity-85";

/** A satellite passing over the Earth's limb, sending a small alert down to a ground station. */
export function HeroOrbit({ stored, alertTag }: { stored: string; alertTag: string }) {
  return (
    <svg
      className="absolute inset-0 size-full max-lg:top-auto max-lg:h-85"
      viewBox="0 0 1440 760"
      preserveAspectRatio="xMaxYMax slice"
      aria-hidden="true"
      focusable="false"
      data-print="hide"
    >
      <g className="fill-on-ink-muted opacity-40">
        {stars.map(([cx, cy, r]) => (
          <circle key={`${cx}-${cy}`} cx={cx} cy={cy} r={r} />
        ))}
      </g>
      <g strokeWidth="1.25">
        <path
          className={cn(line, "opacity-45")}
          strokeDasharray="2 7"
          vectorEffect="non-scaling-stroke"
          d="M520 266C760 162 1120 156 1500 276"
        />
        <path
          className={cn(line, "opacity-20")}
          strokeDasharray="2 7"
          vectorEffect="non-scaling-stroke"
          d="M600 392C860 274 1180 274 1500 378"
        />
        <path className="fill-on-ink-muted opacity-5" d="M-40 760C300 628 1140 628 1480 760V800H-40Z" />
        <path
          className={cn(line, "opacity-45")}
          vectorEffect="non-scaling-stroke"
          d="M-40 760C300 628 1140 628 1480 760"
        />
        <path
          className={cn(line, "opacity-15")}
          vectorEffect="non-scaling-stroke"
          d="M-40 744C300 612 1140 612 1480 744"
        />
        <path
          className="fill-on-ink-muted stroke-on-ink-muted opacity-30"
          fillOpacity="0.07"
          vectorEffect="non-scaling-stroke"
          d="M1033 198L1116 682L1176 690Z"
        />
        <path
          className={cn(line, "opacity-45")}
          strokeDasharray="4 5"
          vectorEffect="non-scaling-stroke"
          d="M1029 199L900 661"
        />
      </g>
      <rect className="fill-on-ink-muted opacity-70" x="896" y="661" width="8" height="8" />
      <g className="motion-safe:animate-satellite-pass">
        <rect
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="1004"
          y="185"
          width="22"
          height="9"
        />
        <rect className="fill-on-ink" x="1028" y="182" width="14" height="15" />
        <rect
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="1044"
          y="185"
          width="22"
          height="9"
        />
      </g>
      <text className={tele} x="1078" y="194">
        {stored}
      </text>
      <text className={tele} x="921.2" y="566.6" transform="rotate(-74.4 921.2 566.6)">
        {alertTag}
      </text>
    </svg>
  );
}
