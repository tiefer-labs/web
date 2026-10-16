// SPDX-License-Identifier: MPL-2.0
"use client";

import { useRef } from "react";
import { cn } from "@/lib/cn";
import { heroSvgClass } from "./globe";
import { type EventId, heroScene, ORBITS, type Stage } from "./hero-scene";
import { usePassProgress } from "./use-pass-progress";

export type EventLabels = Record<Stage, string> & { event: string; responder: string };
export type HeroOrbitLabels = {
  stored: string;
  sent: string;
  /** What has gone to the ground after one alert, two alerts, and so on. */
  sentSizes: string[];
  alertTag: string;
  ais: string;
  events: Record<EventId, EventLabels>;
};

const tele = "font-text text-[11px] leading-none font-medium tracking-[0.14em]";
const RETICLE = "M-14 -8V-14H-8M8 -14H14V-8M14 8V14H8M-8 14H-14V8";
// Drawn with its base on the origin, so it scales up from the ground.
const FLAME =
  "M0 -24C4 -17 10 -13 10 -6C10 0 5.5 3 0 3C-5.5 3 -10 0 -10 -6C-10 -10 -7 -13 -5.5 -16.5C-3.5 -13.5 -2 -12 0 -12C0 -16 -1 -20 0 -24Z";
const CORE =
  "M0 -12C2 -9 4.5 -7 4.5 -3.5C4.5 -0.5 2.5 1 0 1C-2.5 1 -4.5 -0.5 -4.5 -3.5C-4.5 -6 -2.5 -8 -1.5 -9.5C-1 -8 -0.5 -7.5 0 -7.5Z";
const SLICK = "M-24 -3C-21 -12 -6 -15 5 -12C15 -9 27 -11 30 -3C33 5 22 11 9 10C-1 9 -9 15 -18 11C-26 8 -27 2 -24 -3Z";
const SHEEN = "M-12 -3C-8 -8 2 -8 8 -6C14 -4 18 -1 14 3C9 6 0 4 -6 5C-12 5 -15 1 -12 -3Z";
// Pointing up: the bow is at the top.
const HULL = "M0 -8L3.5 -3V7H-3.5V-3Z";
const CLOUD: [number, number, number][] = [
  [-24, 4, 13],
  [-6, -6, 18],
  [14, -2, 15],
  [28, 6, 10],
  [2, 8, 12],
];
const BAR = 250;
const STAGES: Stage[] = ["scanning", "detected", "alerted"];
const EVENTS: EventId[] = ["wildfire", "oilSpill", "vessels"];
const RESPONDER_LABEL: Record<"left" | "below", { x: number; y: number; anchor: "start" | "end" }> = {
  left: { x: -12, y: 4, anchor: "end" },
  below: { x: -6, y: 20, anchor: "start" },
};
/** Where each event's name sits: left of its marker, or below the silent ship, clear of the other ships' tags. */
const EVENT_LABEL: Record<EventId, { x: number; y: number; anchor: "end" | "middle" }> = {
  wildfire: { x: -20, y: -4, anchor: "end" },
  oilSpill: { x: -40, y: -4, anchor: "end" },
  vessels: { x: 0, y: 30, anchor: "middle" },
};

/**
 * Three passes over Europe, on repeat. Each time the satellite images the ground and keeps every frame on board, finds
 * one event (a wildfire, an oil spill, a ship without AIS) and sends only a 48 KB alert down to the people who act on it.
 */
