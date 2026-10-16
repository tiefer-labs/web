// SPDX-License-Identifier: MPL-2.0
import { expect, test } from "vitest";
import { heroScene, POSTER } from "@/components/home/hero-scene";

const steps = Array.from({ length: 200 }, (_, i) => i / 200);

test("the loop has no seam: the satellite is out of sight where the pass wraps around", () => {
  expect(heroScene(0).sat.opacity).toBe(0);
  expect(heroScene(0.9999).sat.opacity).toBe(0);
  expect(heroScene(0).kept).toBe(0);
});

test("frames are kept on board one by one before the alert goes down", () => {
  const kept = steps.map((p) => heroScene(p).kept);
  const firstPacket = steps.findIndex((p) => heroScene(p).link.packet);
  expect(Math.max(...kept)).toBe(4);
  expect(kept[firstPacket]).toBe(4);
  for (let i = 1; i < firstPacket; i++) expect(kept[i]).toBeGreaterThanOrEqual(kept[i - 1] ?? 0);
});

test("the alert packet travels from the satellite to the ground station", () => {
  const packets = steps.flatMap((p) => heroScene(p).link.packet ?? []);
  const { station } = heroScene(0);
  expect(packets.length).toBeGreaterThan(5);
  const distance = (pt: { x: number; y: number }) => Math.hypot(pt.x - station.x, pt.y - station.y);
  expect(distance(packets.at(-1) ?? station)).toBeLessThan(distance(packets[0] ?? station));
});

test("the still frame shows the whole story at once", () => {
  const scene = heroScene(POSTER);
  expect(scene.sat.opacity).toBe(1);
  expect(scene.kept).toBe(4);
  expect(scene.target.reticle.opacity).toBeGreaterThan(0.9);
  expect(scene.link.packet).not.toBeNull();
  expect(scene.link.label.opacity).toBeGreaterThan(0.8);
});
