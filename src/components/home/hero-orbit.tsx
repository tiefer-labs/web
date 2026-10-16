// SPDX-License-Identifier: MPL-2.0
"use client";

import { useRef } from "react";
import { cn } from "@/lib/cn";
import { heroSvgClass } from "./globe";
import { heroScene, type Stage } from "./hero-scene";
import { usePassProgress } from "./use-pass-progress";

export type HeroOrbitLabels = Record<Stage, string> & {
  fire: string;
  responder: string;
  stored: string;
  sent: string;
  alertSize: string;
  alertTag: string;
};

const tele = "font-text text-[11px] leading-none font-medium tracking-[0.14em]";
const RETICLE = "M-14 -8V-14H-8M8 -14H14V-8M14 8V14H8M-8 14H-14V8";
// Drawn with its base on the origin, so it scales up from the ground.
const FLAME =
  "M0 -24C4 -17 10 -13 10 -6C10 0 5.5 3 0 3C-5.5 3 -10 0 -10 -6C-10 -10 -7 -13 -5.5 -16.5C-3.5 -13.5 -2 -12 0 -12C0 -16 -1 -20 0 -24Z";
const CORE =
  "M0 -12C2 -9 4.5 -7 4.5 -3.5C4.5 -0.5 2.5 1 0 1C-2.5 1 -4.5 -0.5 -4.5 -3.5C-4.5 -6 -2.5 -8 -1.5 -9.5C-1 -8 -0.5 -7.5 0 -7.5Z";
const CLOUD: [number, number, number][] = [
  [-24, 4, 13],
  [-6, -6, 18],
  [14, -2, 15],
  [28, 6, 10],
  [2, 8, 12],
];
const BAR = 250;
const STAGES: Stage[] = ["scanning", "detected", "alerted"];

/**
 * One pass over Europe, on repeat: the satellite images the land under it and keeps every frame on board, finds a
 * wildfire in Greece and sends only a 48 KB alert down to the emergency services.
 */
