// SPDX-License-Identifier: MPL-2.0
import { geoContains, geoOrthographic } from "d3-geo";
import { feature } from "topojson-client";
import type { GeometryCollection, Topology } from "topojson-specification";
import { expect, test } from "vitest";
import land50m from "world-atlas/land-50m.json";
import { GLOBE, project } from "@/components/home/globe";
import { type HeroScene, heroScene, PASSES, POSTER, SHIPS } from "@/components/home/hero-scene";

const steps = Array.from({ length: 1500 }, (_, i) => i / 1500);
const scenes = steps.map((p) => heroScene(p));
const EVENTS = ["wildfire", "oilSpill", "vessels"] as const;
/** The status panel, in viewBox units. */
const PANEL = { left: 1102, right: 1388, top: 177, bottom: 301 };

test("projects places exactly like d3's orthographic projection, which drew the land dots", () => {
  const d3 = geoOrthographic().rotate([GLOBE.lambda, GLOBE.phi]).scale(GLOBE.radius).translate([GLOBE.x, GLOBE.y]);
  const places = PASSES.flatMap((pass) => [...pass.targets, pass.responder]);
  for (const { lon, lat } of places) {
    const [x, y] = d3([lon, lat]) ?? [Number.NaN, Number.NaN];
    expect(project(lon, lat)).toMatchObject({ x: expect.closeTo(x, 1), y: expect.closeTo(y, 1), front: true });
  }
});

test("plays the wildfire, the oil spill and the ships in order, then starts again", () => {
  const events = scenes.map((scene) => scene.event).filter((event, i, all) => event !== all[i - 1]);
  expect(events).toEqual([...EVENTS]);
});

test("every pass runs scanning, detected, alerted, and fades out before the next begins", () => {
  const twice = ["scanning", "detected", "alerted", "detected", "alerted", "scanning"];
  for (const event of EVENTS) {
    const pass = scenes.filter((scene) => scene.event === event);
    const stages = pass.map((scene) => scene.stage).filter((stage, i, all) => stage !== all[i - 1]);
    // The ships pass finds two ships without AIS and alerts for each.
    expect(stages).toEqual(event === "vessels" ? twice : ["scanning", "detected", "alerted", "scanning"]);
    expect(pass[0]?.sat.opacity).toBe(0);
    expect(pass.at(-1)?.sat.opacity).toBe(0);
    expect(pass.at(-1)?.store.kept).toBe(0);
  }
});

test("each satellite flies west to east and keeps imaging after it sends the alert", () => {
  for (const event of EVENTS) {
    const pass = scenes.filter((scene) => scene.event === event);
    for (let i = 1; i < pass.length; i++) expect(pass[i]?.sat.x).toBeGreaterThan(pass[i - 1]?.sat.x ?? 0);
    const sending = pass.find((scene) => scene.link.packet);
    const kept = Math.max(...pass.map((scene) => scene.store.kept));
    expect(sending?.store.sent).toBe(0);
    expect(kept).toBe(sending?.store.of);
    expect(kept).toBeGreaterThan(sending?.store.kept ?? kept);
  }
});

test("the ships pass sends one 48 KB alert for each ship without AIS", () => {
  const pass = scenes.filter((scene) => scene.event === "vessels");
  expect(SHIPS.filter((ship) => ship.dark)).toHaveLength(2);
  expect(Math.max(...pass.map((scene) => scene.store.sent))).toBe(2);
  expect(pass.at(-1)?.store.sent).toBe(0);
});

test("every ship is at sea, under the ground track, inside a frame", () => {
  // oxlint-disable-next-line typescript/no-unsafe-type-assertion
  const topology = land50m as unknown as Topology<{ land: GeometryCollection }>;
  const land = feature(topology, topology.objects.land);
  const pass = PASSES.find(({ event }) => event === "vessels");
  for (const ship of SHIPS) {
    expect(geoContains(land, [ship.lon, ship.lat]), `${ship.lon} E is on land`).toBe(false);
    const frame = pass?.frames.find((lon) => Math.abs(lon - ship.lon) < 1.9);
    expect(frame, `no frame over ${ship.lon} E`).toBeDefined();
    expect(Math.abs(ship.lat - (pass?.track(ship.lon) ?? 0))).toBeLessThan(1.5);
  }
});

