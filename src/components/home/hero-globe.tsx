// SPDX-License-Identifier: MPL-2.0
import { GLOBE, heroSvgClass } from "./globe";
import { COASTLINE, GRATICULE, LAND_DOTS } from "./globe-data";
import { orbitPoint } from "./hero-scene";

const orbit = Array.from({ length: 81 }, (_, i) => orbitPoint(i - 30))
  .map(({ x, y }, i) => `${i ? "L" : "M"}${x} ${y}`)
  .join("");

/** The still part of the hero: Europe as a dot map on the curve of the Earth, and the satellite's orbit. */
export function HeroGlobe() {
  const halo = GLOBE.radius + 70;
  return (
    <svg
      className={heroSvgClass}
      viewBox="0 0 1440 760"
      preserveAspectRatio="xMaxYMax slice"
      aria-hidden="true"
      focusable="false"
      data-print="hide"
    >
      <defs>
        <radialGradient id="hero-atmosphere" gradientUnits="userSpaceOnUse" cx={GLOBE.x} cy={GLOBE.y} r={halo}>
          <stop offset={GLOBE.radius / halo} className="[stop-color:var(--color-on-ink-muted)] [stop-opacity:0.2]" />
          <stop offset="1" className="[stop-color:var(--color-on-ink-muted)] [stop-opacity:0]" />
        </radialGradient>
        <linearGradient id="hero-fade" gradientUnits="userSpaceOnUse" x1="420" x2="760" y1="0" y2="0">
          <stop offset="0" stopColor="white" stopOpacity="0" />
          <stop offset="1" stopColor="white" />
        </linearGradient>
        {/* Fades the map out behind the hero copy. */}
        <mask id="hero-fade-mask" maskUnits="userSpaceOnUse" x="0" y="0" width="1440" height="760">
          <rect width="1440" height="760" fill="url(#hero-fade)" />
        </mask>
      </defs>

      <g mask="url(#hero-fade-mask)">
        <circle cx={GLOBE.x} cy={GLOBE.y} r={halo} fill="url(#hero-atmosphere)" />
        <circle className="fill-on-ink-muted opacity-[0.05]" cx={GLOBE.x} cy={GLOBE.y} r={GLOBE.radius} />
        <circle
          className="fill-none stroke-on-ink-muted opacity-50"
          strokeWidth="1.25"
          vectorEffect="non-scaling-stroke"
          cx={GLOBE.x}
          cy={GLOBE.y}
          r={GLOBE.radius}
        />
        <path className="fill-none stroke-on-ink-muted opacity-10" vectorEffect="non-scaling-stroke" d={GRATICULE} />
        <path
          className="fill-on-ink-muted stroke-on-ink-muted opacity-35"
          fillOpacity="0.12"
          strokeWidth="1"
          vectorEffect="non-scaling-stroke"
          d={COASTLINE}
        />
        <path
          className="fill-none stroke-on-ink-muted opacity-55"
          strokeWidth="2.6"
          strokeLinecap="round"
          d={LAND_DOTS}
        />
        <path
          className="fill-none stroke-on-ink-muted opacity-40"
          strokeWidth="1.25"
          strokeDasharray="2 7"
          vectorEffect="non-scaling-stroke"
          d={orbit}
        />
      </g>
    </svg>
  );
}
