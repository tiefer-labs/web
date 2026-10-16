// SPDX-License-Identifier: MPL-2.0

/** One orbital pass, in milliseconds. Every moving part repeats with this period, so the loop has no seam. */
export const PERIOD = 14_000;
/** The moment of the pass shown to visitors who prefer reduced motion: every stage at once. */
export const POSTER = 0.62;

// Earth and both orbits are circles around one centre far below the viewBox, so the limb reads as a curve.
const CENTRE = { x: 720, y: 3628 };
export const EARTH = 2967;
export const ORBIT = EARTH + 455;
export const INNER_ORBIT = EARTH + 330;
export { CENTRE };

const TILES = [3.2, 4.4, 5.6, 6.8];
const TARGET = 8.4;
// Left of the frames the satellite scans, so the downlink crosses the swath's path.
const STATION = 1.9;
const SWATH = 0.75;

// Starts right of the hero copy, so the beam never crosses the text.
const satAngle = (p: number) => 1 + 14 * p;
/** When the satellite is straight above a point at this angle. */
const overhead = (angle: number) => (angle - 1) / 14;

const DETECT = overhead(TARGET);
const SEND = DETECT + 0.04;
const ARRIVE = SEND + 0.11;

type Point = { x: number; y: number };

const round = (v: number) => Math.round(v * 100) / 100;

function polar(radius: number, degrees: number): Point {
  const a = (degrees * Math.PI) / 180;
  return { x: round(CENTRE.x + radius * Math.sin(a)), y: round(CENTRE.y - radius * Math.cos(a)) };
}

const lerp = (a: Point, b: Point, k: number): Point => ({
  x: round(a.x + (b.x - a.x) * k),
  y: round(a.y + (b.y - a.y) * k),
});
const clamp = (v: number) => Math.min(1, Math.max(0, v));
/** 0 before `a`, 1 after `b`, smooth in between. */
function ramp(p: number, a: number, b: number) {
  const v = clamp((p - a) / (b - a));
  return v * v * (3 - 2 * v);
}
/** Rises from `a` to `b`, holds, falls from `c` to `d`. */
const span = (p: number, a: number, b: number, c: number, d: number) => ramp(p, a, b) * (1 - ramp(p, c, d));
const wave = (p: number, cycles: number, phase: number) => 0.5 + 0.5 * Math.sin(2 * Math.PI * (p * cycles + phase));

/** Every position and opacity of the hero illustration at point `p` (0 to 1) of the pass. */
export function heroScene(p: number) {
  const angle = satAngle(p);
  const sat = polar(ORBIT, angle);
  const visible = ramp(p, 0, 0.07) * (1 - ramp(p, 0.95, 1));
  const reset = 1 - ramp(p, 0.93, 0.99);

  const tiles = TILES.map((tileAngle) => {
    const at = overhead(tileAngle);
    const ground = polar(EARTH, tileAngle);
    const rise = ramp(p, at, at + 0.045);
    return {
      ground,
      angle: tileAngle,
      opacity: round(0.22 + (0.25 * ramp(p, at - 0.004, at) + 0.5 * span(p, at - 0.004, at, at, at + 0.04)) * reset),
      frame: rise > 0 && rise < 1 ? lerp(ground, sat, rise) : null,
      kept: rise >= 1 && reset > 0.5,
    };
  });

  const lock = ramp(p, DETECT - 0.025, DETECT + 0.005);
  const station = polar(EARTH, STATION);
  const linkEnd = polar(ORBIT, satAngle(Math.min(p, ARRIVE)));
  const travel = ramp(p, SEND, ARRIVE);
  const ping = (delay: number) => {
    const k = clamp((p - ARRIVE - delay) / 0.12);
    return { r: round(4 + 30 * k), opacity: round(k > 0 && k < 1 ? 0.6 * (1 - k) : 0) };
  };
  const linkAngle = (Math.atan2(linkEnd.y - station.y, linkEnd.x - station.x) * 180) / Math.PI;
  const labelAt = lerp(station, linkEnd, 0.4);
  const normal = ((linkAngle - 90) * Math.PI) / 180;

  return {
    stars: (index: number) => round(0.4 * (0.6 + 0.4 * wave(p, 2 + (index % 3), index * 0.37))),
    sat: { ...sat, angle: round(angle), opacity: round(visible) },
    trail: [0, 1, 2, 3].map((i) => ({
      from: polar(ORBIT, angle - 0.8 * (i + 1)),
      to: polar(ORBIT, angle - 0.8 * i),
      opacity: round(visible * (0.5 - i * 0.12)),
    })),
    swath: {
      points: [sat, polar(EARTH, angle - SWATH), polar(EARTH, angle + SWATH)].map((pt) => `${pt.x},${pt.y}`).join(" "),
      opacity: round(
        visible *
          (0.3 +
            0.35 *
              Math.max(
                ...TILES.map((t) => span(p, overhead(t) - 0.01, overhead(t), overhead(t), overhead(t) + 0.03)),
                span(p, DETECT - 0.02, DETECT, DETECT, DETECT + 0.05),
              )),
      ),
    },
    tiles,
    kept: tiles.filter((tile) => tile.kept).length,
    target: {
      ...polar(EARTH, TARGET),
      angle: TARGET,
      detected: round(span(p, DETECT, DETECT + 0.01, 0.9, 0.95)),
      pulse: { r: round(3 + 7 * ((p * 7) % 1)), opacity: round(0.5 * (1 - ((p * 7) % 1))) },
      reticle: { scale: round(1 + 1.4 * (1 - lock)), opacity: round(lock * (1 - ramp(p, 0.9, 0.95))) },
    },
    link: {
      from: linkEnd,
      to: station,
      opacity: round(0.55 * span(p, SEND - 0.015, SEND, ARRIVE + 0.01, ARRIVE + 0.05)),
      offset: round(-p * 600),
      packet: travel > 0 && travel < 1 ? lerp(sat, station, travel) : null,
      label: {
        x: round(labelAt.x + 14 * Math.cos(normal)),
        y: round(labelAt.y + 14 * Math.sin(normal)),
        angle: round(linkAngle),
        opacity: round(0.85 * span(p, SEND, SEND + 0.02, ARRIVE + 0.06, ARRIVE + 0.1)),
      },
    },
    station: {
      ...station,
      angle: STATION,
      lit: round(span(p, ARRIVE, ARRIVE + 0.01, ARRIVE + 0.08, ARRIVE + 0.14)),
      pings: [ping(0), ping(0.035)],
    },
    companion: { ...polar(INNER_ORBIT, 14 - 17 * p), opacity: round(0.6 * ramp(p, 0, 0.1) * (1 - ramp(p, 0.9, 1))) },
  };
}

export type HeroScene = ReturnType<typeof heroScene>;