test("the satellites never fly behind the status panel", () => {
  const hidden = scenes
    .filter(({ sat }) => sat.opacity > 0)
    .filter(
      ({ sat }) =>
        sat.x > PANEL.left - 30 && sat.x < PANEL.right + 30 && sat.y > PANEL.top - 30 && sat.y < PANEL.bottom + 30,
    )
    .map(({ event, sat }) => `${event} ${sat.x},${sat.y}`);
  expect(hidden).toEqual([]);
});

test("the alert packet travels down to the responder of its event", () => {
  for (const event of EVENTS) {
    const pass = scenes.filter((scene) => scene.event === event);
    const packets = pass.flatMap((scene) => scene.link.packet ?? []);
    const responder = pass[0]?.responder ?? { x: 0, y: 0 };
    const distance = (pt: { x: number; y: number }) => Math.hypot(pt.x - responder.x, pt.y - responder.y);
    expect(packets.length).toBeGreaterThan(5);
    expect(distance(packets.at(-1) ?? responder)).toBeLessThan(distance(packets[0] ?? responder));
  }
});

/** The responder label, about 170 units long, left of its marker or below it. */
function labelBox(responder: HeroScene["responder"]) {
  if (responder.side === "below")
    return { l: responder.x - 6, r: responder.x + 170, t: responder.y + 8, b: responder.y + 24 };
  return { l: responder.x - 182, r: responder.x - 12, t: responder.y - 12, b: responder.y + 6 };
}

test("a satellite never flies across the name of a responder", () => {
  // The satellite and its panels fit in 64 by 64 units.
  const crossed = scenes
    .filter((scene) => scene.responder.named > 0 && scene.sat.opacity > 0)
    .filter((scene) => {
      const box = labelBox(scene.responder);
      const { x, y } = scene.sat;
      return x + 32 > box.l && x - 32 < box.r && y + 32 > box.t && y - 32 < box.b;
    })
    .map(({ event }) => event);
  expect(crossed).toEqual([]);
});

test("no responder label sits under the status panel", () => {
  const covered = scenes
    .filter(({ responder }) => responder.named > 0)
    .filter(({ responder }) => {
      const box = labelBox(responder);
      return box.r > PANEL.left && box.l < PANEL.right && box.b > PANEL.top && box.t < PANEL.bottom;
    })
    .map(({ event }) => event);
  expect(covered).toEqual([]);
});

test("the alert label never runs into the responder label", () => {
  const clashes = scenes.filter(({ link, responder }) => {
    if (link.label.opacity <= 0) return false;
    const box = labelBox(responder);
    const a = (link.label.angle * Math.PI) / 180;
    // The alert label is a band about 90 long and 12 high on its line, drawn above its baseline.
    return Array.from({ length: 11 }, (_, i) => -45 + 9 * i).some((d) =>
      [0, -6, -12].some((h) => {
        const x = link.label.x + d * Math.cos(a) + h * Math.cos(a + Math.PI / 2);
        const y = link.label.y + d * Math.sin(a) + h * Math.sin(a + Math.PI / 2);
        return x > box.l && x < box.r && y > box.t && y < box.b;
      }),
    );
  });
  expect(clashes.map(({ event }) => event)).toEqual([]);
});

test("the alert label always reads left to right", () => {
  const upsideDown = scenes
    .filter(({ link }) => link.label.opacity > 0 && Math.abs(link.label.angle) > 90)
    .map(({ event }) => event);
  expect(upsideDown).toEqual([]);
});

test("the still frame shows the wildfire alert on its way down", () => {
  const scene = heroScene(POSTER);
  expect(scene.event).toBe("wildfire");
  expect(scene.stage).toBe("detected");
  expect(scene.targets[0]?.reticle.opacity).toBeGreaterThan(0.9);
  expect(scene.link.packet).not.toBeNull();
});