export function HeroOrbit({ labels }: { labels: HeroOrbitLabels }) {
  const svg = useRef<SVGSVGElement>(null);
  const s = heroScene(usePassProgress(svg));
  const alarm = s.stage !== "scanning";

  return (
    <svg
      ref={svg}
      className={heroSvgClass}
      viewBox="0 0 1440 760"
      preserveAspectRatio="xMaxYMax slice"
      aria-hidden="true"
      focusable="false"
      data-print="hide"
    >
      {s.frames.map(({ lon, points, opacity }) => (
        <polygon
          key={lon}
          className="fill-on-ink stroke-on-ink"
          fillOpacity="0.06"
          vectorEffect="non-scaling-stroke"
          opacity={opacity}
          points={points}
        />
      ))}

      {s.clouds.map(({ x, y, size }, index) => (
        <g
          // oxlint-disable-next-line react/no-array-index-key
          key={index}
          className="fill-on-ink opacity-15"
          transform={`translate(${x} ${y}) scale(${size})`}
        >
          {CLOUD.map(([cx, cy, r]) => (
            <circle key={`${cx}-${cy}`} cx={cx} cy={cy} r={r} />
          ))}
        </g>
      ))}

      <g transform={`translate(${s.fire.x} ${s.fire.y})`}>
        <circle className="fill-signal" r="15" opacity={s.fire.glow * 0.5} />
        <g transform={`scale(${s.fire.flame * 0.9})`}>
          <path className="fill-signal" d={FLAME} />
          <path className="fill-on-ink opacity-90" d={CORE} />
        </g>
        <path
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          opacity={s.fire.reticle.opacity}
          transform={`translate(0 -8) scale(${s.fire.reticle.scale})`}
          d={RETICLE}
        />
        <text className={cn(tele, "fill-signal")} x="-20" y="-4" textAnchor="end" opacity={s.fire.detected}>
          {labels.fire}
        </text>
      </g>

      <polygon
        className="fill-on-ink-muted stroke-on-ink-muted"
        fillOpacity="0.1"
        vectorEffect="non-scaling-stroke"
        opacity={s.beam.opacity * 0.6}
        points={s.beam.points}
      />
      <polygon
        className="fill-none stroke-on-ink"
        vectorEffect="non-scaling-stroke"
        opacity={s.beam.opacity}
        points={s.beam.footprint}
      />

      <line
        className="fill-none stroke-signal"
        strokeWidth="1.5"
        strokeDasharray="5 5"
        strokeDashoffset={s.link.offset}
        vectorEffect="non-scaling-stroke"
        opacity={s.link.opacity}
        x1={s.link.from.x}
        y1={s.link.from.y}
        x2={s.link.to.x}
        y2={s.link.to.y}
      />
      {s.link.packet && (
        <rect className="fill-signal" x={s.link.packet.x - 4} y={s.link.packet.y - 4} width="8" height="8" />
      )}
      <text
        className={cn(tele, "fill-signal")}
        opacity={s.link.label.opacity}
        x={s.link.label.x}
        y={s.link.label.y}
        textAnchor="middle"
        transform={`rotate(${s.link.label.angle} ${s.link.label.x} ${s.link.label.y})`}
      >
        {labels.alertTag}
      </text>

      <g transform={`translate(${s.responder.x} ${s.responder.y})`}>
        {s.responder.pings.map(({ r, opacity }, index) => (
          <circle
            // oxlint-disable-next-line react/no-array-index-key
            key={index}
            className="fill-none stroke-signal"
            strokeWidth="1.25"
            vectorEffect="non-scaling-stroke"
            r={r}
            opacity={opacity}
          />
        ))}
        <rect className="fill-on-ink-muted" x="-5" y="-5" width="10" height="10" />
        <rect className="fill-signal" x="-5" y="-5" width="10" height="10" opacity={s.responder.lit} />
        {/* Left of the marker: the alert line arrives from the east. */}
        <text
          className={cn(tele, "fill-on-ink")}
          x="-12"
          y="4"
          textAnchor="end"
          opacity={s.responder.named * (0.7 + 0.3 * s.responder.lit)}
        >
          {labels.responder}
        </text>
      </g>

      <g opacity={s.sat.opacity} transform={`translate(${s.sat.x} ${s.sat.y}) rotate(${s.sat.angle})`}>
        <rect
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="-4.5"
          y="-31"
          width="9"
          height="22"
        />
        <rect className="fill-on-ink" x="-7.5" y="-7" width="15" height="14" />
        <rect
          className="fill-none stroke-on-ink"
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
          x="-4.5"
          y="9"
          width="9"
          height="22"
        />
      </g>

      {/* Below the header on wide screens, where the top of the box is cropped. */}
      <g transform="translate(1120 205)" className={tele}>
        <rect
          className="fill-ink stroke-line-on-ink"
          fillOpacity="0.85"
          vectorEffect="non-scaling-stroke"
          x="-18"
          y="-28"
          width="286"
          height="124"
        />
        <rect className={alarm ? "fill-signal" : "fill-on-ink"} y="-9" width="9" height="9" />
        {STAGES.map((stage) => (
          <text
            key={stage}
            className="fill-on-ink text-[12px] transition-opacity duration-300"
            x="18"
            opacity={stage === s.stage ? 1 : 0}
          >
            {labels[stage]}
          </text>
        ))}
        <text className="fill-on-ink-muted" y="30">
          {labels.stored}
        </text>
        <rect
          className="fill-none stroke-line-on-ink"
          vectorEffect="non-scaling-stroke"
          y="38"
          width={BAR}
          height="6"
        />
        <rect className="fill-on-ink-muted" y="38" width={(BAR * s.store.kept) / s.store.of} height="6" />
        <text className="fill-on-ink-muted" y="66">
          {labels.sent}
        </text>
        <rect
          className="fill-none stroke-line-on-ink"
          vectorEffect="non-scaling-stroke"
          y="74"
          width={BAR}
          height="6"
        />
        <rect className="fill-signal" y="74" width={s.store.sent ? 2 : 0} height="6" />
        <text className="fill-signal" x={BAR} y="66" textAnchor="end" opacity={s.store.sent ? 1 : 0}>
          {labels.alertSize}
        </text>
      </g>
    </svg>
  );
}
