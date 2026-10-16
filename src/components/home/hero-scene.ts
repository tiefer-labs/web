// SPDX-License-Identifier: MPL-2.0
import { type Point, project } from "./globe";

export type EventId = "wildfire" | "oilSpill" | "vessels";
export type Stage = "scanning" | "detected" | "alerted";
type Place = { lon: number; lat: number };

/** One pass, in milliseconds. */
export const PASS = 16_000;
/** The three passes play in order, then start again; every moving part repeats with this period. */
export const PERIOD = 3 * PASS;
/** The moment shown to visitors who prefer reduced motion: the wildfire alert on its way down. */
export const POSTER = 0.79 / 3;

type Pass = {
  event: EventId;
  /** Longitudes the satellite crosses, west to east. */
  from: number;
  to: number;
  track: (lon: number) => number;
  /** How far above its ground track the satellite is drawn. */
  lift: number;
  frames: number[];
  /** What the satellite finds, in the order it flies over them; each one sends its own alert. */
  targets: Place[];
  responder: Place & { label: "left" | "below" };
};

/** A curve through 38.6 N at 12 W, 38.3 N at 6 E and 43.4 N at 33 E. */
function shipsTrack(lon: number) {
  return 38.6 - 0.0167 * (lon + 12) + 0.00457 * (lon + 12) * (lon - 6);
}

export const PASSES: Pass[] = [
  {
    event: "wildfire",
    from: -15,
    to: 37,
    track: (lon) => 44.5 - 0.18 * (lon + 16),
    lift: 150,
    frames: [-5.5, -1, 3.5, 8, 12.5, 17, 22, 26.5, 31],
    targets: [{ lon: 22, lat: 37.55 }],
    // The EU Emergency Response Coordination Centre, Brussels.
    responder: { lon: 4.35, lat: 50.85, label: "left" },
  },
  {
    event: "oilSpill",
    // Across Central Europe and down to the Black Sea, below the status panel.
    from: -10,
    to: 46,
    track: (lon) => 50 - 0.15 * (lon + 10),
    lift: 100,
    frames: [-4, 2, 8, 14, 20, 25.5, 30.6, 36],
    targets: [{ lon: 30.6, lat: 43.9 }],
    // The European Maritime Safety Agency, Lisbon. The alert arrives from the east, so the label goes below.
    responder: { lon: -9.14, lat: 38.72, label: "below" },
  },
  {
    event: "vessels",
    // From the Atlantic off Portugal through the western Mediterranean and the Tyrrhenian Sea into the Black Sea.
    // It runs on past the view, so the second alert lands before the pass ends.
    from: -14,
    to: 42,
    track: shipsTrack,
    lift: 85,
    frames: [-10, -5.5, -1.5, 2.2, 6, 10.5, 15, 19.5, 23.5, 27.3, 31, 35],
    targets: [
      { lon: 6, lat: shipsTrack(6) - 0.4 },
      { lon: 31, lat: shipsTrack(31) + 0.2 },
    ],
    // Frontex, the EU border and coast guard agency, Warsaw.
    responder: { lon: 21.01, lat: 52.23, label: "left" },
  },
];

/**
 * Ships in the third pass, each in its own sea and right under the ground track, so a frame covers it: the Atlantic,
 * the western Mediterranean, the Tyrrhenian Sea and the Black Sea. The two at the targets do not broadcast AIS.
 */
export const SHIPS = [
  { lon: -10, lat: shipsTrack(-10) + 0.3, heading: 20 },
  { lon: 6, lat: shipsTrack(6) - 0.4, heading: 250, dark: true },
  { lon: 15, lat: shipsTrack(15), heading: 120 },
  { lon: 31, lat: shipsTrack(31) + 0.2, heading: 60, dark: true },
];
const CLOUDS = [
  { lon: 1, lat: 42.6, size: 1 },
  { lon: 12, lat: 39.8, size: 0.8 },
  { lon: 19, lat: 47.5, size: 1.1 },
  { lon: -4, lat: 49, size: 0.7 },
];
const FRAME_HALF = { lon: 1.9, lat: 1.5 };

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

/** Where the satellite of a pass appears when it is over a longitude. */
export function orbitPoint(pass: Pass, lon: number): Point {
  const ground = project(lon, pass.track(lon));
  return { x: ground.x, y: round(ground.y - pass.lift) };
}

