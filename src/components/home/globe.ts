// SPDX-License-Identifier: MPL-2.0

/** An orthographic view of Europe from space: the globe's centre, radius and rotation, in viewBox units. */
export const GLOBE = { x: 1080, y: 1000, radius: 1000, lambda: -14, phi: -12 } as const;

const RAD = Math.PI / 180;
const COS_PHI = Math.cos(GLOBE.phi * RAD);
const SIN_PHI = Math.sin(GLOBE.phi * RAD);

/** Both hero layers share this box, so the moving parts line up with the map. */
export const heroSvgClass = "absolute inset-0 size-full max-lg:top-auto max-lg:h-85";

export type Point = { x: number; y: number };

const round = (v: number) => Math.round(v * 100) / 100;

/** Projects a place on the globe, as d3's geoOrthographic does. */
export function project(lon: number, lat: number): Point & { front: boolean } {
  const l = (lon + GLOBE.lambda) * RAD;
  const cosLat = Math.cos(lat * RAD);
  const x = cosLat * Math.cos(l);
  const y = cosLat * Math.sin(l);
  const z = Math.sin(lat * RAD);
  const r = GLOBE.radius;
  return {
    x: round(GLOBE.x + r * y),
    y: round(GLOBE.y - r * (z * COS_PHI + x * SIN_PHI)),
    front: x * COS_PHI - z * SIN_PHI > 0,
  };
}