export function HeroOrbit({ labels }: { labels: HeroOrbitLabels }) {
  const svg = useRef<SVGSVGElement>(null);
  const s = heroScene(usePassProgress(svg));
  const alarm = s.stage !== "scanning";
  const side = s.responder.side;

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
      {ORBITS.map((d, index) => (
        <path
          // oxlint-disable-next-line react/no-array-index-key
          key={index}
          className="fill-none stroke-on-ink-muted"
          strokeWidth="1.25"
          strokeDasharray="2 7"
          vectorEffect="non-scaling-stroke"
          opacity={s.orbits[index]}
          d={d}
        />
      ))}

      {s.frames.map(({ lon, points, opacity }) => (
        <polygon
          key={`${s.event}-${lon}`}
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

      <g opacity={s.fire.opacity} transform={`translate(${s.fire.x} ${s.fire.y})`}>
        <circle className="fill-signal" r="15" opacity={s.fire.glow * 0.5} />
        <g transform={`scale(${s.fire.flame * 0.9})`}>
          <path className="fill-signal" d={FLAME} />
          <path className="fill-on-ink opacity-90" d={CORE} />
        </g>
      </g>

      <g opacity={s.spill.opacity}>
        <g transform={`translate(${s.spill.x} ${s.spill.y}) scale(${s.spill.scale})`}>
          <path className="fill-signal stroke-signal" fillOpacity="0.25" vectorEffect="non-scaling-stroke" d={SLICK} />
          <path className="fill-signal opacity-50" d={SHEEN} />
        </g>
        {/* The offshore platform the slick spreads from. */}
        <g
          className="fill-none stroke-on-ink"
          strokeWidth="1.25"
          transform={`translate(${s.spill.rig.x} ${s.spill.rig.y})`}
        >
          <rect vectorEffect="non-scaling-stroke" x="-5" y="-9" width="10" height="5" />
          <path vectorEffect="non-scaling-stroke" d="M-4 -4L-6 4M4 -4L6 4M0 -9V-14" />
        </g>
      </g>

      {s.ships.map((ship, index) => (
        // oxlint-disable-next-line react/no-array-index-key
        <g key={index} opacity={ship.opacity} transform={`translate(${ship.x} ${ship.y})`}>
          {!ship.dark && (
            <circle
              className="fill-none stroke-on-ink-muted"
              vectorEffect="non-scaling-stroke"
              r={ship.signal.r}
              opacity={ship.signal.opacity}
            />
          )}
          <path
            className={ship.dark ? "fill-on-ink-muted" : "fill-on-ink"}
            transform={`rotate(${ship.heading})`}
            d={HULL}
          />
          <rect
            className="fill-none stroke-on-ink-muted"
            vectorEffect="non-scaling-stroke"
            opacity={ship.matched}
            x="-11"
            y="-11"
            width="22"
            height="22"
          />
          <text className={cn(tele, "fill-on-ink-muted")} y="-16" textAnchor="middle" opacity={ship.matched}>
            {labels.ais}
          </text>
        </g>
      ))}

      {s.targets.map((target, index) => (
        // oxlint-disable-next-line react/no-array-index-key
        <g key={index} transform={`translate(${target.x} ${target.y})`}>
          <path
            className="fill-none stroke-on-ink"
            strokeWidth="1.5"
            vectorEffect="non-scaling-stroke"
            opacity={target.reticle.opacity}
            transform={`translate(0 ${s.event === "wildfire" ? -8 : 0}) scale(${target.reticle.scale})`}
            d={RETICLE}
          />
          {EVENTS.map((event) => (
            <text
              key={event}
              className={cn(tele, "fill-signal")}
              x={EVENT_LABEL[event].x}
              y={EVENT_LABEL[event].y}
              textAnchor={EVENT_LABEL[event].anchor}
              opacity={event === s.event ? target.detected : 0}
            >
              {labels.events[event].event}
            </text>
          ))}
        </g>
      ))}

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

      <g opacity={s.responder.marker} transform={`translate(${s.responder.x} ${s.responder.y})`}>
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
        {/* Away from the alert line, on the side it does not arrive from. */}
        {EVENTS.map((event) => (
          <text
            key={event}
            className={cn(tele, "fill-on-ink")}
            x={RESPONDER_LABEL[side].x}
            y={RESPONDER_LABEL[side].y}
            textAnchor={RESPONDER_LABEL[side].anchor}
            opacity={event === s.event ? s.responder.named * (0.7 + 0.3 * s.responder.lit) : 0}
          >
            {labels.events[event].responder}
          </text>
        ))}
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
        {EVENTS.flatMap((event) =>
          STAGES.map((stage) => (
            <text
              key={`${event}-${stage}`}
              className="fill-on-ink text-[12px] transition-opacity duration-300"
              x="18"
              opacity={event === s.event && stage === s.stage ? 1 : 0}
            >
              {labels.events[event][stage]}
            </text>
          )),
        )}
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
        <rect className="fill-signal" y="74" width={2 * s.store.sent} height="6" />
        {labels.sentSizes.map((size, index) => (
          <text
            key={size}
            className="fill-signal"
            x={BAR}
            y="66"
            textAnchor="end"
            opacity={s.store.sent === index + 1 ? 1 : 0}
          >
            {size}
          </text>
        ))}
      </g>
    </svg>
  );
}
