// SPDX-License-Identifier: MPL-2.0
"use client";

import { useEffect, useRef, useState } from "react";
import { cn } from "@/lib/cn";
import { CENTRE, EARTH, heroScene, INNER_ORBIT, ORBIT, PERIOD, POSTER } from "./hero-scene";

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
const tele = "fill-on-ink-muted font-text text-[11px] leading-none font-medium tracking-[0.14em]";
const RETICLE = "M-10 -5V-10H-5M5 -10H10V-5M10 5V10H5M-5 10H-10V5";

/**
 * Plays the pass from its start, in a loop while the hero is on screen. Visitors who prefer reduced motion get the
 * still frame instead.
 */
function usePassProgress(target: React.RefObject<Element | null>) {
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

/**
 * One pass, on repeat: the satellite keeps the frames it scans on board, locks onto an event on the ground and sends
 * a small alert down to the ground station.
 */
export function HeroOrbit({ stored, alertTag }: { stored: string; alertTag: string }) {
  const svg = useRef<SVGSVGElement>(null);
  const s = heroScene(usePassProgress(svg));

  return (
    <svg
      ref={svg}
      className="absolute inset-0 size-full max-lg:top-auto max-lg:h-85"
      viewBox="0 0 1440 760"
      preserveAspectRatio="xMaxYMax slice"
      aria-hidden="true"
      focusable="false"
      data-print="hide"
    >
      <g className="fill-on-ink-muted">
        {stars.map(([cx, cy, r], index) => (
          <circle key={`${cx}-${cy}`} cx={cx} cy={cy} r={r} opacity={s.stars(index)} />
        ))}
      </g>

      <defs>
        <linearGradient id="hero-orbit-fade" gradientUnits="userSpaceOnUse" x1="560" x2="900" y1="0" y2="0">
          <stop offset="0" stopColor="white" stopOpacity="0" />
          <stop offset="1" stopColor="white" />
        </linearGradient>
        {/* Fades the orbits out behind the hero copy. */}
        <mask id="hero-orbit-mask" maskUnits="userSpaceOnUse" x="0" y="0" width="1440" height="760">
          <rect width="1440" height="760" fill="url(#hero-orbit-fade)" />
        </mask>
      </defs>

      <g strokeWidth="1.25">
        <g mask="url(#hero-orbit-mask)">
          <circle
            className={cn(line, "opacity-45")}
            strokeDasharray="2 7"
            vectorEffect="non-scaling-stroke"
            cx={CENTRE.x}
            cy={CENTRE.y}
            r={ORBIT}
          />
          <circle
            className={cn(line, "opacity-20")}
            strokeDasharray="2 7"
            vectorEffect="non-scaling-stroke"
            cx={CENTRE.x}
            cy={CENTRE.y}
            r={INNER_ORBIT}
          />
        </g>
        <circle className="fill-on-ink-muted opacity-5" cx={CENTRE.x} cy={CENTRE.y} r={EARTH} />
        <circle
          className={cn(line, "opacity-15")}
          vectorEffect="non-scaling-stroke"
          cx={CENTRE.x}
          cy={CENTRE.y}
          r={EARTH + 16}
        />
        <circle
          className={cn(line, "opacity-45")}
          vectorEffect="non-scaling-stroke"
          cx={CENTRE.x}
          cy={CENTRE.y}
          r={EARTH}
        />
        <circle
          className={cn(line, "opacity-15")}
          strokeDasharray="1 9"
          vectorEffect="non-scaling-stroke"
          cx={CENTRE.x}
          cy={CENTRE.y}
          r={EARTH - 34}
        />

        {s.trail.map(({ from, to, opacity }, index) => (
          <path
            // oxlint-disable-next-line react/no-array-index-key
            key={index}
            className="fill-none stroke-on-ink"
            vectorEffect="non-scaling-stroke"
            opacity={opacity}
            d={`M${from.x} ${from.y}A${ORBIT} ${ORBIT} 0 0 1 ${to.x} ${to.y}`}
          />
        ))}

        <polygon
          className="fill-on-ink-muted stroke-on-ink-muted"
          fillOpacity="0.07"
          vectorEffect="non-scaling-stroke"
          opacity={s.swath.opacity}
          points={s.swath.points}
        />

        {s.tiles.map(({ ground, angle, opacity }) => (
          <rect
            key={angle}
            className="fill-none stroke-on-ink"
            vectorEffect="non-scaling-stroke"
            opacity={opacity}
            x="-3.5"
            y="-3.5"
            width="7"
            height="7"
            transform={`translate(${ground.x} ${ground.y}) rotate(${angle})`}
          />
        ))}
        {s.tiles.map(
          ({ frame, angle }) =>
            frame && (
              <rect
                key={`frame-${angle}`}
                className="fill-on-ink opacity-80"
                x={frame.x - 2}
                y={frame.y - 2}
                width="4"
                height="4"
              />
            ),
        )}

        <g transform={`translate(${s.target.x} ${s.target.y}) rotate(${s.target.angle})`}>
          <circle
            className={cn(line, "stroke-on-ink")}
            vectorEffect="non-scaling-stroke"
            r={s.target.pulse.r}
            opacity={s.target.pulse.opacity}
          />
          <circle className="fill-on-ink-muted" r="2.5" />
          <circle className="fill-on-ink" r="2.5" opacity={s.target.detected} />
          <path
            className="fill-none stroke-on-ink"
            strokeWidth="1.5"
            vectorEffect="non-scaling-stroke"
            opacity={s.target.reticle.opacity}
            transform={`scale(${s.target.reticle.scale})`}
            d={RETICLE}
          />
        </g>

        <line
          className={line}
          strokeDasharray="4 5"
          strokeDashoffset={s.link.offset}
          vectorEffect="non-scaling-stroke"
          opacity={s.link.opacity}
          x1={s.link.from.x}
          y1={s.link.from.y}
          x2={s.link.to.x}
          y2={s.link.to.y}
        />
        {s.link.packet && (
          <rect className="fill-on-ink" x={s.link.packet.x - 3} y={s.link.packet.y - 3} width="6" height="6" />
        )}

        <g transform={`translate(${s.station.x} ${s.station.y}) rotate(${s.station.angle})`}>
          {s.station.pings.map(({ r, opacity }, index) => (
            <circle
              // oxlint-disable-next-line react/no-array-index-key
              key={index}
              className={line}
              vectorEffect="non-scaling-stroke"
              r={r}
              opacity={opacity}
            />
          ))}
          <rect className="fill-on-ink-muted opacity-70" x="-4" y="-8" width="8" height="8" />
          <rect className="fill-on-ink" x="-4" y="-8" width="8" height="8" opacity={s.station.lit} />
        </g>
      </g>

      <rect
        className="fill-on-ink-muted"
        opacity={s.companion.opacity}
        x={s.companion.x - 1.5}
        y={s.companion.y - 1.5}
        width="3"
        height="3"
      />

      <g opacity={s.sat.opacity}>
        <g transform={`translate(${s.sat.x} ${s.sat.y}) rotate(${s.sat.angle})`}>
          <rect
            className="fill-none stroke-on-ink"
            strokeWidth="1.5"
            vectorEffect="non-scaling-stroke"
            x="-31"
            y="-4.5"
            width="22"
            height="9"
          />
          <rect className="fill-on-ink" x="-7" y="-7.5" width="14" height="15" />
          <rect
            className="fill-none stroke-on-ink"
            strokeWidth="1.5"
            vectorEffect="non-scaling-stroke"
            x="9"
            y="-4.5"
            width="22"
            height="9"
          />
        </g>
        <text className={cn(tele, "opacity-85")} x={s.sat.x + 44} y={s.sat.y + 4}>
          {stored}
        </text>
        {[0, 1, 2, 3].map((slot) => (
          <rect
            key={slot}
            className={slot < s.kept ? "fill-on-ink" : "fill-none stroke-on-ink-muted opacity-50"}
            vectorEffect="non-scaling-stroke"
            x={s.sat.x + 44.5 + slot * 8}
            y={s.sat.y + 12.5}
            width="5"
            height="5"
          />
        ))}
      </g>

      <text
        className={tele}
        opacity={s.link.label.opacity}
        x={s.link.label.x}
        y={s.link.label.y}
        textAnchor="middle"
        transform={`rotate(${s.link.label.angle} ${s.link.label.x} ${s.link.label.y})`}
      >
        {alertTag}
      </text>
    </svg>
  );
}