/** The dotted orbit of each pass, from just before it enters the view to just after it leaves. */
export const ORBITS = PASSES.map((pass) =>
  Array.from({ length: 61 }, (_, i) => orbitPoint(pass, pass.from - 15 + ((pass.to - pass.from + 30) * i) / 60))
    .map(({ x, y }, i) => `${i ? "L" : "M"}${x} ${y}`)
    .join(""),
);

/** The outline of the scene a frame covers on the ground. */
function footprint(pass: Pass, lon: number) {
  const lat = pass.track(lon);
  return [
    project(lon - FRAME_HALF.lon, lat + FRAME_HALF.lat),
    project(lon + FRAME_HALF.lon, lat + FRAME_HALF.lat),
    project(lon + FRAME_HALF.lon, lat - FRAME_HALF.lat),
    project(lon - FRAME_HALF.lon, lat - FRAME_HALF.lat),
  ];
}

/** Where the first target of a pass is on screen, optionally moved by a few degrees. */
function placeOf(index: number, dLon = 0, dLat = 0) {
  const place = PASSES[index]?.targets[0] ?? { lon: 0, lat: 0 };
  return project(place.lon + dLon, place.lat + dLat);
}

/** Every position and opacity of the hero illustration at point `p` (0 to 1) of the whole loop. */
export function heroScene(p: number) {
  const index = Math.min(PASSES.length - 1, Math.floor(p * PASSES.length));
  const pass = PASSES[index] ?? PASSES[0];
  const q = p * PASSES.length - index;
  const lonAt = (k: number) => pass.from + (pass.to - pass.from) * k;
  /** When the satellite is over a longitude, as a point of this pass. */
  const overhead = (lon: number) => (lon - pass.from) / (pass.to - pass.from);

  const lon = lonAt(q);
  const sat = orbitPoint(pass, lon);
  const ahead = orbitPoint(pass, lon + 0.5);
  const visible = ramp(q, 0, 0.05) * (1 - ramp(q, 0.94, 0.99));
  // The scene of each event fades in as its pass starts and out as it ends, so the passes cross-fade.
  const present = ramp(q, 0, 0.06) * (1 - ramp(q, 0.95, 1));
  const reset = 1 - ramp(q, 0.93, 0.99);
  const under = footprint(pass, lon);
  const flash = (at: number) => span(q, at - 0.012, at, at, at + 0.02);

  const frames = pass.frames.map((frameLon) => {
    const at = overhead(frameLon);
    return {
      lon: frameLon,
      points: points(footprint(pass, frameLon)),
      opacity: round((0.3 * ramp(q, at - 0.004, at) + 0.6 * span(q, at - 0.004, at, at, at + 0.03)) * reset),
      kept: q >= at && reset > 0.5,
    };
  });
  const kept = frames.filter((frame) => frame.kept).length;

  const wave = (cycles: number, phase: number) => Math.sin(2 * Math.PI * (p * cycles + phase));
  const responder = project(pass.responder.lon, pass.responder.lat);
  const alerts = pass.targets.map((place) => {
    const detect = overhead(place.lon);
    return { place, detect, send: detect + 0.03, arrive: detect + 0.11 };
  });
  const first = alerts[0] ?? { place: pass.responder, detect: 1, send: 1, arrive: 1 };
  // The alert in flight, or the last one sent.
  const alert = alerts.findLast(({ send }) => q >= send - 0.01) ?? first;
  const linkFrom = orbitPoint(pass, lonAt(Math.min(q, alert.arrive)));
  const travel = ramp(q, alert.send, alert.arrive);
  const lineAngle = (Math.atan2(linkFrom.y - responder.y, linkFrom.x - responder.x) * 180) / Math.PI;
  // Turned by half a circle when the line points left, so the label never reads upside down.
  const linkAngle = lineAngle > 90 ? lineAngle - 180 : lineAngle < -90 ? lineAngle + 180 : lineAngle;
  const labelAt = lerp(responder, linkFrom, 0.5);
  // Below the line: the label hangs from its baseline, so it sits clear of the line and the responder's name.
  const normal = ((linkAngle + 90) * Math.PI) / 180;
  const ping = (delay: number) => {
    const k = clamp((q - alert.arrive - delay) / 0.1);
    return { r: round(5 + 34 * k), opacity: round(k > 0 && k < 1 ? 0.8 * (1 - k) : 0) };
  };
  const sent = alerts.filter(({ arrive }) => q >= arrive).length;
  const only = (event: EventId) => (pass.event === event ? round(present) : 0);

  const ended = q >= 0.97;
  const found = alerts.findLast(({ detect }) => q >= detect);
  const stage: Stage = ended || !found ? "scanning" : q >= found.arrive ? "alerted" : "detected";

  return {
    event: pass.event,
    stage,
    orbits: PASSES.map((_, i) => (i === index ? round(0.4 * present) : 0)),
    sat: {
      x: sat.x,
      y: sat.y,
      angle: round((Math.atan2(ahead.y - sat.y, ahead.x - sat.x) * 180) / Math.PI),
      opacity: round(visible),
    },
    beam: {
      points: points([sat, under[3] ?? sat, under[2] ?? sat]),
      footprint: points(under),
      opacity: round(visible * (0.5 + 0.5 * Math.max(...pass.frames.map((f) => flash(overhead(f)))))),
    },
    frames,
    store: { kept, of: pass.frames.length, sent: reset > 0.5 ? sent : 0 },
    clouds: CLOUDS.map(({ lon: cloudLon, lat, size }, i) => ({
      ...project(cloudLon + 0.6 * wave(3, i * 0.25), lat),
      size,
    })),
    targets: alerts.map(({ place, detect }) => {
      const lock = ramp(q, detect - 0.02, detect + 0.005);
      return {
        ...project(place.lon, place.lat),
        detected: round(span(q, detect, detect + 0.01, 0.95, 0.99)),
        reticle: { scale: round(1 + 1.6 * (1 - lock)), opacity: round(lock * (1 - ramp(q, 0.95, 0.99))) },
      };
    }),
    fire: {
      ...placeOf(0),
      opacity: only("wildfire"),
      flame: round(1 + 0.08 * wave(69, 0) + 0.05 * wave(111, 0.3)),
      glow: round(0.35 + 0.15 * wave(27, 0)),
    },
    spill: {
      ...placeOf(1),
      opacity: only("oilSpill"),
      // The slick keeps spreading through the pass.
      scale: round(0.75 + 0.45 * (index === 1 ? q : 0)),
      rig: placeOf(1, 1.3, 0.6),
    },
    ships: SHIPS.map((ship, i) => {
      const at = overhead(ship.lon);
      const pulse = (q * 8 + i * 0.3) % 1;
      return {
        ...project(ship.lon, ship.lat),
        heading: ship.heading,
        dark: ship.dark === true,
        opacity: only("vessels"),
        // Ships that broadcast AIS are matched as the satellite sees them; the silent ones are the events.
        matched: ship.dark ? 0 : round(index === 2 ? span(q, at, at + 0.01, 0.95, 0.99) : 0),
        signal: { r: round(5 + 9 * pulse), opacity: round(0.5 * (1 - pulse)) },
      };
    }),
    link: {
      from: linkFrom,
      to: responder,
      opacity: round(0.8 * span(q, alert.send - 0.01, alert.send, alert.arrive + 0.01, alert.arrive + 0.05)),
      offset: round(-q * 800),
      packet: travel > 0 && travel < 1 ? lerp(sat, responder, travel) : null,
      label: {
        x: round(labelAt.x + 20 * Math.cos(normal)),
        y: round(labelAt.y + 20 * Math.sin(normal)),
        angle: round(linkAngle),
        opacity: round(span(q, alert.send, alert.send + 0.015, alert.arrive + 0.06, alert.arrive + 0.1)),
      },
    },
    responder: {
      ...responder,
      side: pass.responder.label,
      marker: round(present),
      // Named only once the alert is on its way: earlier, the satellite may fly across this spot.
      named: round(span(q, first.send - 0.02, first.send, 0.95, 0.99)),
      lit: round(span(q, first.arrive, first.arrive + 0.01, 0.95, 0.99)),
      pings: [ping(0), ping(0.04), ping(0.08)],
    },
  };
}

export type HeroScene = ReturnType<typeof heroScene>;
