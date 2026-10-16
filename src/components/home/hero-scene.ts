// SPDX-License-Identifier: MPL-2.0
import { type Point, project } from "./globe";

/** One pass over Europe, in milliseconds. Every moving part repeats with this period, so the loop has no seam. */
export const PERIOD = 16_000;
/** The moment of the pass shown to visitors who prefer reduced motion: every stage at once. */
export const POSTER = 0.79;

export const FIRE = { lon: 22, lat: 37.55 };
// Brussels: the alert crosses Europe to reach the emergency services.
export const RESPONDER = { lon: 4.35, lat: 50.85 };

const FRAMES = [-5.5, -1, 3.5, 8, 12.5, 17, FIRE.lon];
const CLOUDS = [
  { lon: 1, lat: 42.6, size: 1 },
  { lon: 12, lat: 39.8, size: 0.8 },
  { lon: 19, lat: 47.5, size: 1.1 },
  { lon: -4, lat: 49, size: 0.7 },
];
const FRAME_HALF = { lon: 1.9, lat: 1.5 };

/** The ground track: west to east, drifting south over the Mediterranean. */
export const trackLat = (lon: number) => 44.5 - 0.18 * (lon + 16);
// Seen from this high, a real orbit would sit almost on the map; lifting it on screen shows the satellite in space.
const LIFT = 150;
/** Where the satellite appears when it is over a longitude. */
export function orbitPoint(lon: number): Point {
  const ground = project(lon, trackLat(lon));
  return { x: ground.x, y: ground.y - LIFT };
}
const satLon = (p: number) => -15 + 52 * p;
/** When the satellite is over a longitude. */
const overhead = (lon: number) => (lon + 15) / 52;

const DETECT = overhead(FIRE.lon);
const SEND = DETECT + 0.03;
const ARRIVE = SEND + 0.08;

const round = (v: number) => Math.round(v * 100) / 100;
const clamp = (v: number) => Math.min(1, Math.max(0, v));
/** 0 before `a`, 1 after `b`, smooth in between. */
function ramp(p: number, a: number, b: number) {
  const v = clamp((p - a) / (b - a));
  return v * v * (3 - 2 * v);
}
/** Rises from `a` to `b`, holds, falls from `c` to `d`. */
const span = (p: number, a: number, b: number, c: number, d: number) => ramp(p, a, b) * (1 - ramp(p, c, d));
const lerp = (a: Point, b: Point, k: number): Point => ({
  x: round(a.x + (b.x - a.x) * k),
  y: round(a.y + (b.y - a.y) * k),
});
const points = (list: Point[]) => list.map(({ x, y }) => `${x},${y}`).join(" ");

/** The outline of the scene a frame covers on the ground. */
function footprint(lon: number) {
  const lat = trackLat(lon);
  return [
    project(lon - FRAME_HALF.lon, lat + FRAME_HALF.lat),
    project(lon + FRAME_HALF.lon, lat + FRAME_HALF.lat),
    project(lon + FRAME_HALF.lon, lat - FRAME_HALF.lat),
    project(lon - FRAME_HALF.lon, lat - FRAME_HALF.lat),
  ];
}

export type Stage = "scanning" | "detected" | "alerted";

/** Every position and opacity of the hero illustration at point `p` (0 to 1) of the pass. */
export function heroScene(p: number) {
  const lon = satLon(p);
  const sat = orbitPoint(lon);
  const ahead = orbitPoint(lon + 0.5);
  const visible = ramp(p, 0, 0.05) * (1 - ramp(p, 0.94, 0.99));
  const reset = 1 - ramp(p, 0.93, 0.99);
  const under = footprint(lon);
  const flash = (at: number) => span(p, at - 0.012, at, at, at + 0.02);

  const frames = FRAMES.map((frameLon) => {
    const at = overhead(frameLon);
    return {
      lon: frameLon,
      points: points(footprint(frameLon)),
      opacity: round((0.3 * ramp(p, at - 0.004, at) + 0.6 * span(p, at - 0.004, at, at, at + 0.03)) * reset),
      kept: p >= at && reset > 0.5,
    };
  });
  const kept = frames.filter((frame) => frame.kept).length;

  const fire = project(FIRE.lon, FIRE.lat);
  const responder = project(RESPONDER.lon, RESPONDER.lat);
  const lock = ramp(p, DETECT - 0.02, DETECT + 0.005);
  const wave = (cycles: number, phase: number) => Math.sin(2 * Math.PI * (p * cycles + phase));
  const linkLon = satLon(Math.min(p, ARRIVE));
  const linkFrom = orbitPoint(linkLon);
  const travel = ramp(p, SEND, ARRIVE);
  const linkAngle = (Math.atan2(linkFrom.y - responder.y, linkFrom.x - responder.x) * 180) / Math.PI;
  const labelAt = lerp(responder, linkFrom, 0.5);
  const normal = ((linkAngle - 90) * Math.PI) / 180;
  const ping = (delay: number) => {
    const k = clamp((p - ARRIVE - delay) / 0.1);
    return { r: round(5 + 34 * k), opacity: round(k > 0 && k < 1 ? 0.8 * (1 - k) : 0) };
  };

  const ended = p >= 0.97;
  const stage: Stage = ended ? "scanning" : p >= ARRIVE ? "alerted" : p >= DETECT ? "detected" : "scanning";

  return {
    stage,
    sat: {
      x: sat.x,
      y: sat.y,
      angle: round((Math.atan2(ahead.y - sat.y, ahead.x - sat.x) * 180) / Math.PI),
      opacity: round(visible),
    },
    beam: {
      points: points([sat, under[3] ?? sat, under[2] ?? sat]),
      footprint: points(under),
      opacity: round(visible * (0.5 + 0.5 * Math.max(...FRAMES.map((f) => flash(overhead(f)))))),
    },
    frames,
    store: { kept, of: FRAMES.length, sent: p >= ARRIVE && reset > 0.5 },
    clouds: CLOUDS.map(({ lon: cloudLon, lat, size }, index) => ({
      ...project(cloudLon + 0.6 * wave(1, index * 0.25), lat),
      size,
    })),
    fire: {
      ...fire,
      flame: round(1 + 0.08 * wave(23, 0) + 0.05 * wave(37, 0.3)),
      glow: round(0.35 + 0.15 * wave(9, 0)),
      detected: round(span(p, DETECT, DETECT + 0.01, 0.95, 0.99)),
      reticle: { scale: round(1 + 1.6 * (1 - lock)), opacity: round(lock * (1 - ramp(p, 0.95, 0.99))) },
    },
    link: {
      from: linkFrom,
      to: responder,
      opacity: round(0.8 * span(p, SEND - 0.01, SEND, ARRIVE + 0.01, ARRIVE + 0.05)),
      offset: round(-p * 800),
      packet: travel > 0 && travel < 1 ? lerp(sat, responder, travel) : null,
      label: {
        x: round(labelAt.x + 14 * Math.cos(normal)),
        y: round(labelAt.y + 14 * Math.sin(normal)),
        angle: round(linkAngle),
        opacity: round(span(p, SEND, SEND + 0.015, ARRIVE + 0.06, ARRIVE + 0.1)),
      },
    },
    responder: {
      ...responder,
      // Named only once the alert is on its way: earlier, the satellite flies across this spot.
      named: round(span(p, SEND - 0.02, SEND, 0.95, 0.99)),
      lit: round(span(p, ARRIVE, ARRIVE + 0.01, 0.95, 0.99)),
      pings: [ping(0), ping(0.04), ping(0.08)],
    },
  };
}

export type HeroScene = ReturnType<typeof heroScene>;
