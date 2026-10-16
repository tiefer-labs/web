// SPDX-License-Identifier: MPL-2.0
import { geoOrthographic } from "d3-geo";
import { expect, test } from "vitest";
import { GLOBE, project } from "@/components/home/globe";
import { FIRE, heroScene, POSTER, RESPONDER } from "@/components/home/hero-scene";

const steps = Array.from({ length: 400 }, (_, i) => i / 400);

test("projects places exactly like d3's orthographic projection, which drew the land dots", () => {
  const d3 = geoOrthographic().rotate([GLOBE.lambda, GLOBE.phi]).scale(GLOBE.radius).translate([GLOBE.x, GLOBE.y]);
  for (const place of [
    [-9.1, 38.7],
    [2.35, 48.85],
    [FIRE.lon, FIRE.lat],
    [RESPONDER.lon, RESPONDER.lat],
    [37.6, 55.7],
  ] as const) {
    const [x, y] = d3([...place]) ?? [Number.NaN, Number.NaN];
    expect(project(place[0], place[1])).toMatchObject({
      x: expect.closeTo(x, 1),
      y: expect.closeTo(y, 1),
      front: true,
    });
  }
});

test("the loop has no seam: the satellite is out of sight and the store empty where the pass wraps", () => {
  for (const p of [0, 0.9999]) {
    expect(heroScene(p).sat.opacity).toBe(0);
    expect(heroScene(p).store.kept).toBe(0);
    expect(heroScene(p).stage).toBe("scanning");
  }
});

test("the satellite flies west to east and keeps every frame before anything is sent", () => {
  const xs = steps.map((p) => heroScene(p).sat.x);
  for (let i = 1; i < xs.length; i++) expect(xs[i]).toBeGreaterThan(xs[i - 1] ?? 0);
  const firstPacket = steps.find((p) => heroScene(p).link.packet) ?? 1;
  const { store } = heroScene(firstPacket);
  expect(store.kept).toBe(store.of);
  expect(store.sent).toBe(false);
});

test("the stages run in order: scanning, wildfire detected, emergency services alerted", () => {
  const stages = steps.map((p) => heroScene(p).stage).filter((stage, i, all) => stage !== all[i - 1]);
  expect(stages).toEqual(["scanning", "detected", "alerted", "scanning"]);
});

test("the alert packet travels down to the emergency services", () => {
  const packets = steps.flatMap((p) => heroScene(p).link.packet ?? []);
  const { responder } = heroScene(0);
  const distance = (pt: { x: number; y: number }) => Math.hypot(pt.x - responder.x, pt.y - responder.y);
  expect(packets.length).toBeGreaterThan(10);
  expect(distance(packets.at(-1) ?? responder)).toBeLessThan(distance(packets[0] ?? responder));
});

test("the still frame shows the whole story at once", () => {
  const scene = heroScene(POSTER);
  expect(scene.stage).toBe("detected");
  expect(scene.sat.opacity).toBe(1);
  expect(scene.fire.reticle.opacity).toBeGreaterThan(0.9);
  expect(scene.store.kept).toBe(scene.store.of);
  expect(scene.link.packet).not.toBeNull();
  expect(scene.link.label.opacity).toBe(1);
});

test("the emergency services are named only while the satellite is far from them", () => {
  const gaps = steps
    .map((p) => heroScene(p))
    .filter(({ sat, responder }) => responder.named > 0 && sat.opacity > 0)
    .map(({ sat, responder }) => Math.hypot(sat.x - responder.x, sat.y - responder.y));
  expect(gaps.length).toBeGreaterThan(0);
  expect(Math.min(...gaps)).toBeGreaterThan(180);
});

test("the alert label never runs into the emergency services label", () => {
  // The responder label sits left of its marker; the alert label is about 90 units long on its line.
  for (const p of steps) {
    const { link, responder } = heroScene(p);
    if (link.label.opacity <= 0) continue;
    const a = (link.label.angle * Math.PI) / 180;
    for (let d = -45; d <= 45; d += 9) {
      const x = link.label.x + d * Math.cos(a);
      const y = link.label.y + d * Math.sin(a);
      const inside = x > responder.x - 172 && x < responder.x - 12 && y > responder.y - 12 && y < responder.y + 6;
      expect(inside).toBe(false);
    }
  }
});
