// SPDX-License-Identifier: MPL-2.0
import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";
import "./shared";

afterEach(() => {
  cleanup();
  document.body.innerHTML = "";
});
